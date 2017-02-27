package alb

import "testing"

// whitebox: test add station and add task
func TestAddStationToLine(t *testing.T) {
	line := NewLine("TestAddStationToLine")
	station := NewStation(100)

	line.AddStation(station)

	got := len(line.stations)
	if got != 1 {
		t.Errorf("len(line.AddStation(%d).stations) = 1, got %d", station.ID, got)
	}

	gotS, ok := line.stations[station.ID]
	if !ok || gotS != station {
		t.Errorf("line.AddStation(%d).stations[0] = %d, got %v", station.ID, station.ID, gotS)
	}
}

func TestStationInLine(t *testing.T) {
	line := NewLine("TestStationInLine")

	station1 := NewStation(1)
	station2 := NewStation(2)
	station3 := NewStation(3)

	_ = line.AddStation(station1)
	_ = line.AddStation(station2)

	var tests = []struct {
		station *Station
		want    *Station
	}{
		{station1, station1},
		{station2, station2},
		{station3, nil},
	}

	for _, test := range tests {
		if got := line.Station(test.station.ID); got != test.want {
			if test.want != nil {
				t.Errorf("line.Station(%d) = nil, got station %d", test.station.ID, got.ID)
			} else {
				t.Errorf("line.Station(%d) = %d, got station %d", test.station.ID, test.want.ID, got.ID)
			}
		}
	}
}

func TestStationsInLine(t *testing.T) {
	line1 := NewLine("1")
	line2 := NewLine("2")

	station1 := NewStation(1)
	station2 := NewStation(2)

	stations1 := []*Station{station1, station2}

	err := line2.AddStations(stations1)
	if err != nil {
		t.Errorf("AddStations returned an error, %s", err)
	}

	var tests = []struct {
		line *Line
		want []*Station
	}{
		{line1, []*Station{}},
		{line2, stations1},
	}

	for _, test := range tests {
		got := test.line.Stations()

		if len(got) != len(test.want) {
			t.Errorf("line.Stations(%v) = %v, got %v", test.want, test.want, got)
		}

		for i, station := range got {
			if station != test.want[i] {
				t.Errorf("line.Stations(%v)[%d] = %d, got %d", test.want, i, test.want[i].ID, station.ID)
			}
		}
	}
}

// TODO(ah): test that stations returned in order by id
//func TestStationsInLineInOrder(t *testing.T) {
//}

func TestTaskInLine(t *testing.T) {
	line := NewLine("TestTaskInLine")

	task1 := NewTask(1, 10.0)
	task2 := NewTask(2, 10.0)
	task3 := NewTask(3, 10.0)

	_ = line.AddTask(task1)
	_ = line.AddTask(task2)

	var tests = []struct {
		task *Task
		want *Task
	}{
		{task1, task1},
		{task2, task2},
		{task3, nil},
	}

	for _, test := range tests {
		if got := line.Task(test.task.ID); got != test.want {
			if test.want != nil {
				t.Errorf("line.Task(%d) = nil, got task %d", test.task.ID, got.ID)
			} else {
				t.Errorf("line.Task(%d) = %d, got Task %d", test.task.ID, test.want.ID, got.ID)
			}
		}
	}
}

func TestTasksInLine(t *testing.T) {
	line1 := NewLine("1")
	line2 := NewLine("2")

	task1 := NewTask(1, 10.0)
	task2 := NewTask(2, 10.0)

	tasks1 := []*Task{task1, task2}

	_ = line2.AddTasks(tasks1)

	var tests = []struct {
		line *Line
		want []*Task
	}{
		{line1, nil},
		{line2, tasks1},
	}

	for _, test := range tests {
		got := test.line.Tasks()
		if test.want == nil && got != nil {
			t.Errorf("line.Tasks(nil) = nil, got tasks %v", got)
		}

		if len(got) != len(test.want) {
			t.Errorf("line.Tasks(%v) = %v, got %v", test.want, test.want, got)
		}

		for i, task := range got {
			if task != test.want[i] {
				t.Errorf("line.Tasks(%v)[%d] = %d, got %d", test.want, i, test.want[i].ID, task.ID)
			}
		}
	}
}

// TODO(ah): test that tasks returned in order by id
//func TestTasksInLineInOrder(t *testing.T) {
//}
