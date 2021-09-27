package models

type Backup struct {
	Model
	TaskID uint
	Task   Task
}
