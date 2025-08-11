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

// CreateTicket tạo ticket mới trong column với checklist
func (ctrl *Controller) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Tạo ticket trước
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

	// Tạo checklists nếu có
	var checklists []ChecklistDTO
	if req.Checklists != nil && len(req.Checklists) > 0 {
		for i, checklistReq := range req.Checklists {
			checklist := &entity.Checklist{
				TicketID:  ticket.ID,
				Title:     checklistReq.Title,
				Completed: false,
				Position:  i + 1,
			}

			if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
				utils.JSON500(c, "Failed to create checklist: "+err.Error())
				return
			}

			checklists = append(checklists, ChecklistDTO{
				ID:        checklist.ID,
				TicketID:  checklist.TicketID,
				Title:     checklist.Title,
				Completed: checklist.Completed,
				Position:  checklist.Position,
				CreatedAt: checklist.CreatedAt,
				UpdatedAt: checklist.UpdatedAt,
			})
		}
	}

	response := TicketWithChecklistsDTO{
		ID:          ticket.ID,
		TicketNo:    ticket.TicketNo,
		ColumnID:    ticket.ColumnID,
		Title:       ticket.Title,
		Description: ticket.Description,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		Position:    ticket.Position,
		Checklists:  checklists,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket created successfully",
		"data":    response,
	})
}

// UpdateTicket cập nhật thông tin ticket với checklist
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

	// Cập nhật thông tin ticket
	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.DueDate != nil {
		ticket.DueDate = *req.DueDate
	}
	if req.Priority != nil {
		ticket.Priority = *req.Priority
	}

	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Xử lý checklist nếu có trong request
	var checklists []ChecklistDTO
	if req.Checklists != nil {
		// Xóa tất cả checklist cũ của ticket
		if err := ctrl.Repository.DeleteChecklistsByTicketID(ticket.ID); err != nil {
			utils.JSON500(c, "Failed to delete old checklists: "+err.Error())
			return
		}

		// Tạo checklist mới
		for i, checklistReq := range req.Checklists {
			checklist := &entity.Checklist{
				TicketID:  ticket.ID,
				Title:     *checklistReq.Title,
				Completed: false,
				Position:  i + 1,
			}

			if checklistReq.Completed != nil {
				checklist.Completed = *checklistReq.Completed
			}

			if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
				utils.JSON500(c, "Failed to create checklist: "+err.Error())
				return
			}

			checklists = append(checklists, ChecklistDTO{
				ID:        checklist.ID,
				TicketID:  checklist.TicketID,
				Title:     checklist.Title,
				Completed: checklist.Completed,
				Position:  checklist.Position,
				CreatedAt: checklist.CreatedAt,
				UpdatedAt: checklist.UpdatedAt,
			})
		}
	} else {
		// Nếu không có checklist trong request, lấy checklist hiện tại
		existingChecklists, err := ctrl.Repository.GetChecklistsByTicketID(ticket.ID)
		if err == nil {
			for _, checklist := range existingChecklists {
				checklists = append(checklists, ChecklistDTO{
					ID:        checklist.ID,
					TicketID:  checklist.TicketID,
					Title:     checklist.Title,
					Completed: checklist.Completed,
					Position:  checklist.Position,
					CreatedAt: checklist.CreatedAt,
					UpdatedAt: checklist.UpdatedAt,
				})
			}
		}
	}

	response := TicketWithChecklistsDTO{
		ID:          ticket.ID,
		TicketNo:    ticket.TicketNo,
		ColumnID:    ticket.ColumnID,
		Title:       ticket.Title,
		Description: ticket.Description,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		Position:    ticket.Position,
		Checklists:  checklists,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket updated successfully",
		"data":    response,
	})
}

// DeleteTicket xóa ticket và tất cả checklist liên quan
func (ctrl *Controller) DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	// Xóa tất cả checklist của ticket trước
	if err := ctrl.Repository.DeleteChecklistsByTicketID(id); err != nil {
		utils.JSON500(c, "Failed to delete checklists: "+err.Error())
		return
	}

	// Xóa ticket
	if err := ctrl.Repository.DeleteTicket(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket and related checklists deleted successfully",
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

// GetTicketWithChecklists lấy thông tin ticket với checklist và assignees
func (ctrl *Controller) GetTicketWithChecklists(c *gin.Context) {
	id := c.Param("id")

	ticket, assignments, err := ctrl.Repository.GetTicketWithAssignments(id)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	// Lấy checklists của ticket
	checklists, err := ctrl.Repository.GetChecklistsByTicketID(id)
	if err != nil {
		utils.JSON500(c, "Failed to get checklists: "+err.Error())
		return
	}

	var checklistsDTO []ChecklistDTO
	for _, checklist := range checklists {
		checklistsDTO = append(checklistsDTO, ChecklistDTO{
			ID:        checklist.ID,
			TicketID:  checklist.TicketID,
			Title:     checklist.Title,
			Completed: checklist.Completed,
			Position:  checklist.Position,
			CreatedAt: checklist.CreatedAt,
			UpdatedAt: checklist.UpdatedAt,
		})
	}

	response := TicketWithChecklistsDTO{
		ID:          ticket.ID,
		TicketNo:    ticket.TicketNo,
		ColumnID:    ticket.ColumnID,
		Title:       ticket.Title,
		Description: ticket.Description,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		Position:    ticket.Position,
		Checklists:  checklistsDTO,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	utils.JSON200(c, gin.H{
		"data":      response,
		"assignees": assignments,
	})
}
