package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/10Pines/tracker/v2/internal/shared"
	"github.com/10Pines/tracker/v2/pkg/tracker"

	"github.com/10Pines/tracker/v2/internal/logic"
)

type createTask struct {
	Name       string `json:"name" binding:"required"`
	Datapoints int    `json:"datapoints" binding:"required" validate:"gt=0"`
	Tolerance  int    `json:"tolerance" validate:"gte=0"`
}

type notifyCriteria string

const (
	always  notifyCriteria = "ALWAYS"
	onError notifyCriteria = "ON_ERROR"
	never   notifyCriteria = "NEVER"
)

type createReport struct {
	Notify notifyCriteria `json:"notify"`
}

// NewRouter returns a configured router with all application routes
func NewRouter(l logic.Logic, key string) *gin.Engine {
	router := gin.New()
	router.Use(requestLogger())

	router.GET("/healthz/ready", func(g *gin.Context) {
		g.Status(200)
	})

	api := router.Group("/api")
	api.Use(apiKeyRequired(key))

	api.POST("/backups", func(g *gin.Context) {
		var params tracker.CreateBackup
		if err := g.BindJSON(&params); err != nil {
			return
		}
		create := logic.CreateBackup{
			TaskName: params.TaskName,
		}
		backup, err := l.CreateBackup(create)
		if err != nil {
			respondInternalServerError(g, err)
			return
		}
		g.JSON(http.StatusCreated, gin.H{"id": backup.ID})
	})

	api.POST("/tasks", func(g *gin.Context) {
		var params createTask
		if err := g.BindJSON(&params); err != nil {
			return
		}
		create := logic.CreateTask{
			Name:       params.Name,
			Tolerance:  params.Tolerance,
			Datapoints: params.Datapoints,
		}
		task, err := l.CreateTask(create)
		if err != nil {
			respondInternalServerError(g, err)
			return
		}
		g.JSON(http.StatusCreated, gin.H{"id": task.ID})
	})

	api.POST("/reports", func(g *gin.Context) {
		var params createReport
		if err := g.BindJSON(&params); err != nil {
			return
		}
		if params.Notify == "" {
			params.Notify = never
		}
		now := time.Now()
		report, err := l.CreateReport(now)
		if err != nil {
			respondInternalServerError(g, err)
			return
		}
		if params.Notify == always || params.Notify == onError && !report.IsOK() {
			if err = l.NotifyReport(report); err != nil {
				respondInternalServerError(g, err)
				return
			}
		}
		if params.Notify == never {
			log.Println("Skipping notification")
		}
		g.JSON(http.StatusOK, asJSON(report))
	})
	return router
}

func asJSON(report shared.Report) gin.H {
	var status string
	if report.IsOK() {
		status = "OK"
	} else {
		status = "ERR"
	}

	tasks := make([]gin.H, 0)
	for _, taskStatus := range report.Statuses() {
		task := gin.H{
			"name":    taskStatus.Task.Name,
			"status":  taskStatus.BackupCount,
			"isReady": taskStatus.Ready,
		}

		if !taskStatus.LastBackup.IsZero() {
			task["lastBackup"] = taskStatus.LastBackup.Format(time.RFC3339)
		}

		tasks = append(tasks, task)
	}

	return gin.H{
		"time":   report.Timestamp.Format(time.RFC3339),
		"status": status,
		"tasks":  tasks,
	}
}

func requestLogger() gin.HandlerFunc {
	loggerConfig := gin.LoggerConfig{
		SkipPaths: []string{
			"/healthz/ready",
		},
	}
	requestLogger := gin.LoggerWithConfig(loggerConfig)
	return requestLogger
}

func respondInternalServerError(g *gin.Context, err error) {
	g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
