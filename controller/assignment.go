package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateAssignment tạo assignment mới cho ticket
func (ctrl *Controller) CreateAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignment] Create new assignment request received")

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignment] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra ticket có tồn tại không
	_, err := ctrl.Repository.GetTicketByID(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignment] Ticket not found: %s", req.TicketID)
		utils.JSON404(c, "Ticket not found")
		return
	}

	assignment := &entity.TaskAssignment{
		TicketID:     req.TicketID,
		UserID:       req.UserID,
		UserFullName: req.UserFullName,
	}

	if err := ctrl.Repository.CreateAssignment(assignment); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Assignment] Failed to create assignment")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Assignment] Assignment created successfully: %s", assignment.ID)
	utils.JSON200(c, gin.H{
		"message": "Assignment created successfully",
		"data":    assignment,
	})
}

// UpdateAssignment cập nhật thông tin assignment
func (ctrl *Controller) UpdateAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Assignment] Update assignment request received for ID: %s", id)

	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Assignment] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	assignment, err := ctrl.Repository.GetAssignmentByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Assignment] Assignment not found: %s", id)
		utils.JSON404(c, "Assignment not found")
		return
	}

	if req.UserFullName != "" {
		assignment.UserFullName = req.UserFullName
	}

	if err := ctrl.Repository.UpdateAssignment(assignment); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Assignment] Failed to update assignment: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Assignment] Assignment updated successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Assignment updated successfully",
		"data":    assignment,
	})
}

// DeleteAssignment xóa assignment
func (ctrl *Controller) DeleteAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignment] Delete assignment request received for ID: %s", id)

	if err := ctrl.Repository.DeleteAssignment(id); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignment] Failed to delete assignment: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignment] Assignment deleted successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Assignment deleted successfully",
	})
}

// GetTicketAssignments lấy danh sách assignments của ticket
func (ctrl *Controller) GetTicketAssignments(c *gin.Context) {
	ctx := c.Request.Context()
	ticketID := c.Param("ticket_id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket Assignments] Get assignments request received for ticket ID: %s", ticketID)

	assignments, err := ctrl.Repository.GetAssignmentsByTicketID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Ticket Assignments] Failed to get assignments for ticket: %s", ticketID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Ticket Assignments] Retrieved %d assignments for ticket: %s", len(assignments), ticketID)
	utils.JSON200(c, gin.H{
		"data": assignments,
	})
}

// DeleteAssignmentsByUserID xóa tất cả assignments của một user
func (ctrl *Controller) DeleteAssignmentsByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignments By User ID] Delete assignments request received for user ID: %s", userID)

	if userID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Delete Assignments By User ID] User ID is required")
		utils.JSON400(c, "User ID is required")
		return
	}

	if err := ctrl.Repository.DeleteAssignmentsByUserID(userID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Assignments By User ID] Failed to delete assignments for user: %s", userID)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Assignments By User ID] All assignments deleted successfully for user: %s", userID)
	utils.JSON200(c, gin.H{
		"message": "All assignments for user deleted successfully",
	})
}
