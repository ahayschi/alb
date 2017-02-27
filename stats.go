package alb

import (
	"fmt"
	"math"
)

func Efficiency(line *Line, time float64) float64 {
	ttime := line.TaskTime()
	return ttime / (time * float64(line.NActiveStations())) * 100
}

func SmoothnessIndex(line *Line, time float64) float64 {
	var idx float64
	for _, station := range line.stations {
		if !station.Active() {
			continue
		}
		ttime := station.Time()
		idx += math.Pow(time-ttime, 2)
	}
	return math.Sqrt(idx)
}

func PrintMeasurements(line *Line, time float64) {
	fmt.Printf("%s\n", line.Name)
	fmt.Printf("cycle_time=%.2f\n", time)
	fmt.Printf("theoretical_min=%d\n", int(math.Ceil(line.TaskTime()/time)))
	fmt.Printf("measured_min=%d\n", line.NActiveStations())
	fmt.Printf("line_efficiency=%.1f%%\n", Efficiency(line, time))
	fmt.Printf("smoothness_index=%.1f\n", SmoothnessIndex(line, time))
}

func PrintStations(line *Line) {
	for _, station := range line.Stations() {
		if station.Active() {
			fmt.Println(station)
		}
	}
}

func PrintFreeTasks(line *Line) {
	var tasks string
	for _, task := range line.FreeTasks() {
		tasks += fmt.Sprintf("%d ", task.ID)
	}
	fmt.Printf("free_tasks=%s\n", tasks)
}

func PrintTaskVector(line *Line) {
	var tasks string
	for _, task := range line.AssignedTasks() {
		tasks += fmt.Sprintf("%d ", task.Assignment().ID)
	}
	fmt.Printf("%s\n", tasks)
}
