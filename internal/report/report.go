package report

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/logic"
	"github.com/10Pines/tracker/v2/internal/models"
)

const days = 24 * time.Hour

type TaskStatus struct {
	Task        models.Task
	BackupCount int64
	Expected    int64
	Ready       bool
}

type Report struct {
	timestamp time.Time
	statuses  []TaskStatus
}

func newReport(timestamp time.Time) Report {
	return Report{
		timestamp: timestamp,
	}
}

func (r *Report) Got(task models.Task, backupCount int64) {
	expectedBackupCount := int64(task.Datapoints - task.Tolerance)
	r.statuses = append(r.statuses, TaskStatus{
		Task:        task,
		BackupCount: backupCount,
		Expected:    expectedBackupCount,
		Ready:       r.isReady(task),
	})
}

func (r Report) isReady(task models.Task) bool {
	daysUntilSufficientDatapoints := time.Duration(task.Datapoints) * days
	notBefore := task.CreatedAt.Add(daysUntilSufficientDatapoints)
	return r.timestamp.Unix() >= notBefore.Unix()
}

func (r Report) TaskCount() int {
	return len(r.statuses)
}

func (r Report) Statuses() []TaskStatus {
	return r.statuses
}

func (r *Report) Timestamp() time.Time {
	return r.timestamp
}

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
		backupCount, err := countBackups(task, now, db)
		if err != nil {
			return Report{}, err
		}
		report.Got(task, backupCount)
	}
	return report, nil
}

func countBackups(task models.Task, now time.Time, db *gorm.DB) (int64, error) {
	sinceOffset := time.Duration(task.Datapoints) * days
	since := now.Add(-sinceOffset)
	return logic.CountBackupsByTaskIDAndCreatedAfter(db, task.ID, since)
}
