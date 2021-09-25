package http

import (
	"errors"
	"github.com/10Pines/tracker/internal/logic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type createTask struct {
	Name       string `json:"name" binding:"required"`
	Datapoints int    `json:"datapoints" binding:"required" validate:"gt=0"`
	Tolerance  int    `json:"tolerance" validate:"gte=0"`
}

func NewRouter(l logic.Logic) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())

	tasks := router.Group("/api/tasks")

	tasks.POST("", func(g *gin.Context) {
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

	tasks.POST("/:taskID/jobs", func(g *gin.Context) {
		taskID, err := extractID(g, "taskID")
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		job, err := l.CreateJob(taskID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			g.JSON(http.StatusNotFound, gin.H{})
			return
		}
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusOK, gin.H{"id": job.ID})
	})
	return router
}

func extractID(g *gin.Context, paramName string) (uint, error) {
	param := g.Param(paramName)
	id, err := strconv.ParseUint(param, 10, 32)
	return uint(id), err
}
