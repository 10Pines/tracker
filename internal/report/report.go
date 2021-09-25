package report

import (
	"github.com/10Pines/tracker/internal/models"
	"github.com/10Pines/tracker/internal/storage"
	"log"
	"time"
)

const days = 24 * time.Hour

type TaskStatus struct {
	Task     models.Task
	JobCount int64
	Expected int64
	Ready    bool
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

func (r *Report) Got(task models.Task, jobCount int64) {
	expectedJobCount := int64(task.Datapoints - task.Tolerance)
	r.statuses = append(r.statuses, TaskStatus{
		Task:     task,
		JobCount: jobCount,
		Expected: expectedJobCount,
		Ready:    r.isReady(task),
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

func Run(db storage.DB) (Report, error) {
	now := time.Now()
	var tasks []models.Task
	err := db.AllTasksSortedByIDASC(&tasks)
	if err != nil {
		return Report{}, err
	}
	report := newReport(now)
	for _, task := range tasks {
		log.Println()
		jobCount, err := countJobs(task, now, db)
		if err != nil {
			return Report{}, err
		}
		report.Got(task, jobCount)
	}
	return report, nil
}

func countJobs(task models.Task, now time.Time, db storage.DB) (int64, error) {
	sinceOffset := time.Duration(task.Datapoints) * days
	since := now.Add(-sinceOffset)
	return db.CountJobsByTaskIDAndCreatedAfter(task.ID, since)
}
