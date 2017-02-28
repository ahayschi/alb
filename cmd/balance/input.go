package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/parallelworks/alb"
)

type Params struct {
	InputFile string
	CycleTime float64
}

func GetParams() (*Params, error) {
	if len(os.Args) != 3 {
		return nil, errors.New("params: two args required")
	}

	args := os.Args[1:]
	p := &Params{InputFile: args[0]}
	params := strings.Split(args[1], ",")

	if len(params)%2 != 0 {
		return nil, errors.New("params: requires both key and value")
	}

	for i := 0; i < len(params); i = i + 2 {
		key := params[i]
		value := params[i+1]

		if key == "cycle_time" {
			cycleTime, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("params: parse cycle_time:  %s", err)
			}
			p.CycleTime = cycleTime
		}
	}

	return p, nil
}

func GetStream(p *Params) (io.Reader, error) {
	file, err := os.Open(p.InputFile)
	if err != nil {
		return nil, fmt.Errorf("stream: %s", err)
	}

	return file, nil
}

func getLines(in io.Reader) ([]string, error) {
	var lines []string

	spec := bufio.NewScanner(in)
	for spec.Scan() {
		lines = append(lines, strings.Trim(spec.Text(), "\n\r"))
	}

	if err := spec.Err(); err != nil {
		return nil, fmt.Errorf("lines: file: %s", err)
	}

	if len(lines) == 0 {
		return nil, errors.New("lines: file is empty")
	}

	return lines, nil
}

func taskById(id int, tasks []*alb.Task) *alb.Task {
	for _, task := range tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func ParseIn2File(in io.Reader) ([]*alb.Task, []*alb.Station, error) {
	lines, err := getLines(in)
	if err != nil {
		return nil, nil, err
	}

	nTasks, err := strconv.ParseInt(lines[0], 10, 32)
	if err != nil {
		return nil, nil, fmt.Errorf("parse: number of tasks: %s", err)
	}

	times := lines[1 : nTasks+1]
	predIds := lines[nTasks+1:]

	tasks := make([]*alb.Task, nTasks)
	stations := make([]*alb.Station, nTasks)
	for i, t := range times {
		parts := strings.Split(t, ",")
		if err != nil {
			return nil, nil, fmt.Errorf("parse: times: %s", err)
		}
		if len(parts) == 2 {
			taskId, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, nil, fmt.Errorf("parse: taskId: %s", err)
			}
			taskTime, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, nil, fmt.Errorf("parse: taskTime: %s", err)
			}
			task := alb.NewTask(taskId, taskTime)
			tasks[i] = task
		} else if len(parts) == 1 {
			taskTime, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("parse: task time: line %d: %s", i+1, err)
			}
			task := alb.NewTask(i+1, taskTime)
			tasks[i] = task
		} else {
			return nil, nil, fmt.Errorf("parse: preds: %s", err)
		}

		station := alb.NewStation(i + 1)
		stations[i] = station

	}

	for _, ids := range predIds {
		parts := strings.Split(ids, ",")
		if err != nil {
			return nil, nil, fmt.Errorf("parse: preds: %s", err)
		}

		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("parse: preds: %s", err)
		}

		predId, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, fmt.Errorf("parse: pred: %s", err)
		}

		taskId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("parse: pred: %s", err)
		}

		if predId == -1 && taskId == -1 {
			continue
		}

		pred := taskById(predId, tasks)
		if pred == nil {
			return nil, nil, errors.New("parse: pred: attempting to assign an unknown task")
		}

		task := taskById(taskId, tasks)
		if task == nil {
			return nil, nil, errors.New("parse pred: attempting to assign pred to unknown task")
		}

		task.AddPred(pred)
	}

	return tasks, stations, nil
}

// A hack to constrain the problem to always be valid for a paced line.
func ValidateParams(params *Params, line *alb.Line) {
	for _, task := range line.Tasks() {
		ttime := task.Time()
		if ttime > params.CycleTime {
			log.WithFields(log.Fields{
				"task":       task.ID,
				"task_time":  ttime,
				"cycle_time": params.CycleTime,
			}).Warnf("Cycle time is being bumped to task_time")
			params.CycleTime = ttime
		}
	}
}
