package alb

import (
	"fmt"
	"math"
)

// Station is a place on an assembly line where tasks are performed.
type Station struct {
	ID     int
	tasks  []*Task
	active bool
}

// NewStation returns an initialized Station pointer.
func NewStation(id int) *Station {
	return &Station{
		ID: id,
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
	//return len(s.tasks) > 0
	return s.active
}

func (s *Station) Activate() {
	s.active = true
}

func (s *Station) Disable() {
	s.active = false
}

// Task returns a task by id.
func (s *Station) Task(id int) *Task {
	for _, task := range s.tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// Tasks returns an array of tasks assigned to the station, in order that they were assigned.
func (s *Station) Tasks() []*Task {
	return s.tasks
}

// Load returns the set of tasks assigned to the station.
func (s *Station) Load() []*Task {
	return s.Tasks()
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

	s.tasks = append(s.tasks, task)
	task.Assign(s)
	return nil
}

// WithdrawTask removes a task from the station.
func (s *Station) WithdrawTask(id int) error {
	tasks := s.tasks[:0]
	for _, task := range s.tasks {
		if task.ID != id {
			tasks = append(tasks, task)
			continue
		}

		err := task.Withdraw()
		if err != nil {
			return err
		}
	}

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

	s.tasks = make([]*Task, 0)
	return nil
}

// Time returns the station time (total task time of the tasks assigned
// to the station).
func (s *Station) Time() float64 {
	var total float64
	for _, task := range s.tasks {
		total += task.Time()
	}
	return total
}

// IdleTime returns the absolute difference between the given cycle time
// and the station time.
func (s *Station) IdleTime(time float64) float64 {
	return math.Abs(time - s.Time())
}
