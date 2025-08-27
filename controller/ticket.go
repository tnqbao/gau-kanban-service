package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateTicket tạo ticket mới
func (ctrl *Controller) CreateTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Ticket] Create new ticket request received")

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra column có tồn tại không
	_, err := ctrl.Repository.GetByID(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Column not found: %s", req.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Lấy position cuối cùng trong column
	maxPosition, err := ctrl.Repository.GetMaxTicketPositionInColumn(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to get max ticket position in column: %s", req.ColumnID)
		utils.JSON500(c, err.Error())
		return
	}

	// Tạo ticket number theo format TASK-XXXX
	ticketNo, err := ctrl.Repository.GenerateTicketNumber()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to generate ticket number")
		utils.JSON500(c, err.Error())
		return
	}

	ticket := &entity.Ticket{
		TicketNo: ticketNo,
		ColumnID: req.ColumnID,
		Title:    req.Title,
		Position: maxPosition + 1,
		// DueDate is nil by default (pointer type), which will be NULL in database
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateTicket(ticket); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Ticket] Failed to create ticket")
		utils.JSON500(c, err.Error())
		return
	}

	// Lấy ticket với assignments và checklists
	ticketWithDetails, err := ctrl.Repository.GetTicketWithDetails(ticket.ID)
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Ticket] Ticket created successfully: %s", ticket.ID)
	utils.JSON200(c, gin.H{
		"message": "Ticket created successfully",
		"data":    ticketWithDetails,
	})
}

// GetTickets lấy danh sách tickets
func (ctrl *Controller) GetTickets(c *gin.Context) {
	ctx := c.Request.Context()
	columnID := c.Query("column_id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Tickets] Get tickets request received with column_id: %s", columnID)

	var tickets []entity.Ticket
	var err error

	if columnID != "" {
		tickets, err = ctrl.Repository.GetTicketsByColumnID(columnID)
	} else {
		tickets, err = ctrl.Repository.GetAllTickets()
	}

	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Tickets] Failed to get tickets")
		utils.JSON500(c, err.Error())
		return
	}

	// Lấy tickets với assignments và checklists
	var ticketsWithDetails []interface{}
	for _, ticket := range tickets {
		ticketDetails, err := ctrl.Repository.GetTicketWithDetails(ticket.ID)
		if err != nil {
			// Log error but continue with other tickets
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Tickets] Failed to get ticket details for %s", ticket.ID)
			continue
		}
		ticketsWithDetails = append(ticketsWithDetails, ticketDetails)
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Tickets] Retrieved %d tickets successfully", len(ticketsWithDetails))
	utils.JSON200(c, gin.H{
		"data": ticketsWithDetails,
	})
}

// GetTicketByID lấy ticket theo ID
func (ctrl *Controller) GetTicketByID(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket By ID] Get ticket request received for ID: %s", ticketID)

	ticket, err := ctrl.Repository.GetTicketWithDetails(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket By ID] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket By ID] Ticket retrieved successfully: %s", ticketID)
	utils.JSON200(c, gin.H{
		"data": ticket,
	})
}

// UpdateTicket cập nhật ticket
func (ctrl *Controller) UpdateTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket] Update ticket request received for ID: %s", ticketID)

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Cập nhật các field nếu có trong request
	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.DueDate != nil {
		ticket.DueDate = req.DueDate
	}
	if req.Priority != nil {
		ticket.Priority = req.Priority
	}

	ticket.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to update ticket: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	// Xử lý assignments nếu có
	if req.Assignments != nil {
		// Xóa assignments cũ
		if err := ctrl.Repository.DeleteAssignmentsByTicketID(ticketID); err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to delete old assignments for ticket: %s", ticketID)
			utils.JSON500(c, err.Error())
			return
		}

		// Tạo assignments mới
		for _, assignReq := range req.Assignments {
			assignment := &entity.TaskAssignment{
				TicketID:     ticketID,
				UserID:       assignReq.UserID,
				UserFullName: assignReq.UserFullName,
			}
			if err := ctrl.Repository.CreateAssignment(assignment); err != nil {
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to create assignment for ticket: %s", ticketID)
				utils.JSON500(c, err.Error())
				return
			}
		}
	}

	// Xử lý checklists nếu có
	if req.Checklists != nil {
		// Xóa checklists cũ
		if err := ctrl.Repository.DeleteChecklistsByTicketID(ticketID); err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to delete old checklists for ticket: %s", ticketID)
			utils.JSON500(c, err.Error())
			return
		}

		// Tạo checklists mới
		for i, checklistReq := range req.Checklists {
			checklist := &entity.Checklist{
				TicketID:  ticketID,
				Title:     checklistReq.Title,
				Completed: checklistReq.Completed,
				Position:  i + 1,
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			}
			if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to create checklist for ticket: %s", ticketID)
				utils.JSON500(c, err.Error())
				return
			}
		}
	}

	// Lấy ticket với details sau khi update
	ticketWithDetails, err := ctrl.Repository.GetTicketWithDetails(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket] Failed to get updated ticket details: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket] Ticket updated successfully: %s", ticketID)
	utils.JSON200(c, gin.H{
		"message": "Ticket updated successfully",
		"data":    ticketWithDetails,
	})
}

// DeleteTicket xóa ticket
func (ctrl *Controller) DeleteTicket(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Ticket] Delete ticket request received for ID: %s", ticketID)

	if err := ctrl.Repository.DeleteTicket(ticketID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Ticket] Failed to delete ticket: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Ticket] Ticket deleted successfully: %s", ticketID)
	utils.JSON200(c, gin.H{
		"message": "Ticket deleted successfully",
	})
}

// UpdateTicketPosition cập nhật vị trí ticket trong column
func (ctrl *Controller) UpdateTicketPosition(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket Position] Update ticket position request received for ID: %s", ticketID)

	var req UpdateTicketPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket Position] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.UpdateTicketPosition(ticketID, req.ColumnID, req.Position); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Ticket Position] Failed to update ticket position: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Ticket Position] Ticket position updated successfully: %s to position %d", ticketID, req.Position)
	utils.JSON200(c, gin.H{
		"message": "Ticket position updated successfully",
	})
}

// MoveTicketToColumn di chuyển ticket sang column khác
func (ctrl *Controller) MoveTicketToColumn(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket To Column] Move ticket to column request received")

	var req MoveTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket To Column] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.MoveTicketToColumn(req.TicketID, req.ColumnID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket To Column] Failed to move ticket %s to column %s", req.TicketID, req.ColumnID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket To Column] Ticket moved successfully: %s to column %s", req.TicketID, req.ColumnID)
	utils.JSON200(c, gin.H{
		"message": "Ticket moved successfully",
	})
}

// MoveTicketWithPosition di chuyển ticket sang column khác với position cụ thể
func (ctrl *Controller) MoveTicketWithPosition(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket With Position] Move ticket with position request received")

	var req MoveTicketWithPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket With Position] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra ticket có tồn tại không
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket With Position] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Kiểm tra column có tồn tại không
	_, err = ctrl.Repository.GetByID(req.ColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket With Position] Column not found: %s", req.ColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	if err := ctrl.Repository.MoveTicketToColumnWithPosition(req.TicketID, req.ColumnID, req.Position); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Move Ticket With Position] Failed to move ticket %s to column %s with position %d", req.TicketID, req.ColumnID, req.Position)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Move Ticket With Position] Ticket moved with position successfully: %s to column %s at position %d", req.TicketID, req.ColumnID, req.Position)
	utils.JSON200(c, gin.H{
		"message": "Ticket moved with position successfully",
	})
}

// ChangeTicketPosition thay đổi vị trí ticket với xử lý nâng cao (trong cùng column hoặc giữa các column)
func (ctrl *Controller) ChangeTicketPosition(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Change Ticket Position] Change ticket position request received for ID: %s", ticketID)

	var req ChangeTicketPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Ticket Position] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate new position is positive
	if req.NewPosition < 1 {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Change Ticket Position] Invalid position: %d", req.NewPosition)
		utils.JSON400(c, "Position must be greater than 0")
		return
	}

	// Check if ticket exists
	currentTicket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Ticket Position] Ticket not found: %s", ticketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Check if new column exists
	_, err = ctrl.Repository.GetByID(req.NewColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Ticket Position] Column not found: %s", req.NewColumnID)
		utils.JSON404(c, "Column not found")
		return
	}

	// Get max position in target column to validate new position
	maxPosition, err := ctrl.Repository.GetMaxTicketPositionInColumn(req.NewColumnID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Ticket Position] Failed to get max ticket position in column: %s", req.NewColumnID)
		utils.JSON500(c, "Failed to get max ticket position: "+err.Error())
		return
	}

	// If moving to same column, max position should not change
	// If moving to different column, allow position up to max+1 (for appending)
	if currentTicket.ColumnID != req.NewColumnID && req.NewPosition > maxPosition+1 {
		req.NewPosition = maxPosition + 1 // Set to last position + 1 if exceeds
	} else if currentTicket.ColumnID == req.NewColumnID && req.NewPosition > maxPosition {
		req.NewPosition = maxPosition // Set to last position if exceeds within same column
	}

	// Change ticket position
	if err := ctrl.Repository.ChangeTicketPosition(ticketID, req.NewColumnID, req.NewPosition); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Change Ticket Position] Failed to change ticket position: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	// Log the change details
	if currentTicket.ColumnID == req.NewColumnID {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Change Ticket Position] Ticket position changed within same column: %s from position %d to %d", ticketID, currentTicket.Position, req.NewPosition)
	} else {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Change Ticket Position] Ticket moved between columns: %s from column %s (position %d) to column %s (position %d)", ticketID, currentTicket.ColumnID, currentTicket.Position, req.NewColumnID, req.NewPosition)
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket position changed successfully",
		"data": gin.H{
			"ticket_id":     ticketID,
			"old_column_id": currentTicket.ColumnID,
			"old_position":  currentTicket.Position,
			"new_column_id": req.NewColumnID,
			"new_position":  req.NewPosition,
		},
	})
}
