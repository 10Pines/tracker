package logic

import (
	"github.com/10Pines/tracker/internal/models"
	"github.com/10Pines/tracker/internal/storage"
)

type Logic struct {
	db storage.DB
}

type CreateTask struct {
	Name       string
	Tolerance  int
	Datapoints int
}

func New(db storage.DB) Logic {
	return Logic{
		db,
	}
}

func (o Logic) CreateTask(create CreateTask) (models.Task, error) {
	task := models.NewTask(create.Name, create.Datapoints, create.Tolerance)
	err := o.db.CreateTask(&task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (o Logic) CreateJob(taskID uint) (models.Job, error) {
	var task models.Task
	err := o.db.FindTaskByID(taskID, &task)
	if err != nil {
		return models.Job{}, err
	}
	job := task.NewJob()
	err = o.db.SaveTask(&task)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}
