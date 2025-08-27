package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// GetTagColors trả về mapping màu sắc cho các tag
func (ctrl *Controller) GetTagColors(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Tag Colors] Get tag colors request received")

	labels, err := ctrl.Repository.GetAllLabel()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Tag Colors] Failed to get labels")
		utils.JSON500(c, err.Error())
		return
	}

	tagColors := make(map[string]string)
	for _, label := range labels {
		tagColors[label.Name] = label.Color
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Tag Colors] Retrieved %d tag colors successfully", len(tagColors))
	utils.JSON200(c, gin.H{
		"data": tagColors,
	})
}
