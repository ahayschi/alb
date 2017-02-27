package alb

import "testing"

func TestTaskTime(t *testing.T) {
	var tests = []struct {
		time float64
		want float64
	}{
		{10.0, 10.0},
		{0.0, 0.0},
		{0.5, 0.5},
	}

	for _, test := range tests {
		task := NewTask(1, test.time)
		got := task.Time()
		if got != test.want {
			t.Errorf("task.Time(%f) = %f, got %f", test.time, test.want, got)
		}
	}
}

func TestPredInTask(t *testing.T) {
	task1 := NewTask(1, 10.0)
	pred2 := NewTask(2, 10.0)
	pred3 := NewTask(3, 10.0)

	task1.AddPred(pred2)
	task1.AddPred(pred3)

	var tests = []struct {
		id   int
		want *Task
	}{
		{2, pred2},
		{3, pred3},
		{4, nil},
	}

	for _, test := range tests {
		if got := task1.Pred(test.id); got != test.want {
			t.Errorf("task.Pred(%d) = %v, got %v", test.id, test.want, got)
		}
	}
}

func TestPredsInTask(t *testing.T) {
	task1 := NewTask(1, 10.0)
	task2 := NewTask(2, 10.0)
	pred3 := NewTask(3, 10.0)
	pred4 := NewTask(4, 10.0)

	task1.AddPred(pred3)
	task1.AddPred(pred4)

	var tests = []struct {
		task *Task
		want []*Task
	}{
		{task1, []*Task{pred3, pred4}},
		{task2, []*Task{}},
	}

	for _, test := range tests {
		got := test.task.Preds()

		if len(test.want) != len(got) {
			t.Errorf("expected preds len %d != got preds len %d", len(test.want), len(got))
		}

		for _, want := range test.want {
			present := false
			for _, task := range got {
				if want == task {
					present = true
				}
			}

			if !present {
				t.Errorf("task %d not found", want.ID)
			}
		}
	}
}

func TestIsAssigned(t *testing.T) {
	station1 := NewStation(1)
	task1 := NewTask(1, 10.0)

	got := task1.IsAssigned()
	if got {
		t.Error("task.IsAssigned() = false, got true")
	}

	task1.Assign(station1)

	got = task1.IsAssigned()
	if !got {
		t.Error("task.IsAssigned() = true, got false")
	}
}

func TestAssignToStation(t *testing.T) {
	station1 := NewStation(1)
	task1 := NewTask(1, 10.0)

	got := task1.Assign(station1)
	if got != nil {
		t.Errorf("task.Assign(%d) = nil, got %s", station1.ID, got)
	}

	got = task1.Assign(station1)
	if got == nil {
		t.Errorf("task.Assign(%d) = error, got nil", station1.ID)
	}
}

func TestUnassignFromStation(t *testing.T) {
	station1 := NewStation(1)
	task1 := NewTask(1, 10.0)

	got := task1.Assign(station1)
	if got != nil {
		t.Errorf("task.Assign(%d) = nil, got %s", station1.ID, got)
	}

	got = task1.Withdraw()
	if got != nil {
		t.Errorf("task.Withdraw() = nil, got %s", got)
	}

	got = task1.Withdraw()
	if got == nil {
		t.Errorf("task.Withdraw(%d) = error, got nil", station1.ID)
	}
}
