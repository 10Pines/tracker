package report

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/logic"
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
func (r *Report) Got(task models.Task, stats logic.BackupStats) {
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

//// Timestamp returns the report creation Timestamp
//func (r *Report) Timestamp() time.Time {
//	return r.Timestamp
//}

// IsOK returns whether all tasks are in an OK state
func (r *Report) IsOK() bool {
	for _, status := range r.statuses {
		if !status.IsOK() {
			return false
		}
	}
	return true
}

// Run generates a report containing every task status
func Run(db *gorm.DB) (Report, error) {
	now := time.Now()
	var tasks []models.Task
	err := logic.AllTasksSortedByIDASC(db, &tasks)
	if err != nil {
		return Report{}, err
	}
	report := newReport(now)
	for _, task := range tasks {
		log.Println()
		backupStats, err := backupStats(task, now, db)
		if err != nil {
			return Report{}, err
		}
		report.Got(task, backupStats)
	}
	return report, nil
}

func backupStats(task models.Task, now time.Time, db *gorm.DB) (logic.BackupStats, error) {
	sinceOffset := time.Duration(task.Datapoints) * days
	since := now.Add(-sinceOffset)
	return logic.BackupsStatsByTaskID(db, task.ID, since)
}
