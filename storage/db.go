package storage

import (
	"github.com/10Pines/tracker/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"time"
)

type DB interface {
	CreateTask(task *models.Task) error
	FindTaskByID(taskID uint, task *models.Task) error
	SaveTask(task *models.Task) error
	AllTasksSortedByIDASC(tasks *[]models.Task) error
	CountJobsByTaskIDAndCreatedAfter(taskID uint, since time.Time) (int64, error)
}

type SQL struct {
	*gorm.DB
}

func DBPath() string {
	dbPath := os.Getenv("DB_PATH")
	return path.Join(dbPath, "tracker.db")
}

func MustNewSQL(path string) DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	err = db.AutoMigrate(&models.Job{}, &models.Task{})
	if err != nil {
		log.Fatal(err)
	}
	return &SQL{
		db,
	}
}

func (db *SQL) CreateTask(task *models.Task) error {
	err := db.Create(task).Error
	return err
}

func (db *SQL) FindTaskByID(taskID uint, task *models.Task) error {
	return db.First(task, taskID).Error
}

func (db *SQL) SaveTask(task *models.Task) error {
	return db.Save(task).Error
}

func (db *SQL) AllTasksSortedByIDASC(tasks *[]models.Task) error {
	return db.Find(tasks).Order("id ASC").Error
}

func (db *SQL) CountJobsByTaskIDAndCreatedAfter(taskID uint, since time.Time) (int64, error) {
	var jobCount int64
	err := db.Model(&models.Job{}).
		Where("task_id = ? AND created_at > ?", taskID, since).
		Count(&jobCount).Error
	return jobCount, err
}
