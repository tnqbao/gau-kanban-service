package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateColumn tạo column mới cho board
func (ctrl *Controller) CreateColumn(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Column] Create column request received")

	var req CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Column] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(req.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Board not found: %s", req.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create Column] User %s does not have access to board %s", userIDStr, req.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Get max order for this board and auto-increment with proper spacing for fractional ordering
	maxOrder, err := ctrl.Repository.GetMaxOrderByBoardID(req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Failed to get max order for board: %s", req.BoardID)
		utils.JSON500(c, "Failed to calculate column order")
		return
	}

	// Create column with adaptive spacing for optimal fractional ordering
	newOrder := maxOrder + 1000 // Use 1000 spacing for new columns to ensure plenty of room
	if maxOrder == 0 {
		newOrder = 1000 // First column starts at 1000
	}

	column := &entity.Column{
		BoardID:  req.BoardID,
		Title:    req.Title,
		Position: newOrder,
		WipLimit: req.WipLimit,
	}

	if err := ctrl.Repository.CreateColumn(column); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Column] Failed to create column")
		utils.JSON500(c, "Failed to create column")
		return
	}

	// Initialize empty tickets array for response
	column.Tickets = []entity.Ticket{}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Column] Column created successfully with position %d: %s", column.Position, column.ID)
	utils.JSON200(c, gin.H{
		"message": "Column created successfully",
		"data":    column,
	})
}

// UpdateColumn cập nhật thông tin column
func (ctrl *Controller) UpdateColumn(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column] Update column request received for ID: %s", columnID)

	var req UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Column] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get column and check if it exists
	column, err := ctrl.Repository.GetColumnByID(columnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Column not found: %s", columnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Update Column] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Update column fields if provided
	if req.Title != "" {
		column.Title = req.Title
	}
	if req.WipLimit != nil {
		column.WipLimit = req.WipLimit
	}

	if err := ctrl.Repository.UpdateColumn(column); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Column] Failed to update column: %s", columnID)
		utils.JSON500(c, "Failed to update column")
		return
	}

	// Initialize empty tickets array for response
	column.Tickets = []entity.Ticket{}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Column] Column updated successfully: %s", columnID)
	utils.JSON200(c, gin.H{
		"message": "Column updated successfully",
		"data":    column,
	})
}

// DeleteColumn xóa column
func (ctrl *Controller) DeleteColumn(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Delete column request received for ID: %s", columnID)

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Column] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get column and check if it exists
	column, err := ctrl.Repository.GetColumnByID(columnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Column not found: %s", columnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Delete Column] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	if err := ctrl.Repository.DeleteColumn(columnID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Column] Failed to delete column: %s", columnID)
		utils.JSON500(c, "Failed to delete column")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Column] Column deleted successfully: %s", columnID)
	utils.JSON200(c, gin.H{
		"message": "Column deleted successfully",
	})
}

// GetColumnByID lấy thông tin column theo ID
func (ctrl *Controller) GetColumnByID(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column] Get column request received for ID: %s", columnID)

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Get Column] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get column and check if it exists
	column, err := ctrl.Repository.GetColumnWithTicketsByID(columnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Column] Column not found: %s", columnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Column] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Column] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Get Column] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Initialize empty tickets array if nil
	if column.Tickets == nil {
		column.Tickets = []entity.Ticket{}
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Column] Column retrieved successfully: %s", columnID)
	utils.JSON200(c, gin.H{
		"message": "Column retrieved successfully",
		"data":    column,
	})
}

// ReorderColumn sắp xếp lại vị trí column (Fractional Ordering - Optimal)
func (ctrl *Controller) ReorderColumn(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Param("id") // Lấy column ID từ URL parameter
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Reorder Column] Reorder column request received for column: %s", columnID)

	var req ReorderColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate position: check if trying to reorder column relative to itself
	if len(req.Position) > 6 && req.Position[:6] == "after:" {
		targetID := req.Position[6:]
		if targetID == columnID {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Reorder Column] Cannot reorder column after itself: %s", columnID)
			utils.JSON400(c, "Cannot reorder column after itself")
			return
		}
	}

	if len(req.Position) > 7 && req.Position[:7] == "before:" {
		targetID := req.Position[7:]
		if targetID == columnID {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Reorder Column] Cannot reorder column before itself: %s", columnID)
			utils.JSON400(c, "Cannot reorder column before itself")
			return
		}
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Reorder Column] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get column and check if it exists - also get board_id from column
	column, err := ctrl.Repository.GetColumnByID(columnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Column not found: %s", columnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Use board_id from the column (no need to pass in request)
	boardID := column.BoardID

	// Check if column is already in the desired position
	if req.Position == "first" {
		// Check if this column is already first (has minimum order)
		minOrder, err := ctrl.Repository.GetMinOrderByBoardID(boardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Failed to get min order")
			utils.JSON500(c, "Failed to check current position")
			return
		}
		if column.Position == minOrder {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Reorder Column] Column %s is already first", columnID)
			utils.JSON200(c, gin.H{
				"message": "Column is already in first position",
				"data": gin.H{
					"board_id":      boardID,
					"column_id":     columnID,
					"current_order": column.Position,
					"position":      req.Position,
				},
			})
			return
		}
	}

	if req.Position == "last" {
		// Check if this column is already last (has maximum order)
		maxOrder, err := ctrl.Repository.GetMaxOrderByBoardID(boardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Failed to get max order")
			utils.JSON500(c, "Failed to check current position")
			return
		}
		if column.Position == maxOrder {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Reorder Column] Column %s is already last", columnID)
			utils.JSON200(c, gin.H{
				"message": "Column is already in last position",
				"data": gin.H{
					"board_id":      boardID,
					"column_id":     columnID,
					"current_order": column.Position,
					"position":      req.Position,
				},
			})
			return
		}
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(boardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Board not found: %s", boardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Reorder Column] User %s does not have access to board %s", userIDStr, boardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Calculate new order position using fractional ordering
	newOrder, err := ctrl.Repository.GetColumnOrderPosition(boardID, req.Position)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Failed to calculate new order position")
		utils.JSON500(c, "Failed to calculate new position")
		return
	}

	// Update column order (only 1 UPDATE query - optimal!)
	if err := ctrl.Repository.UpdateColumnOrder(columnID, newOrder); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Reorder Column] Failed to update column order")
		utils.JSON500(c, "Failed to reorder column")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Reorder Column] Column reordered successfully: %s to position %s (order: %d)", columnID, req.Position, newOrder)
	utils.JSON200(c, gin.H{
		"message": "Column reordered successfully",
		"data": gin.H{
			"board_id":  boardID,
			"column_id": columnID,
			"new_order": newOrder,
			"position":  req.Position,
		},
	})
}
