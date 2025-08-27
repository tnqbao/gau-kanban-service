package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateColumn tạo column mới
func (ctrl *Controller) CreateColumn(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Column] Create new column request received")

	var req CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Get the current max position and set the new column's position to max + 1
	maxPosition, err := ctrl.Repository.GetMaxColumnPosition()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Failed to get max column position")
		utils.JSON500(c, "Failed to get max column position: "+err.Error())
		return
	}

	column := &entity.Column{
		Title:    req.Title,
		Position: maxPosition + 1,
	}

	if err := ctrl.Repository.CreateColumn(column); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Failed to create column")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Column] Column created successfully: %s", column.ID)
	utils.JSON200(c, gin.H{
		"message": "Column created successfully",
		"data":    column,
	})
}

// UpdateColumn cập nhật thông tin column
func (ctrl *Controller) UpdateColumn(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column] Update column request received for ID: %s", id)

	var req UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	column, err := ctrl.Repository.GetColumnByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Column not found: %s", id)
		utils.JSON404(c, "Column not found")
		return
	}

	if req.Title != "" {
		column.Title = req.Title
	}
	if req.Position != nil {
		column.Position = *req.Position
	}

	if err := ctrl.Repository.UpdateColumn(column); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Failed to update column: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column] Column updated successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Column updated suctrlessfully",
		"data":    column,
	})
}

// DeleteColumn xóa column
func (ctrl *Controller) DeleteColumn(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Delete column request received for ID: %s", id)

	if err := ctrl.Repository.DeleteColumn(id); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Failed to delete column: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Column deleted successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Column deleted suctrlessfully",
	})
}

// GetColumns lấy danh sách tất cả columns với tickets
func (ctrl *Controller) GetColumns(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns] Get all columns request received")

	columns, err := ctrl.Repository.GetAllColumnWithTickets()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns] Failed to get columns")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns] Retrieved %d columns successfully", len(columns))
	utils.JSON200(c, gin.H{
		"data": columns,
	})
}

// GetColumnById lấy thông tin column theo ID với tickets
func (ctrl *Controller) GetColumnById(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column By ID] Get column request received for ID: %s", id)

	column, err := ctrl.Repository.GetColumnByIdWithTickets(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Column By ID] Column not found: %s", id)
		utils.JSON404(c, "Column not found")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column By ID] Column retrieved successfully: %s", id)
	utils.JSON200(c, gin.H{
		"data": column,
	})
}

// UpdateColumnPosition thay đổi vị trí column
func (ctrl *Controller) UpdateColumnPosition(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column Position] Update column position request received for ID: %s", id)

	var req UpdateColumnPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column Position] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.UpdateColumnPosition(id, req.Position); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column Position] Failed to update column position: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column Position] Column position updated successfully: %s to position %d", id, req.Position)
	utils.JSON200(c, gin.H{
		"message": "Column position updated suctrlessfully",
	})
}
