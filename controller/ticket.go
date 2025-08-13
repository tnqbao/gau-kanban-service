package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateTicket tạo ticket mới
func (ctrl *Controller) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra column có tồn tại không
	_, err := ctrl.Repository.GetColumnByID(req.ColumnID)
	if err != nil {
		utils.JSON404(c, "Column not found")
		return
	}

	// Lấy position cuối cùng trong column
	maxPosition, err := ctrl.Repository.GetMaxTicketPositionInColumn(req.ColumnID)
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Tạo ticket number theo format TASK-XXXX
	ticketNo, err := ctrl.Repository.GenerateTicketNumber()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	ticket := &entity.Ticket{
		TicketNo:    ticketNo,
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Position:    maxPosition + 1, // Đặt ở cuối column
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateTicket(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Tạo assignments nếu có
	if req.Assignments != nil && len(req.Assignments) > 0 {
		for _, assignReq := range req.Assignments {
			assignment := &entity.TaskAssignment{
				TicketID:     ticket.ID,
				UserID:       assignReq.UserID,
				UserFullName: assignReq.UserFullName,
			}
			if err := ctrl.Repository.CreateAssignment(assignment); err != nil {
				// Log error but don't fail the whole operation
				fmt.Printf("Failed to create assignment: %v\n", err)
			}
		}
	}

	// Tạo checklists nếu có
	if req.Checklists != nil && len(req.Checklists) > 0 {
		for i, checklistReq := range req.Checklists {
			checklist := &entity.Checklist{
				TicketID:  ticket.ID,
				Title:     checklistReq.Title,
				Completed: false,
				Position:  i + 1,
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			}
			if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
				// Log error but don't fail the whole operation
				fmt.Printf("Failed to create checklist: %v\n", err)
			}
		}
	}

	// Lấy ticket với assignments và checklists
	ticketWithDetails, err := ctrl.Repository.GetTicketWithDetails(ticket.ID)
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket created successfully",
		"data":    ticketWithDetails,
	})
}

// GetTickets lấy danh sách tickets
func (ctrl *Controller) GetTickets(c *gin.Context) {
	columnID := c.Query("column_id")

	var tickets []entity.Ticket
	var err error

	if columnID != "" {
		tickets, err = ctrl.Repository.GetTicketsByColumnID(columnID)
	} else {
		tickets, err = ctrl.Repository.GetAllTickets()
	}

	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Lấy tickets với assignments và checklists
	var ticketsWithDetails []interface{}
	for _, ticket := range tickets {
		ticketDetails, err := ctrl.Repository.GetTicketWithDetails(ticket.ID)
		if err != nil {
			// Log error but continue with other tickets
			fmt.Printf("Failed to get ticket details for %s: %v\n", ticket.ID, err)
			continue
		}
		ticketsWithDetails = append(ticketsWithDetails, ticketDetails)
	}

	utils.JSON200(c, gin.H{
		"data": ticketsWithDetails,
	})
}

// GetTicketByID lấy ticket theo ID
func (ctrl *Controller) GetTicketByID(c *gin.Context) {
	ticketID := c.Param("id")

	ticket, err := ctrl.Repository.GetTicketWithDetails(ticketID)
	if err != nil {
		utils.JSON404(c, "Ticket not found")
		return
	}

	utils.JSON200(c, gin.H{
		"data": ticket,
	})
}

// UpdateTicket cập nhật ticket
func (ctrl *Controller) UpdateTicket(c *gin.Context) {
	ticketID := c.Param("id")
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ticket, err := ctrl.Repository.GetTicketByID(ticketID)
	if err != nil {
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
		ticket.DueDate = *req.DueDate
	}
	if req.Priority != nil {
		ticket.Priority = *req.Priority
	}

	ticket.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateTicket(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Xử lý assignments nếu có
	if req.Assignments != nil {
		// Xóa assignments cũ
		if err := ctrl.Repository.DeleteAssignmentsByTicketID(ticketID); err != nil {
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
				utils.JSON500(c, err.Error())
				return
			}
		}
	}

	// Xử lý checklists nếu có
	if req.Checklists != nil {
		// Xóa checklists cũ
		if err := ctrl.Repository.DeleteChecklistsByTicketID(ticketID); err != nil {
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
				utils.JSON500(c, err.Error())
				return
			}
		}
	}

	// Lấy ticket với details sau khi update
	ticketWithDetails, err := ctrl.Repository.GetTicketWithDetails(ticketID)
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket updated successfully",
		"data":    ticketWithDetails,
	})
}

// DeleteTicket xóa ticket
func (ctrl *Controller) DeleteTicket(c *gin.Context) {
	ticketID := c.Param("id")

	if err := ctrl.Repository.DeleteTicket(ticketID); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket deleted successfully",
	})
}

// UpdateTicketPosition cập nhật vị trí ticket trong column
func (ctrl *Controller) UpdateTicketPosition(c *gin.Context) {
	ticketID := c.Param("id")
	var req UpdateTicketPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.UpdateTicketPosition(ticketID, req.ColumnID, req.Position); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket position updated successfully",
	})
}

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
		"message": "Ticket moved successfully",
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
