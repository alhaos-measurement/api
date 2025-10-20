package controller

import (
	"github.com/alhaos-measurement/api/internal/logger"
	"github.com/alhaos-measurement/api/internal/model"
	"github.com/alhaos-measurement/api/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	repo   *repository.Repository
	logger *logger.Logger
}

// New construct Controller struct
func New(repo *repository.Repository, logger *logger.Logger) *Controller {
	return &Controller{repo: repo, logger: logger}
}

// RegisterRoutes register api routes at router
func (c *Controller) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		measure := api.Group("/measure")
		{
			measure.GET("/", c.MeasureGetHandler)
			measure.POST("/", c.MeasurePostHandler)
		}
	}

}

// MeasurePostHandler read measure json
func (c *Controller) MeasurePostHandler(ctx *gin.Context) {
	var m model.Measure
	err := ctx.ShouldBindJSON(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.repo.AddMeasure(&m)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, m)
}

func (c *Controller) MeasureGetHandler(context *gin.Context) {
	var req model.LastSensorMeasure
	err := context.ShouldBindQuery(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.logger.Error("zero sensor id provided: " + err.Error())
		return
	}

	if req.SensorID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "sensor_id is required"})
		c.logger.Error("zero sensor id provided")
		return
	}

	measure, err := c.repo.GetLastMeasure(req.SensorID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, measure)
}
