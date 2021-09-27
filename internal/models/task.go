package models

type Task struct {
	Model
	Name       string `gorm:"uniqueIndex"`
	Tolerance  int
	Datapoints int
	Backups    []Backup
}

func (t *Task) CreateBackup() Backup {
	backup := Backup{}
	backup.TaskID = t.ID
	t.Backups = append(t.Backups, backup)
	return backup
}

func NewTask(name string, datapoints, tolerance int) Task {
	return Task{
		Name:       name,
		Tolerance:  tolerance,
		Datapoints: datapoints,
	}
}
