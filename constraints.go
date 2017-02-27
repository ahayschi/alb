package alb

type Constraint interface {
	Valid(*Task, *Station) bool
}

type OnlyActiveStations struct {
}

func (c *OnlyActiveStations) Valid(task *Task, station *Station) bool {
	return station.Active()
}

type SingleTaskAssignment struct {
}

func (c *SingleTaskAssignment) Valid(task *Task, station *Station) bool {
	return !task.IsAssigned()
}

type RestrictedStationTime struct {
	Time float64
}

func (c *RestrictedStationTime) Valid(task *Task, station *Station) bool {
	return task.Time()+station.Time() <= c.Time
}

type PredecessorsStartToStart struct {
}

func (c *PredecessorsStartToStart) Valid(task *Task, station *Station) bool {
	preds := task.Preds()
	for _, pred := range preds {
		if !pred.IsAssigned() {
			return false
		}
	}

	return true
}
