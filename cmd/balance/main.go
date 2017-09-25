package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/parallelworks/alb"
)

func main() {
	var (
		filename  = flag.String("file", "", "input in2 file")
		cycleTime = flag.Float64("cycle", 60.0, "cycle time of line")
		heuristic = flag.String("heuristic", "LongestTaskTime", "balancing heuristic")
	)

	flag.Parse()
	log.SetLevel(log.DebugLevel)

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	input, err := GetStream(*filename)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	line := alb.NewLine(*filename)

	tasks, stations, err := ParseIn2File(input)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	line.AddTasks(tasks)
	line.AddStations(stations)

	ctime, err := ValidateLine(line, *cycleTime)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	constraints := []alb.Constraint{
		&alb.SingleTaskAssignment{},
		&alb.RestrictedStationTime{Time: ctime},
		&alb.PredecessorsStartToStart{},
	}
	line.AddConstraints(constraints)

	// Balance
	h := stoh(*heuristic)
	err = line.BalanceByStationId(h)
	if err != nil {
		log.Fatalf("balance: %s", err)
	}

	alb.PrintMeasurements(line, ctime)
	alb.PrintFreeTasks(line)
	alb.PrintStations(line)
	alb.PrintTaskVector(line)
}
