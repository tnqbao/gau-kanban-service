package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateChecklist creates a new checklist item for a ticket
func (ctrl *Controller) CreateChecklist(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Checklist] Create new checklist request received")

	var req CreateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Checklist] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Checklist] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Checklist] User ID: %s, Ticket ID: %s", userIDStr, req.TicketID)

	// Verify ticket exists
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Checklist] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get next order for this ticket
	nextOrder, err := ctrl.Repository.GetNextChecklistOrder(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Checklist] Failed to get next order")
		utils.JSON500(c, "Failed to get next order")
		return
	}

	// Create checklist item
	checklist := &entity.Checklist{
		TicketID: req.TicketID,
		Title:    req.Title,
		Order:    nextOrder,
		Status:   "pending",
	}

	if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Checklist] Failed to create checklist")
		utils.JSON500(c, "Failed to create checklist")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Checklist] Checklist created successfully: %s", checklist.ID)
	utils.JSON201(c, gin.H{
		"message": "Checklist created successfully",
		"data":    checklist,
	})
}

// UpdateChecklist updates an existing checklist item
func (ctrl *Controller) UpdateChecklist(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Checklist] Update checklist request received")

	checklistID := c.Param("id")
	if checklistID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Checklist] Checklist ID is required")
		utils.JSON400(c, "Checklist ID is required")
		return
	}

	var req UpdateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Checklist] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get checklist and verify it exists
	checklist, err := ctrl.Repository.GetChecklistByID(checklistID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Checklist not found: %s", checklistID)
		utils.JSON404(c, "Checklist not found")
		return
	}

	// Get ticket to check board access
	ticket, err := ctrl.Repository.GetTicketByID(checklist.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Ticket not found: %s", checklist.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to check board access
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Update Checklist] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Validate status values
	validStatuses := map[string]bool{"pending": true, "completed": true}
	if !validStatuses[req.Status] {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Checklist] Invalid status: %s", req.Status)
		utils.JSON400(c, "Invalid status. Must be one of: pending, completed")
		return
	}

	// Update checklist status
	checklist.Status = req.Status

	if err := ctrl.Repository.UpdateChecklist(checklist); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Checklist] Failed to update checklist")
		utils.JSON500(c, "Failed to update checklist")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Checklist] Checklist updated successfully: %s", checklist.ID)
	utils.JSON200(c, gin.H{
		"message": "Checklist updated successfully",
		"data":    checklist,
	})
}

// DeleteChecklist deletes a checklist item
func (ctrl *Controller) DeleteChecklist(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Checklist] Delete checklist request received")

	checklistID := c.Param("id")
	if checklistID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Checklist] Checklist ID is required")
		utils.JSON400(c, "Checklist ID is required")
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Checklist] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get checklist and verify it exists
	checklist, err := ctrl.Repository.GetChecklistByID(checklistID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Checklist not found: %s", checklistID)
		utils.JSON404(c, "Checklist not found")
		return
	}

	// Get ticket to check board access
	ticket, err := ctrl.Repository.GetTicketByID(checklist.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Ticket not found: %s", checklist.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to check board access
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Delete Checklist] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Delete checklist from database
	if err := ctrl.Repository.DeleteChecklist(checklistID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Checklist] Failed to delete checklist")
		utils.JSON500(c, "Failed to delete checklist")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Checklist] Checklist deleted successfully: %s", checklist.ID)
	utils.JSON200(c, gin.H{
		"message": "Checklist deleted successfully",
		"data": gin.H{
			"id":        checklist.ID,
			"ticket_id": checklist.TicketID,
			"title":     checklist.Title,
		},
	})
}
