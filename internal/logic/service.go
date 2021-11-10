package logic

import (
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
	"github.com/10Pines/tracker/v2/internal/reporter"
	"github.com/10Pines/tracker/v2/internal/shared"
)

// Logic contains all business logic
type Logic struct {
	db       *gorm.DB
	reporter reporter.Reporter
}

// CreateTask represents the params necessary to create a new task
type CreateTask struct {
	Name       string
	Tolerance  int
	Datapoints int
}

// CreateBackup represents the params necessary to create a new backup
type CreateBackup struct {
	TaskName string
}

// New returns a Logic instance
func New(db *gorm.DB, reporter reporter.Reporter) Logic {
	return Logic{
		db:       db,
		reporter: reporter,
	}
}

// CreateTask saves a new task
func (l Logic) CreateTask(create CreateTask) (models.Task, error) {
	task := models.NewTask(create.Name, create.Datapoints, create.Tolerance)
	err := createTask(l.db, &task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// CreateBackup saves a new backup, creating the task if missing
func (l Logic) CreateBackup(create CreateBackup) (models.Backup, error) {
	var backup models.Backup
	err := l.db.Transaction(func(tx *gorm.DB) error {
		task, err := findByTaskNameOrCreate(tx, create.TaskName)
		if err != nil {
			return err
		}
		backup = task.CreateBackup()
		err = saveBackup(tx, &backup)
		return err
	})
	if err != nil {
		return models.Backup{}, err
	}
	return backup, nil
}

// CreateReport creates a report analyzing all tasks
func (l Logic) CreateReport(now time.Time) (shared.Report, error) {
	var tasks []models.Task
	err := allTasksSortedByIDASC(l.db, &tasks)
	if err != nil {
		return shared.Report{}, err
	}
	report := shared.NewReport(now)
	for _, task := range tasks {
		stats, err := getStats(task, now, l.db)
		if err != nil {
			return shared.Report{}, err
		}
		report.Got(task, stats)
	}
	return report, nil
}

// NotifyReport broadcast the given report to configured sinks
func (l Logic) NotifyReport(report shared.Report) error {
	return l.reporter.Process(report)
}

func getStats(task models.Task, now time.Time, db *gorm.DB) (shared.BackupStats, error) {
	sinceOffset := time.Duration(task.Datapoints) * shared.Day
	since := now.Add(-sinceOffset)
	return backupsStatsByTaskID(db, task.ID, since)
}
