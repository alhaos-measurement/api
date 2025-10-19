package controller

import (
	"github.com/alhaos-measurement/api/internal/model"
	"github.com/alhaos-measurement/api/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	repo *repository.Repository
}

// New construct Controller struct
func New(repo *repository.Repository) *Controller {
	return &Controller{repo: repo}
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

	// get request struct instance
	var req model.LastSensorMeasure
	err := context.ShouldBindQuery(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	measure, err := c.repo.GetLastMeasure(req.SensorID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, measure)
}
