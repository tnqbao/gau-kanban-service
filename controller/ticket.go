package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// MoveTicketToColumn di chuyển ticket sang column khác
func (ctrl *Controller) MoveTicketToColumn(c *gin.Context) {
	var req MoveTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.MoveTicketToColumn(req.TicketID, req.ColumnID); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket moved suctrlessfully",
	})
}

// MoveTicketWithPosition di chuyển ticket sang column khác với position cụ thể
func (ctrl *Controller) MoveTicketWithPosition(c *gin.Context) {
	var req MoveTicketWithPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra ticket có tồn tại không
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Kiểm tra column có tồn tại không
	_, err = ctrl.Repository.GetColumnByID(req.ColumnID)
	if err != nil {
		utils.JSON404(c, "Column not found")
		return
	}

	if err := ctrl.Repository.MoveTicketToColumnWithPosition(req.TicketID, req.ColumnID, req.Position); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket moved with position successfully",
	})
}

// CreateTicket tạo ticket mới trong column
func (ctrl *Controller) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ticket := &entity.Ticket{
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}

	if err := ctrl.Repository.CreateTicket(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket created suctrlessfully",
		"data":    ticket,
	})
}

// UpdateTicket cập nhật thông tin ticket
func (ctrl *Controller) UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ticket, err := ctrl.Repository.GetTicketByID(id)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	if req.Title != "" {
		ticket.Title = req.Title
	}
	if req.Description != "" {
		ticket.Description = req.Description
	}
	if req.DueDate != "" {
		ticket.DueDate = req.DueDate
	}
	if req.Priority != "" {
		ticket.Priority = req.Priority
	}

	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket updated suctrlessfully",
		"data":    ticket,
	})
}

// DeleteTicket xóa ticket
func (ctrl *Controller) DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.Repository.DeleteTicket(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket deleted suctrlessfully",
	})
}

// UpdateTicketPosition cập nhật vị trí của ticket trong column
func (ctrl *Controller) UpdateTicketPosition(c *gin.Context) {
	ticketID := c.Param("id")
	var req UpdateTicketPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra ticket có tồn tại không
	_, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	if err := ctrl.Repository.UpdateTicketPosition(ticketID, req.Position); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket position updated successfully",
	})
}

// GetTicket lấy thông tin ticket theo ID với assignees
func (ctrl *Controller) GetTicket(c *gin.Context) {
	id := c.Param("id")

	ticket, assignments, err := ctrl.Repository.GetTicketWithAssignments(id)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Tạo response với ticket và assignees
	response := gin.H{
		"id":          ticket.ID,
		"ticket_no":   ticket.TicketNo,
		"column_id":   ticket.ColumnID,
		"title":       ticket.Title,
		"description": ticket.Description,
		"due_date":    ticket.DueDate,
		"priority":    ticket.Priority,
		"position":    ticket.Position,
		"created_at":  ticket.CreatedAt,
		"updated_at":  ticket.UpdatedAt,
		"assignees":   assignments,
	}

	utils.JSON200(c, gin.H{
		"data": response,
	})
}

// GetTickets lấy danh sách tất cả tickets với assignees
func (ctrl *Controller) GetTickets(c *gin.Context) {
	tickets, err := ctrl.Repository.GetAllTicket()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Lấy assignees cho từng ticket
	var ticketsWithAssignees []gin.H
	for _, ticket := range tickets {
		assignments, err := ctrl.Repository.GetAssignmentsByTicketID(ticket.ID)
		if err != nil {
			// Nếu có lỗi khi lấy assignments, vẫn trả về ticket nhưng assignees rỗng
			assignments = []entity.TaskAssignment{}
		}

		ticketResponse := gin.H{
			"id":          ticket.ID,
			"ticket_no":   ticket.TicketNo,
			"column_id":   ticket.ColumnID,
			"title":       ticket.Title,
			"description": ticket.Description,
			"due_date":    ticket.DueDate,
			"priority":    ticket.Priority,
			"position":    ticket.Position,
			"created_at":  ticket.CreatedAt,
			"updated_at":  ticket.UpdatedAt,
			"assignees":   assignments,
		}
		ticketsWithAssignees = append(ticketsWithAssignees, ticketResponse)
	}

	utils.JSON200(c, gin.H{
		"data": ticketsWithAssignees,
	})
}
