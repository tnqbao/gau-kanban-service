package controller

import (
	"gau-kanban-service/entity"
	"gau-kanban-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateChecklist tạo checklist mới cho ticket
func (ctrl *Controller) CreateChecklist(ctx *gin.Context) {
	var req CreateChecklistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Lấy vị trí tiếp theo cho checklist trong ticket
	maxPosition, err := ctrl.Repository.GetMaxChecklistPosition(req.TicketID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to get max position", err.Error())
		return
	}

	checklist := &entity.Checklist{
		TicketID:  req.TicketID,
		Title:     req.Title,
		Completed: false,
		Position:  maxPosition + 1,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateChecklist(checklist); err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to create checklist", err.Error())
		return
	}

	response := ChecklistDTO{
		ID:        checklist.ID,
		TicketID:  checklist.TicketID,
		Title:     checklist.Title,
		Completed: checklist.Completed,
		Position:  checklist.Position,
		CreatedAt: checklist.CreatedAt,
		UpdatedAt: checklist.UpdatedAt,
	}

	utils.SendSuccessResponse(ctx, http.StatusCreated, "Checklist created successfully", response)
}

// GetChecklistsByTicketID lấy tất cả checklist của một ticket
func (ctrl *Controller) GetChecklistsByTicketID(ctx *gin.Context) {
	ticketID := ctx.Param("ticketId")
	if ticketID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Ticket ID is required", "")
		return
	}

	checklists, err := ctrl.Repository.GetChecklistsByTicketID(ticketID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to get checklists", err.Error())
		return
	}

	var response []ChecklistDTO
	for _, checklist := range checklists {
		response = append(response, ChecklistDTO{
			ID:        checklist.ID,
			TicketID:  checklist.TicketID,
			Title:     checklist.Title,
			Completed: checklist.Completed,
			Position:  checklist.Position,
			CreatedAt: checklist.CreatedAt,
			UpdatedAt: checklist.UpdatedAt,
		})
	}

	utils.SendSuccessResponse(ctx, http.StatusOK, "Checklists retrieved successfully", response)
}

// UpdateChecklist cập nhật checklist
func (ctrl *Controller) UpdateChecklist(ctx *gin.Context) {
	checklistID := ctx.Param("id")
	if checklistID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Checklist ID is required", "")
		return
	}

	var req UpdateChecklistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	checklist, err := ctrl.Repository.GetChecklistByID(checklistID)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusNotFound, "Checklist not found", err.Error())
		return
	}

	// Cập nhật các field nếu có trong request
	if req.Title != nil {
		checklist.Title = *req.Title
	}
	if req.Completed != nil {
		checklist.Completed = *req.Completed
	}
	if req.Position != nil {
		checklist.Position = *req.Position
	}

	checklist.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateChecklist(checklist); err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to update checklist", err.Error())
		return
	}

	response := ChecklistDTO{
		ID:        checklist.ID,
		TicketID:  checklist.TicketID,
		Title:     checklist.Title,
		Completed: checklist.Completed,
		Position:  checklist.Position,
		CreatedAt: checklist.CreatedAt,
		UpdatedAt: checklist.UpdatedAt,
	}

	utils.SendSuccessResponse(ctx, http.StatusOK, "Checklist updated successfully", response)
}

// UpdateChecklistPosition cập nhật vị trí checklist
func (ctrl *Controller) UpdateChecklistPosition(ctx *gin.Context) {
	checklistID := ctx.Param("id")
	if checklistID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Checklist ID is required", "")
		return
	}

	var req UpdateChecklistPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := ctrl.Repository.UpdateChecklistPosition(checklistID, req.Position); err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to update checklist position", err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, http.StatusOK, "Checklist position updated successfully", nil)
}

// DeleteChecklist xóa checklist
func (ctrl *Controller) DeleteChecklist(ctx *gin.Context) {
	checklistID := ctx.Param("id")
	if checklistID == "" {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "Checklist ID is required", "")
		return
	}

	if err := ctrl.Repository.DeleteChecklist(checklistID); err != nil {
		utils.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete checklist", err.Error())
		return
	}

	utils.SendSuccessResponse(ctx, http.StatusOK, "Checklist deleted successfully", nil)
}
