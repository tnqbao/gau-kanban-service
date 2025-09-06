package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateAssignee assigns a member to a ticket
func (ctrl *Controller) CreateAssignee(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignee] Create new assignee request received")

	var req CreateAssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Assignee] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignee] User ID: %s, Ticket ID: %s, Target User ID: %s", userIDStr, req.TicketID, req.UserID)

	// Verify ticket exists and get board info
	ticket, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to get board_id
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create Assignee] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Find member by user_id and board_id
	member, err := ctrl.Repository.GetMemberByUserIDAndBoardID(req.UserID, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] User %s is not a member of board %s", req.UserID, column.BoardID)
		utils.JSON404(c, "User is not a member of this board")
		return
	}

	// Check if assignment already exists
	assigneeExists, err := ctrl.Repository.CheckTicketAssigneeExists(req.TicketID, member.ID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Failed to check existing assignment")
		utils.JSON500(c, "Failed to check existing assignment")
		return
	}
	if assigneeExists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Assignee] User already assigned to ticket")
		utils.JSON409(c, "User is already assigned to this ticket")
		return
	}

	// Create assignment
	assignee := &entity.TicketAssignee{
		TicketID: req.TicketID,
		MemberID: member.ID,
	}

	if err := ctrl.Repository.CreateTicketAssignee(assignee); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Failed to create assignment")
		utils.JSON500(c, "Failed to create assignment")
		return
	}

	// Load member details for response
	assignee.Member = *member

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignee] Assignment created successfully: %s", assignee.ID)
	utils.JSON201(c, gin.H{
		"message": "Member assigned successfully",
		"data":    assignee,
	})
}

// DeleteAssignee removes a member assignment from a ticket
func (ctrl *Controller) DeleteAssignee(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignee] Delete assignee request received")

	var req CreateAssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Assignee] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignee] User ID: %s, Ticket ID: %s, Target User ID: %s", userIDStr, req.TicketID, req.UserID)

	// Verify ticket exists and get board info
	ticket, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to get board_id
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Find member by user_id and board_id
	member, err := ctrl.Repository.GetMemberByUserIDAndBoardID(req.UserID, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] User %s is not a member of board %s", req.UserID, column.BoardID)
		utils.JSON404(c, "User is not a member of this board")
		return
	}

	// Check if assignment exists
	assigneeExists, err := ctrl.Repository.CheckTicketAssigneeExists(req.TicketID, member.ID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Failed to check existing assignment")
		utils.JSON500(c, "Failed to check existing assignment")
		return
	}
	if !assigneeExists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Assignee] Assignment not found")
		utils.JSON404(c, "Assignment not found")
		return
	}

	// Delete assignment
	if err := ctrl.Repository.DeleteTicketAssignee(req.TicketID, member.ID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Failed to delete assignment")
		utils.JSON500(c, "Failed to delete assignment")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignee] Assignment deleted successfully for ticket: %s, user: %s", req.TicketID, req.UserID)
	utils.JSON200(c, gin.H{
		"message": "Assignee removed successfully",
	})
}
