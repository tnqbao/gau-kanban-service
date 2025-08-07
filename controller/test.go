package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/utils"
)

func (ctrl *Controller) TestDeployment(c *gin.Context) {
	utils.JSON200(c, gin.H{"message": "Deployment test successful"})
}

func (ctrl *Controller) CheckHealth(c *gin.Context) {
	utils.JSON200(c, gin.H{"status": "running"})
}
