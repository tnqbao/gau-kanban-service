package controller

import (
	"time"

	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"

	"github.com/gin-gonic/gin"
)

// CreateChecklist tạo checklist mới cho ticket
func (ctrl *Controller) CreateChecklist(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Create Checklist] Create new checklist request received")

	var req CreateChecklistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Create Checklist] Invalid request body")
		utils.JSON400(ctx, "Invalid request body")
		return
	}

	// Lấy vị trí tiếp theo cho checklist trong ticket
	maxPosition, err := ctrl.Repository.GetMaxChecklistPosition(req.TicketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Create Checklist] Failed to get max checklist position for ticket: %s", req.TicketID)
		utils.JSON500(ctx, err.Error())
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
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Create Checklist] Failed to create checklist")
		utils.JSON500(ctx, err.Error())
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

	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Create Checklist] Checklist created successfully: %s", checklist.ID)
	utils.JSON200(ctx, gin.H{
		"message": "Checklist created successfully",
		"data":    response,
	})
}

// GetChecklistsByTicketID lấy tất cả checklist của một ticket
func (ctrl *Controller) GetChecklistsByTicketID(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	ticketID := ctx.Param("ticketId")
	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Get Checklists By Ticket ID] Get checklists request received for ticket ID: %s", ticketID)

	if ticketID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(requestCtx, "[Get Checklists By Ticket ID] Ticket ID is required")
		utils.JSON400(ctx, "Ticket ID is required")
		return
	}

	checklists, err := ctrl.Repository.GetChecklistsByTicketID(ticketID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Get Checklists By Ticket ID] Failed to get checklists for ticket: %s", ticketID)
		utils.JSON500(ctx, err.Error())
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

	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Get Checklists By Ticket ID] Retrieved %d checklists for ticket: %s", len(response), ticketID)
	utils.JSON200(ctx, gin.H{
		"message": "Checklists retrieved successfully",
		"data":    response,
	})
}

// UpdateChecklist cập nhật checklist
func (ctrl *Controller) UpdateChecklist(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	checklistID := ctx.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Update Checklist] Update checklist request received for ID: %s", checklistID)

	if checklistID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(requestCtx, "[Update Checklist] Checklist ID is required")
		utils.JSON400(ctx, "Checklist ID is required")
		return
	}

	var req UpdateChecklistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Update Checklist] Invalid request body")
		utils.JSON400(ctx, "Invalid request body")
		return
	}

	checklist, err := ctrl.Repository.GetChecklistByID(checklistID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Update Checklist] Checklist not found: %s", checklistID)
		utils.JSON404(ctx, "Checklist not found")
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
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Update Checklist] Failed to update checklist: %s", checklistID)
		utils.JSON500(ctx, err.Error())
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

	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Update Checklist] Checklist updated successfully: %s", checklistID)
	utils.JSON200(ctx, gin.H{
		"message": "Checklist updated successfully",
		"data":    response,
	})
}

// UpdateChecklistPosition cập nhật vị trí checklist
func (ctrl *Controller) UpdateChecklistPosition(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	checklistID := ctx.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Update Checklist Position] Update checklist position request received for ID: %s", checklistID)

	if checklistID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(requestCtx, "[Update Checklist Position] Checklist ID is required")
		utils.JSON400(ctx, "Checklist ID is required")
		return
	}

	var req UpdateChecklistPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Update Checklist Position] Invalid request body")
		utils.JSON400(ctx, "Invalid request body")
		return
	}

	if err := ctrl.Repository.UpdateChecklistPosition(checklistID, req.Position); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Update Checklist Position] Failed to update checklist position: %s", checklistID)
		utils.JSON500(ctx, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Update Checklist Position] Checklist position updated successfully: %s to position %d", checklistID, req.Position)
	utils.JSON200(ctx, gin.H{
		"message": "Checklist position updated successfully",
	})
}

// DeleteChecklist xóa checklist
func (ctrl *Controller) DeleteChecklist(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	checklistID := ctx.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Delete Checklist] Delete checklist request received for ID: %s", checklistID)

	if checklistID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(requestCtx, "[Delete Checklist] Checklist ID is required")
		utils.JSON400(ctx, "Checklist ID is required")
		return
	}

	if err := ctrl.Repository.DeleteChecklist(checklistID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(requestCtx, err, "[Delete Checklist] Failed to delete checklist: %s", checklistID)
		utils.JSON500(ctx, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(requestCtx, "[Delete Checklist] Checklist deleted successfully: %s", checklistID)
	utils.JSON200(ctx, gin.H{
		"message": "Checklist deleted successfully",
	})
}
