package logic

import (
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
)

func createTask(db *gorm.DB, task *models.Task) error {
	err := db.Create(task).Error
	return err
}

func saveBackup(db *gorm.DB, backup *models.Backup) error {
	err := db.Save(backup).Error
	return err
}

// AllTasksSortedByIDASC retrieves all tasks sorted by ID DESC
func AllTasksSortedByIDASC(db *gorm.DB, tasks *[]models.Task) error {
	err := db.Find(tasks).Order("id ASC").Error
	return err
}

func findByTaskNameOrCreate(db *gorm.DB, taskName string) (models.Task, error) {
	task := models.Task{}
	err := db.Where(&models.Task{Name: taskName}).FirstOrCreate(&task, taskDefaults(taskName)).Error
	return task, err
}

func taskDefaults(taskName string) models.Task {
	return models.Task{
		Name:       taskName,
		Tolerance:  0,
		Datapoints: 7,
	}
}

// BackupStats contains information regarding a task backups
type BackupStats struct {
	CountWithinDatapoints int64
	LastBackup            time.Time
}

// BackupsStatsByTaskID returns the backups' stats of said task after the given time
func BackupsStatsByTaskID(db *gorm.DB, taskID uint, since time.Time) (BackupStats, error) {
	var backupCount int64
	err := db.Model(&models.Backup{}).
		Where("task_id = ? AND created_at > ?", taskID, since).
		Count(&backupCount).Error
	if err != nil {
		return BackupStats{}, err
	}

	var lastBackup models.Backup
	err = db.Last(&lastBackup).Error
	if err != nil {
		return BackupStats{}, err
	}

	return BackupStats{
		LastBackup:            lastBackup.CreatedAt,
		CountWithinDatapoints: backupCount,
	}, nil
}
