package alb

import (
	"errors"
	"fmt"
	"sort"
)

// Task is a physical task performed at a station on the assembly line.
type Task struct {
	ID   int
	time float64
	predecessors map[int]*Task
	assignment   *Station
}

// NewTask returns an initialized Task pointer.
func NewTask(id int, time float64) *Task {
	return &Task{
		ID:           id,
		time:         time,
		predecessors: make(map[int]*Task, 0),
	}
}

// Time returns the task's completion time
func (t *Task) Time() float64 {
	return t.time
}

// Pred returns a task's predecessor by id.
func (t *Task) Pred(id int) *Task {
	task, _ := t.predecessors[id]
	return task
}

// Preds returns an array of predecessor tasks, sorted by task ID.
func (t *Task) Preds() []*Task {
	var keys []int
	for k := range t.predecessors {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var preds []*Task
	for _, k := range keys {
		preds = append(preds, t.predecessors[k])
	}

	return preds
}

// AddPred adds a task to the task's list of predecessors.
func (t *Task) AddPred(task *Task) {
	if t.Pred(task.ID) == nil {
		t.predecessors[task.ID] = task
	}
}

// IsAssigned checks if the task has a current station assignment.
func (t *Task) IsAssigned() bool {
	return t.assignment != nil
}

func (t *Task) Assignment() *Station {
	return t.assignment
}

// Assign adds the task to a station.
func (t *Task) Assign(station *Station) error {
	if t.assignment != nil {
		return fmt.Errorf("task already assigned to station %d", t.assignment.ID)
	}

	t.assignment = station
	return nil
}

// Withdraw removes the task from its currently assigned station.
func (t *Task) Withdraw() error {
	if t.assignment == nil {
		return errors.New("task not currently assigned to a station")
	}

	t.assignment = nil
	return nil
}
