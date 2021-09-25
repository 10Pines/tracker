package models

type Task struct {
	Model
	Name       string
	Tolerance  int
	Datapoints int
	Jobs       []Job
}

func (t *Task) NewJob() Job {
	job := Job{}
	t.Jobs = append(t.Jobs, job)
	return job
}

func NewTask(name string, datapoints, tolerance int) Task {
	return Task{
		Name:       name,
		Tolerance:  tolerance,
		Datapoints: datapoints,
	}
}
