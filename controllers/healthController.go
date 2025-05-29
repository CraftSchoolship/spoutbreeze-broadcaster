package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck godoc
// @Summary      Health Check
// @Description  Returns the health status of the application
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (h *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// ReadinessCheck godoc
// @Summary      Readiness Check
// @Description  Checks if the application is ready to handle traffic
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /readiness [get]
func (h *HealthController) ReadinessCheck(c *gin.Context) {
	// Here, you can add checks for your services like DB connection, etc.
	// For example:
	// - Check DB connection
	// - Check cache
	// - Check other dependencies

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}