package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// GetTagColors trả về mapping màu sắc cho các tag
func (ctrl *Controller) GetTagColors(c *gin.Context) {
	labels, err := ctrl.Repository.GetAllLabel()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	tagColors := make(map[string]string)
	for _, label := range labels {
		tagColors[label.Name] = label.Color
	}

	utils.JSON200(c, gin.H{
		"data": tagColors,
	})
}
