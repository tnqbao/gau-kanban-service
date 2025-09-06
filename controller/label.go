package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateLabel creates a new label for a board
func (ctrl *Controller) CreateLabel(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label] Create new label request received")

	var req CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Label] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label] User ID: %s, Board ID: %s", userIDStr, req.BoardID)

	// Verify board exists and user has access
	_, err := ctrl.Repository.GetBoardByID(req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label] Board not found: %s", req.BoardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Check if user is a member of the board
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label] Failed to check board membership")
		utils.JSON500(c, "Failed to check board access")
		return
	}
	if !isMember {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Label] User not authorized to access board")
		utils.JSON403(c, "Access denied: not a member of this board")
		return
	}

	// Create label
	label := &entity.Label{
		BoardID:     req.BoardID,
		Name:        req.Title,
		Color:       "#6B7280", // Default gray color
		Description: "",
	}

	if err := ctrl.Repository.CreateLabel(label); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label] Failed to create label")
		utils.JSON500(c, "Failed to create label")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label] Label created successfully: %s", label.ID)
	utils.JSON201(c, gin.H{
		"message": "Label created successfully",
		"data":    label,
	})
}

// CreateLabelTicket assigns a label to a ticket
func (ctrl *Controller) CreateLabelTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label Ticket] Create new label ticket assignment request received")

	var req CreateLabelTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Label Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label Ticket] User ID: %s, Label ID: %s, Ticket ID: %s", userIDStr, req.LabelID, req.TicketID)

	// Verify label exists
	label, err := ctrl.Repository.GetLabelByID(req.LabelID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Label not found: %s", req.LabelID)
		utils.JSON404(c, "Label not found")
		return
	}

	// Verify ticket exists
	_, err = ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Check if user has access to the board (through label's board)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, label.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Failed to check board membership")
		utils.JSON500(c, "Failed to check board access")
		return
	}
	if !isMember {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Label Ticket] User not authorized to access board")
		utils.JSON403(c, "Access denied: not a member of this board")
		return
	}

	// Check if label is already assigned to ticket
	labelExists, err := ctrl.Repository.CheckTicketLabelExists(req.TicketID, req.LabelID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Failed to check existing label assignment")
		utils.JSON500(c, "Failed to check existing label assignment")
		return
	}
	if labelExists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Label Ticket] Label already assigned to ticket")
		utils.JSON409(c, "Label is already assigned to this ticket")
		return
	}

	// Create ticket label assignment
	ticketLabel := &entity.TicketLabel{
		TicketID: req.TicketID,
		LabelID:  req.LabelID,
	}

	if err := ctrl.Repository.CreateTicketLabel(ticketLabel); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Label Ticket] Failed to assign label to ticket")
		utils.JSON500(c, "Failed to assign label to ticket")
		return
	}

	// Load label details for response
	ticketLabel.Label = *label

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Label Ticket] Label assigned to ticket successfully: %s", ticketLabel.ID)
	utils.JSON201(c, gin.H{
		"message": "Label assigned to ticket successfully",
		"data":    ticketLabel,
	})
}
