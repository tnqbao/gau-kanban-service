package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateTicket tạo ticket mới
func (ctrl *Controller) CreateTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Ticket] Create ticket request received")

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get column and verify it exists
	column, err := ctrl.Repository.GetColumnByID(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Column not found: %s", req.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create Ticket] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Generate ticket number
	ticketNumber, err := ctrl.Repository.GenerateNextTicketNumber()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to generate ticket number")
		utils.JSON500(c, "Failed to generate ticket number")
		return
	}

	// Get max order for this column
	maxOrder, err := ctrl.Repository.GetMaxOrderByColumnID(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to get max order for column: %s", req.ColumnID)
		utils.JSON500(c, "Failed to calculate ticket order")
		return
	}

	// Create ticket with user provided title and default values
	ticket := &entity.Ticket{
		ColumnID:     req.ColumnID,
		Title:        req.Title,
		TicketNumber: ticketNumber,
		Description:  "",
		Status:       "open",
		Priority:     "medium",
		Position:     maxOrder + 1000, // Use 1000 spacing for fractional ordering
		DueDate:      nil,
	}

	if err := ctrl.Repository.CreateTicket(ticket); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to create ticket")
		utils.JSON500(c, "Failed to create ticket")
		return
	}

	// Initialize empty arrays for relationships
	ticket.Labels = []entity.TicketLabel{}
	ticket.Assignees = []entity.TicketAssignee{}
	ticket.Checklists = []entity.Checklist{}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Ticket] Ticket created successfully: %s (%s)", ticket.ID, ticket.TicketNumber)
	utils.JSON200(c, gin.H{
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}

// UpdateTicket updates an existing ticket
func (ctrl *Controller) UpdateTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket] Update ticket request received")

	ticketID := c.Param("id")
	if ticketID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Ticket] Ticket ID is required")
		utils.JSON400(c, "Ticket ID is required")
		return
	}

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get ticket and verify it exists
	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to check board access
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Update Ticket] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Update fields if provided
	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.Priority != nil {
		// Validate priority values
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
		if !validPriorities[*req.Priority] {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Ticket] Invalid priority: %s", *req.Priority)
			utils.JSON400(c, "Invalid priority. Must be one of: low, medium, high, critical")
			return
		}
		ticket.Priority = *req.Priority
	}
	if req.Status != nil {
		// Validate status values
		validStatuses := map[string]bool{"open": true, "in_progress": true, "in_review": true, "done": true, "blocked": true}
		if !validStatuses[*req.Status] {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Ticket] Invalid status: %s", *req.Status)
			utils.JSON400(c, "Invalid status. Must be one of: open, in_progress, in_review, done, blocked")
			return
		}
		ticket.Status = *req.Status
	}
	if req.DueDate != nil {
		if *req.DueDate == "" {
			// Clear due date
			ticket.DueDate = nil
		} else {
			// Parse due date
			dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Invalid due date format: %s", *req.DueDate)
				utils.JSON400(c, "Invalid due date format. Use ISO 8601 format (e.g., 2023-12-31T23:59:59Z)")
				return
			}
			ticket.DueDate = &dueDate
		}
	}

	// Update ticket in database
	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to update ticket")
		utils.JSON500(c, "Failed to update ticket")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket] Ticket updated successfully: %s (%s)", ticket.ID, ticket.TicketNumber)
	utils.JSON200(c, gin.H{
		"message": "Ticket updated successfully",
		"data":    ticket,
	})
}

// DeleteTicket deletes a ticket
func (ctrl *Controller) DeleteTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Ticket] Delete ticket request received")

	ticketID := c.Param("id")
	if ticketID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Ticket] Ticket ID is required")
		utils.JSON400(c, "Ticket ID is required")
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get ticket and verify it exists
	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to check board access
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Delete Ticket] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Delete ticket from database
	if err := ctrl.Repository.DeleteTicket(ticketID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Failed to delete ticket")
		utils.JSON500(c, "Failed to delete ticket")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Ticket] Ticket deleted successfully: %s (%s)", ticket.ID, ticket.TicketNumber)
	utils.JSON200(c, gin.H{
		"message": "Ticket deleted successfully",
		"data": gin.H{
			"id":            ticket.ID,
			"ticket_number": ticket.TicketNumber,
			"title":         ticket.Title,
		},
	})
}

// GetTicketByID retrieves a ticket by ID with all related data
func (ctrl *Controller) GetTicketByID(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket] Get ticket by ID request received")

	ticketID := c.Param("id")
	if ticketID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Get Ticket] Ticket ID is required")
		utils.JSON400(c, "Ticket ID is required")
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Get Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get ticket and verify it exists
	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get column to check board access
	column, err := ctrl.Repository.GetColumnByID(ticket.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Column not found: %s", ticket.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, column.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(column.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Board not found: %s", column.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Get Ticket] User %s does not have access to board %s", userIDStr, column.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Get assignees for this ticket
	assignees, err := ctrl.Repository.GetTicketAssigneesByTicketID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Failed to get assignees for ticket: %s", ticketID)
		// Don't fail the request, just set empty array
		assignees = []entity.TicketAssignee{}
	}

	// Get labels for this ticket
	labels, err := ctrl.Repository.GetTicketLabelsByTicketID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Failed to get labels for ticket: %s", ticketID)
		// Don't fail the request, just set empty array
		labels = []entity.TicketLabel{}
	}

	// Get checklists for this ticket
	checklists, err := ctrl.Repository.GetChecklistsByTicketID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket] Failed to get checklists for ticket: %s", ticketID)
		// Don't fail the request, just set empty array
		checklists = []entity.Checklist{}
	}

	// Set the relationships on the ticket
	ticket.Assignees = assignees
	ticket.Labels = labels
	ticket.Checklists = checklists

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket] Ticket retrieved successfully: %s (%s)", ticket.ID, ticket.TicketNumber)
	utils.JSON200(c, gin.H{
		"message": "Ticket retrieved successfully",
		"data":    ticket,
	})
}

// MoveTicket moves a ticket to a new column and/or position
func (ctrl *Controller) MoveTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket] Move ticket request received")

	ticketID := c.Param("id")
	if ticketID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] Ticket ID is required")
		utils.JSON400(c, "Ticket ID is required")
		return
	}

	var req MoveTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get ticket and verify it exists
	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Get target column and verify it exists
	targetColumn, err := ctrl.Repository.GetColumnByID(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Target column not found: %s", req.ColumnID)
		utils.JSON404(c, "Target column not found")
		return
	}

	// Check if user has access to the target board
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, targetColumn.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(targetColumn.BoardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Board not found: %s", targetColumn.BoardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Move Ticket] User %s does not have access to board %s", userIDStr, targetColumn.BoardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Check for self-reference in position
	if req.Position == "after:"+ticketID || req.Position == "before:"+ticketID {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] Cannot move ticket relative to itself")
		utils.JSON400(c, "Cannot move ticket relative to itself")
		return
	}

	var newPosition int

	// Calculate new position based on the position parameter
	switch {
	case req.Position == "first":
		// Get minimum position in target column
		minPosition, err := ctrl.Repository.GetMinPositionByColumnID(req.ColumnID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to get min position")
			utils.JSON500(c, "Failed to calculate position")
			return
		}

		if minPosition == 0 {
			// Column is empty, set first position
			newPosition = 1000
		} else {
			newPosition = minPosition - 1000
		}

	case req.Position == "last":
		// Get maximum position in target column
		maxPosition, err := ctrl.Repository.GetMaxPositionByColumnID(req.ColumnID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to get max position")
			utils.JSON500(c, "Failed to calculate position")
			return
		}

		if maxPosition == 0 {
			// Column is empty, set first position
			newPosition = 1000
		} else {
			newPosition = maxPosition + 1000
		}

	case len(req.Position) > 6 && req.Position[:6] == "after:":
		referenceTicketID := req.Position[6:]

		// Get reference ticket
		refTicket, err := ctrl.Repository.GetTicketByID(referenceTicketID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Reference ticket not found: %s", referenceTicketID)
			utils.JSON404(c, "Reference ticket not found")
			return
		}

		// Verify reference ticket is in the same target column
		if refTicket.ColumnID != req.ColumnID {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] Reference ticket is not in target column")
			utils.JSON400(c, "Reference ticket must be in the target column")
			return
		}

		// Get next ticket position after reference
		nextPosition, err := ctrl.Repository.GetNextPositionAfter(req.ColumnID, refTicket.Position)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to get next position")
			utils.JSON500(c, "Failed to calculate position")
			return
		}

		if nextPosition == 0 {
			// Reference ticket is last, place after it
			newPosition = refTicket.Position + 1000
		} else {
			// Place between reference and next ticket
			newPosition = (refTicket.Position + nextPosition) / 2
		}

	case len(req.Position) > 7 && req.Position[:7] == "before:":
		referenceTicketID := req.Position[7:]

		// Get reference ticket
		refTicket, err := ctrl.Repository.GetTicketByID(referenceTicketID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Reference ticket not found: %s", referenceTicketID)
			utils.JSON404(c, "Reference ticket not found")
			return
		}

		// Verify reference ticket is in the same target column
		if refTicket.ColumnID != req.ColumnID {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] Reference ticket is not in target column")
			utils.JSON400(c, "Reference ticket must be in the target column")
			return
		}

		// Get previous ticket position before reference
		prevPosition, err := ctrl.Repository.GetPreviousPositionBefore(req.ColumnID, refTicket.Position)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to get previous position")
			utils.JSON500(c, "Failed to calculate position")
			return
		}

		if prevPosition == 0 {
			// Reference ticket is first, place before it
			newPosition = refTicket.Position - 1000
		} else {
			// Place between previous and reference ticket
			newPosition = (prevPosition + refTicket.Position) / 2
		}

	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Move Ticket] Invalid position format: %s", req.Position)
		utils.JSON400(c, "Invalid position format. Use 'first', 'last', 'after:ticket_id', or 'before:ticket_id'")
		return
	}

	// Update ticket with new column and position
	ticket.ColumnID = req.ColumnID
	ticket.Position = newPosition

	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket] Failed to move ticket")
		utils.JSON500(c, "Failed to move ticket")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket] Ticket moved successfully: %s (%s) to column %s at position %d",
		ticket.ID, ticket.TicketNumber, req.ColumnID, newPosition)

	utils.JSON200(c, gin.H{
		"message": "Ticket moved successfully",
		"data":    ticket,
	})
}
