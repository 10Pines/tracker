package models

type JobStatus = int

type Job struct {
	Model
	TaskID uint
	Task   Task
}
