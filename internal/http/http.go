package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/10Pines/tracker/v2/pkg/tracker"

	"github.com/10Pines/tracker/v2/internal/logic"
)

type createTask struct {
	Name       string `json:"name" binding:"required"`
	Datapoints int    `json:"datapoints" binding:"required" validate:"gt=0"`
	Tolerance  int    `json:"tolerance" validate:"gte=0"`
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
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusCreated, gin.H{"id": task.ID})
	})
	return router
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
