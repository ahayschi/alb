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

var (
	heuristicMap = map[string]alb.Heuristic{
		"ShortestTaskTime": alb.ShortestTaskTime,
		"LongestTaskTime":  alb.LongestTaskTime,
	}
)

func stoh(heuristic string) alb.Heuristic {
	return heuristicMap[heuristic]
}

func GetStream(filename string) (io.Reader, error) {
	file, err := os.Open(filename)
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

// ValidateLine is a temporary hack to validate 2 conditions:
// 	(1) Paced line
//		If violated, we coerce the line to be valid and log a warning
//
//	(2) Global work < Global capacity
//		If violated, we return an error.
func ValidateLine(line *alb.Line, cycleTime float64) (float64, error) {
	validCycleTime := cycleTime

	// Coerce to paced line
	for _, task := range line.Tasks() {
		ttime := task.Time()
		if ttime > cycleTime {
			log.WithFields(log.Fields{
				"task":       task.ID,
				"task_time":  ttime,
				"cycle_time": cycleTime,
			}).Warnf("Cycle time is being bumped to task_time")
			validCycleTime = ttime
		}
	}

	// Check global work and capacity
	total := line.NStations()
	globalWork := line.TaskTime()
	globalWorkCapacity := float64(total) * validCycleTime
	if globalWork > globalWorkCapacity {
		err := fmt.Sprintf("ss=%d, t=%f, tt=%f", total, validCycleTime, globalWork)
		return validCycleTime, fmt.Errorf("validate: global work exceeds global capacity (%s)", err)
	}

	return validCycleTime, nil
}
