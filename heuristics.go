package alb

type Heuristic func([]*Task) *Task

func ShortestTaskTime(tasks []*Task) *Task {
	var min *Task
	for _, task := range tasks {
		if min == nil {
			min = task
			continue
		}

		if task.Time() < min.Time() {
			min = task
		}
	}
	return min
}

// Preferred
func LongestTaskTime(tasks []*Task) *Task {
	var max *Task
	for _, task := range tasks {
		if max == nil {
			max = task
			continue
		}

		if task.Time() > max.Time() {
			max = task
		}
	}
	return max
}

// NOTE: Experimental, not tested
func TimeOfSuccessors(task *Task, tasks []*Task) float64 {
	total := task.Time()
	for _, t := range tasks {
		if t.Pred(task.ID) != nil {
			total += TimeOfSuccessors(t, tasks)
		}
	}
	return total
}

// NOTE: Experimental, not tested
func LongestSuccesssorTime(tasks []*Task) *Task {
	var maxLength float64
	var maxTask *Task
	for _, task := range tasks {
		if maxTask == nil {
			maxTask = task
			continue
		}

		length := TimeOfSuccessors(task, tasks)
		if length > maxLength {
			maxTask = task
		}
	}
	return maxTask
}

// NOTE: Experimental, not tested
func ShortestSuccessorTime(tasks []*Task) *Task {
	var minLength float64
	var minTask *Task
	for _, task := range tasks {
		if minTask == nil {
			minTask = task
			continue
		}

		length := TimeOfSuccessors(task, tasks)
		if length < minLength {
			minTask = task
		}
	}
	return minTask
}

func NSuccessors(task *Task, tasks []*Task) int {
	total := 1
	for _, t := range tasks {
		if t.Pred(task.ID) != nil {
			total += NSuccessors(t, tasks)
		}
	}
	return total
}

// NOTE: Experimental, not tested
func MostSuccessors(tasks []*Task) *Task {
	var maxLength int
	var maxTask *Task
	for _, task := range tasks {
		if maxTask == nil {
			maxTask = task
			continue
		}

		length := NSuccessors(task, tasks)
		if length > maxLength {
			maxTask = task
		}
	}
	return maxTask
}

// NOTE: Experimental, not tested
func LeastSuccessors(tasks []*Task) *Task {
	var minLength int
	var minTask *Task
	for _, task := range tasks {
		if minTask == nil {
			minTask = task
			continue
		}

		length := NSuccessors(task, tasks)
		if length < minLength {
			minTask = task
		}
	}
	return minTask
}
