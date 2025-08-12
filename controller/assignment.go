package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateAssignment tạo assignment mới cho ticket
func (ctrl *Controller) CreateAssignment(c *gin.Context) {
	var req
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

	assignment := &entity.TaskAssignment{
		TicketID:     req.TicketID,
		UserID:       req.UserID,
		UserFullName: req.UserFullName,
	}

	if err := ctrl.Repository.CreateAssignment(assignment); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Assignment created successfully",
		"data":    assignment,
	})
}

// UpdateAssignment cập nhật thông tin assignment
func (ctrl *Controller) UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	var req
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	assignment, err := ctrl.Repository.GetAssignmentByID(id)
	if err != nil {
		utils.JSON404(c, "Assignment not found")
		return
	}

	if req. != "" {
		assignment.UserFullName = req.UserFullName
	}

	if err := ctrl.Repository.UpdateAssignment(assignment); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Assignment updated successfully",
		"data":    assignment,
	})
}

// DeleteAssignment xóa assignment
func (ctrl *Controller) DeleteAssignment(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.Repository.DeleteAssignment(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Assignment deleted successfully",
	})
}

// GetTicketAssignments lấy danh sách assignments của ticket
func (ctrl *Controller) GetTicketAssignments(c *gin.Context) {
	ticketID := c.Param("ticket_id")

	assignments, err := ctrl.Repository.GetAssignmentsByTicketID(ticketID)
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"data": assignments,
	})
}

// DeleteAssignmentsByUserID xóa tất cả assignments của một user
func (ctrl *Controller) DeleteAssignmentsByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		utils.JSON400(c, "User ID is required")
		return
	}

	if err := ctrl.Repository.DeleteAssignmentsByUserID(userID); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "All assignments for user deleted successfully",
	})
}
