package logic

import (
	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
)

// Logic contains all business logic
type Logic struct {
	db *gorm.DB
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
func New(db *gorm.DB) Logic {
	return Logic{
		db,
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
