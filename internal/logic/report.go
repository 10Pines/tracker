package logic

import (
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
)

const days = 24 * time.Hour

// TaskStatus represents the result of checking a task
type TaskStatus struct {
	Task        models.Task
	BackupCount int64
	Expected    int64
	Ready       bool
	LastBackup  time.Time
}

// IsOK returns whether the task is in an OK state
func (s TaskStatus) IsOK() bool {
	return !s.Ready || s.Expected <= s.BackupCount
}

// TaskHasBackups returns if the task has any backups at all
func (s TaskStatus) TaskHasBackups() bool {
	return !s.LastBackup.IsZero()
}

// Report represents the results of performing a check on all tasks
type Report struct {
	Timestamp time.Time
	statuses  []TaskStatus
}

func newReport(timestamp time.Time) Report {
	return Report{
		Timestamp: timestamp,
	}
}

// Got tracks a task with the given backup stats
func (r *Report) Got(task models.Task, stats backupStats) {
	expectedBackupCount := int64(task.Datapoints - task.Tolerance)
	status := TaskStatus{
		Task:        task,
		BackupCount: stats.CountWithinDatapoints,
		Expected:    expectedBackupCount,
		Ready:       r.isReady(task),
		LastBackup:  stats.LastBackup,
	}
	r.statuses = append(r.statuses, status)
}

func (r Report) isReady(task models.Task) bool {
	daysUntilSufficientDatapoints := time.Duration(task.Datapoints) * days
	notBefore := task.CreatedAt.Add(daysUntilSufficientDatapoints)
	return r.Timestamp.Unix() >= notBefore.Unix()
}

// TaskCount returns the report tasks status count
func (r Report) TaskCount() int {
	return len(r.statuses)
}

// Statuses returns the report tasks status
func (r Report) Statuses() []TaskStatus {
	return r.statuses
}

// IsOK returns whether all tasks are in an OK state
func (r *Report) IsOK() bool {
	for _, status := range r.statuses {
		if !status.IsOK() {
			return false
		}
	}
	return true
}

func getStats(task models.Task, now time.Time, db *gorm.DB) (backupStats, error) {
	sinceOffset := time.Duration(task.Datapoints) * days
	since := now.Add(-sinceOffset)
	return backupsStatsByTaskID(db, task.ID, since)
}
