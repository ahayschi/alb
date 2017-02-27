package alb

import (
	"fmt"
	"sort"
)

// Station is a place on an assembly line where tasks are performed.
type Station struct {
	ID    int
	tasks map[int]*Task
}

// NewStation returns an initialized Station pointer.
func NewStation(id int) *Station {
	return &Station{
		ID:    id,
		tasks: make(map[int]*Task),
	}
}

// String converts the station to a string representation.
func (s *Station) String() string {
	var str string
	var tasks string
	for _, task := range s.Tasks() {
		tasks += fmt.Sprintf("%d ", task.ID)
	}

	str += fmt.Sprintf("Station %d:\tTaskTime %.2f\tTasks %s", s.ID, s.Time(), tasks)
	return str
}

// Active indicates whether the station is in use or not in use.
func (s *Station) Active() bool {
	return len(s.tasks) > 0
}

// Task returns a task by id.
func (s *Station) Task(id int) *Task {
	task, _ := s.tasks[id]
	return task
}

// Tasks returns an array of tasks assigned to the station, sorted by task ID.
func (s *Station) Tasks() []*Task {
	var keys []int
	for k := range s.tasks {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var tasks []*Task
	for _, k := range keys {
		tasks = append(tasks, s.tasks[k])
	}

	return tasks
}

// NTasks returns the number of tasks assigned to the station.
func (s *Station) NTasks() int {
	return len(s.tasks)
}

// AssignTask adds a given task to the station.
func (s *Station) AssignTask(task *Task) error {
	t := s.Task(task.ID)
	if t != nil {
		return fmt.Errorf("station already assigned task %d", task.ID)
	}

	s.tasks[task.ID] = task
	task.Assign(s)
	return nil
}

// WithdrawTask removes a task from the station.
func (s *Station) WithdrawTask(id int) error {
	t := s.Task(id)
	if t == nil {
		return fmt.Errorf("station not currently assigned task %d", id)
	}

	err := t.Withdraw()
	if err != nil {
		return err
	}

	delete(s.tasks, id)
	return nil
}

// WithdrawTasks removes all tasks from the station.
func (s *Station) WithdrawTasks() error {
	for _, task := range s.tasks {
		err := task.Withdraw()
		if err != nil {
			return err
		}
	}

	s.tasks = make(map[int]*Task)
	return nil
}

// Time returns the station time (total task time of the tasks assigned to the station).
func (s *Station) Time() float64 {
	var total float64
	for _, task := range s.tasks {
		total += task.Time()
	}
	return total
}
