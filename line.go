package alb

import (
	"fmt"
	"reflect"
	"sort"
)

// Line is an assembly line with stations and tasks.
type Line struct {
	Name        string
	stations    map[int]*Station
	tasks       map[int]*Task
	constraints []Constraint
}

// NewLine returns an initialized Line pointer.
func NewLine(name string) *Line {
	return &Line{
		Name:     name,
		stations: make(map[int]*Station),
		tasks:    make(map[int]*Task),
	}
}

// AddConstraint adds a constraint to balance the line.
func (l *Line) AddConstraint(c Constraint) {
	l.constraints = append(l.constraints, c)
}

// AddConstraints adds multiple constraints to balance the line.
func (l *Line) AddConstraints(cs []Constraint) {
	l.constraints = append(l.constraints, cs...)
}

// ReplaceConstraint replaces an existing constraint on the line by its type.
func (l *Line) ReplaceConstraint(c Constraint) error {
	dirty := false
	for i, constraint := range l.constraints {
		if reflect.TypeOf(c) == reflect.TypeOf(constraint) {
			l.constraints[i] = c
			dirty = true
			break
		}
	}

	if !dirty {
		return fmt.Errorf("no existing constraint of type %T found to replace", c)
	}

	return nil
}

// RemoveConstraints removes all constraints from the line.
func (l *Line) RemoveConstraints() {
	l.constraints = nil
}

// Station returns a station belonging to the line by id.
func (l *Line) Station(id int) *Station {
	station, _ := l.stations[id]
	return station
}

// Stations returns an array of stations belonging to the line, sorted by station ID.
func (l *Line) Stations() []*Station {
	var keys []int
	for k := range l.stations {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var stations []*Station
	for _, k := range keys {
		stations = append(stations, l.stations[k])
	}

	return stations
}

// AddStation adds a station to the line.
func (l *Line) AddStation(station *Station) error {
	s := l.Station(station.ID)
	if s != nil {
		return fmt.Errorf("line already has station %d", station.ID)
	}

	l.stations[station.ID] = station
	return nil
}

// AddStations adds stations to the line.
func (l *Line) AddStations(stations []*Station) error {
	for _, station := range stations {
		err := l.AddStation(station)
		if err != nil {
			return err
		}
	}

	return nil
}

// ActiveStations returns all active stations on the line.
func (l *Line) ActiveStations() []*Station {
	var stations []*Station
	for _, station := range l.stations {
		if station.Active() {
			stations = append(stations, station)
		}
	}
	return stations
}

// NActiveStations returns the number of active stations on the line.
func (l *Line) NActiveStations() int {
	var n int
	for _, station := range l.stations {
		if station.Active() {
			n++
		}
	}
	return n
}

// Task returns a task belonging to the line by id.
func (l *Line) Task(id int) *Task {
	task, _ := l.tasks[id]
	return task
}

// Tasks returns an array of tasks belonging to the line, sorted by task ID.
func (l *Line) Tasks() []*Task {
	var keys []int
	for k := range l.tasks {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var tasks []*Task
	for _, k := range keys {
		tasks = append(tasks, l.tasks[k])
	}

	return tasks
}

// FreeTasks returns all tasks without an assignment.
func (l *Line) FreeTasks() []*Task {
	var tasks []*Task
	for _, task := range l.Tasks() {
		if !task.IsAssigned() {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// AssignedTasks returns all tasks with an assignment.
func (l *Line) AssignedTasks() []*Task {
	var tasks []*Task
	for _, task := range l.Tasks() {
		if task.IsAssigned() {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// AddTask adds a task to the line.
func (l *Line) AddTask(task *Task) error {
	t := l.Task(task.ID)
	if t != nil {
		return fmt.Errorf("line already has task %d", task.ID)
	}

	l.tasks[task.ID] = task
	return nil
}

// AddTasks adds stations to the line.
func (l *Line) AddTasks(tasks []*Task) error {
	for _, task := range tasks {
		err := l.AddTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}

// TaskTime calculates the total task time over all tasks on the line.
func (l *Line) TaskTime() float64 {
	var total float64
	for _, task := range l.tasks {
		total += task.Time()
	}
	return total
}

// StationTime calculates the total station time for all stations on the line.
func (l *Line) StationTime() float64 {
	var total float64
	for _, station := range l.stations {
		total += station.Time()
	}
	return total
}

// UnassignTasks unassigns all tasks from all stations on the line.
func (l *Line) UnassignTasks() error {
	for _, station := range l.stations {
		err := station.WithdrawTasks()
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidAssignment checks to see if the given task can be assigned to the
// given station. Validity is dependent on the line's current constraints.
func (l *Line) ValidAssignment(taskID, stationID int) bool {
	task := l.Task(taskID)
	station := l.Station(stationID)

	if task == nil || station == nil {
		return false
	}

	for _, constraint := range l.constraints {
		if !constraint.Valid(task, station) {
			return false
		}
	}

	return true
}

// ValidAssignments returns all tasks on the line that can be assigned
// to the given station. It calls ValidAssignment for each task to
// check that the assignment would not violate any of the line's constraints.
func (l *Line) ValidAssignments(stationID int) []*Task {
	var tasks []*Task
	for _, task := range l.tasks {
		if l.ValidAssignment(task.ID, stationID) {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// Balance assigns tasks to stations on the line until all tasks have a
// valid assignment. It uses the given heuristic function for deciding
// the order in which tasks are assigned.
func (l *Line) Balance(fn Heuristic) error {
	for _, station := range l.Stations() {
		candidates := l.ValidAssignments(station.ID)
		for len(candidates) > 0 {
			best := fn(candidates)
			err := station.AssignTask(best)
			if err != nil {
				return err
			}

			candidates = l.ValidAssignments(station.ID)
		}
	}

	return nil
}
