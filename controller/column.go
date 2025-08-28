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

	// Validate that the board exists
	_, err := ctrl.Repository.GetBoardByID(req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Board not found: %s", req.BoardID)
		utils.JSON404(c, "Board not found")
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
		BoardID:  req.BoardID,
		Title:    req.Title,
		Position: maxPosition + 1,
	}

	if err := ctrl.Repository.Create(column); err != nil {
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

// GetColumns lấy tất cả columns với tickets
func (ctrl *Controller) GetColumns(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns] Get all columns with tickets request received")

	columns, err := ctrl.Repository.GetAllWithTickets()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns] Failed to get columns with tickets")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns] Successfully retrieved %d columns with tickets", len(columns))
	utils.JSON200(c, gin.H{
		"message": "Columns retrieved successfully",
		"data":    columns,
	})
}

// GetColumnById lấy column theo ID với tickets
func (ctrl *Controller) GetColumnById(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column By ID] Get column by ID request received for ID: %s", id)

	column, err := ctrl.Repository.GetColumnById(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Column By ID] Column not found: %s", id)
		utils.JSON404(c, "Column not found")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column By ID] Column retrieved successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Column retrieved successfully",
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

	column, err := ctrl.Repository.GetByID(id)
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

	if err := ctrl.Repository.Update(column); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Failed to update column: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column] Column updated successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Column updated successfully",
		"data":    column,
	})
}

// DeleteColumn xóa column
func (ctrl *Controller) DeleteColumn(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Delete column request received for ID: %s", id)

	if err := ctrl.Repository.Delete(id); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Failed to delete column: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Column deleted successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Column deleted successfully",
	})
}

// GetColumnsByBoardId lấy tất cả columns của một board theo board ID
func (ctrl *Controller) GetColumnsByBoardId(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id") // Changed from "boardId" to "id" to match new route structure
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns By Board ID] Get columns by board ID request received for board: %s", boardID)

	// Validate that the board exists
	_, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns By Board ID] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	columns, err := ctrl.Repository.GetColumnsByBoardId(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns By Board ID] Failed to get columns for board: %s", boardID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns By Board ID] Successfully retrieved %d columns for board: %s", len(columns), boardID)
	utils.JSON200(c, gin.H{
		"message": "Columns retrieved successfully",
		"data":    columns,
	})
}

// GetColumnsByBoardIdWithTickets lấy tất cả columns của một board với tickets
func (ctrl *Controller) GetColumnsByBoardIdWithTickets(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id") // Changed from "boardId" to "id" to match new route structure
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns By Board ID With Tickets] Get columns with tickets by board ID request received for board: %s", boardID)

	// Validate that the board exists
	_, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns By Board ID With Tickets] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	columns, err := ctrl.Repository.GetColumnsByBoardIdWithTickets(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Columns By Board ID With Tickets] Failed to get columns with tickets for board: %s", boardID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Columns By Board ID With Tickets] Successfully retrieved %d columns with tickets for board: %s", len(columns), boardID)
	utils.JSON200(c, gin.H{
		"message": "Columns with tickets retrieved successfully",
		"data":    columns,
	})
}

// ChangeColumnPosition thay đổi vị trí column với xử lý nâng cao
func (ctrl *Controller) ChangeColumnPosition(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Change Column Position] Change column position request received for ID: %s", columnID)

	var req ChangeColumnPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Column Position] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate new position is positive
	if req.NewPosition < 1 {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Change Column Position] Invalid position: %d", req.NewPosition)
		utils.JSON400(c, "Position must be greater than 0")
		return
	}

	// Check if column exists
	column, err := ctrl.Repository.GetByID(columnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Column Position] Column not found: %s", columnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Get max position to validate new position
	maxPosition, err := ctrl.Repository.GetMaxColumnPosition()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Column Position] Failed to get max column position")
		utils.JSON500(c, "Failed to get max column position: "+err.Error())
		return
	}

	// Validate new position doesn't exceed max
	if req.NewPosition > maxPosition {
		req.NewPosition = maxPosition // Set to last position if exceeds
	}

	// Change column position
	if err := ctrl.Repository.ChangeColumnPosition(columnID, req.NewPosition); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Column Position] Failed to change column position: %s", columnID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Change Column Position] Column position changed successfully: %s from position %d to %d", columnID, column.Position, req.NewPosition)
	utils.JSON200(c, gin.H{
		"message": "Column position changed successfully",
		"data": gin.H{
			"column_id":    columnID,
			"old_position": column.Position,
			"new_position": req.NewPosition,
		},
	})
}
