package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/parallelworks/alb"
)

func main() {
	params, err := GetParams()
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	input, err := GetStream(params)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	line := alb.NewLine(params.InputFile)

	tasks, stations, err := ParseIn2File(input)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	line.AddTasks(tasks)
	line.AddStations(stations)

	constraints := []alb.Constraint{
		&alb.SingleTaskAssignment{},
		&alb.RestrictedTaskTime{Time: params.CycleTime},
		&alb.PredecessorsStartToStart{},
	}
	line.AddConstraints(constraints)

	err = line.Balance(alb.LongestTaskTime)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	alb.PrintMeasurements(line, params.CycleTime)
	alb.PrintStations(line)

	//var tasks string
	//
	//for _, task := range line.Tasks() {
	//	tasks += fmt.Sprintf("%d ", task.Assignment().ID)
	//}
	//fmt.Println(tasks)
}
