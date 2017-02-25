package alb

import "testing"

func TestIsActive(t *testing.T) {
	station1 := NewStation(1)
	station2 := NewStation(2)

	task1 := NewTask(1, 10.0)

	station2.AssignTask(task1)

	var tests = []struct {
		station *Station
		want    bool
	}{
		{station1, false},
		{station2, true},
	}

	for _, test := range tests {
		got := test.station.Active()
		if got != test.want {
			t.Errorf("station.Active() = %t, got %t", test.want, got)
		}
	}
}

// TODO(ah): test task in station (also assign)
//func TestTaskInStation(t *testing.T) {
//
//}

// TODO(ah): test tasks in station (also assign)
//func TestTasksInStation(t *testing.T) {
//
//}

// TODO(ah): test number of tasks in station (also assign)
//func TestNumTasksInStation(t *testing.T) {
//
//}

// TODO(ah): test unassigning task in station
//func TestUnassignTaskInStation(t *testing.T) {
//
//}
