package models

// Backup represents a successful backup
type Backup struct {
	Model
	TaskID uint
	Task   Task
}
