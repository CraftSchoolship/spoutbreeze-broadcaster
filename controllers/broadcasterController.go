package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"spoutbreeze/models"
	"spoutbreeze/services"
)


// JoinBBB godoc
// @Summary      Join BBB
// @Description  Join a BigBlueButton session
// @Tags         Broadcaster
// @Accept       json
// @Produce      json
// @Param        request body models.BroadcasterRequest true "Broadcaster Request"
// @Success      200 {object} models.BroadcasterResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /broadcaster/joinBBB [post]
func JoinBBB(c *gin.Context) {
	var request models.BroadcasterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.ProcessBroadcasterRequest(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Broadcasting session started successfully"})
}
