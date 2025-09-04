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
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignee] User ID: %s, Ticket ID: %s, Member ID: %s", userIDStr, req.TicketID, req.MemberID)

	// Verify ticket exists
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Verify member exists
	member, err := ctrl.Repository.GetMemberByID(req.MemberID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Member not found: %s", req.MemberID)
		utils.JSON404(c, "Member not found")
		return
	}

	// Check if assignment already exists
	assigneeExists, err := ctrl.Repository.CheckTicketAssigneeExists(req.TicketID, req.MemberID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignee] Failed to check existing assignment")
		utils.JSON500(c, "Failed to check existing assignment")
		return
	}
	if assigneeExists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Assignee] Member already assigned to ticket")
		utils.JSON409(c, "Member is already assigned to this ticket")
		return
	}

	// Create assignment
	assignee := &entity.TicketAssignee{
		TicketID: req.TicketID,
		MemberID: req.MemberID,
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
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignee] User ID: %s, Ticket ID: %s, Member ID: %s", userIDStr, req.TicketID, req.MemberID)

	// Verify ticket exists
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Verify member exists
	_, err = ctrl.Repository.GetMemberByID(req.MemberID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Member not found: %s", req.MemberID)
		utils.JSON404(c, "Member not found")
		return
	}

	// Check if assignment exists
	assigneeExists, err := ctrl.Repository.CheckTicketAssigneeExists(req.TicketID, req.MemberID)
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
	if err := ctrl.Repository.DeleteTicketAssignee(req.TicketID, req.MemberID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignee] Failed to delete assignment")
		utils.JSON500(c, "Failed to delete assignment")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignee] Assignment deleted successfully for ticket: %s, member: %s", req.TicketID, req.MemberID)
	utils.JSON200(c, gin.H{
		"message": "Assignee removed successfully",
	})
}
