package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateColumn tạo column mới
func (ctrl *Controller) CreateColumn(c *gin.Context) {
	var req CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	column := &entity.Column{
		Name:     req.Name,
		Position: req.Position,
	}

	if err := ctrl.Repository.CreateColumn(column); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column created suctrlessfully",
		"data":    column,
	})
}

// UpdateColumn cập nhật thông tin column
func (ctrl *Controller) UpdateColumn(c *gin.Context) {
	id := c.Param("id")
	var req UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	column, err := ctrl.Repository.GetColumnByID(id)
	if err != nil {
		utils.JSON404(c, "Column not found")
		return
	}

	if req.Name != "" {
		column.Name = req.Name
	}
	if req.Position != nil {
		column.Position = *req.Position
	}

	if err := ctrl.Repository.UpdateColumn(column); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column updated suctrlessfully",
		"data":    column,
	})
}

// DeleteColumn xóa column
func (ctrl *Controller) DeleteColumn(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.Repository.DeleteColumn(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column deleted suctrlessfully",
	})
}

// GetColumns lấy danh sách tất cả columns
func (ctrl *Controller) GetColumns(c *gin.Context) {
	columns, err := ctrl.Repository.GetAllColumn()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"data": columns,
	})
}

// UpdateColumnPosition thay đổi vị trí column
func (ctrl *Controller) UpdateColumnPosition(c *gin.Context) {
	id := c.Param("id")
	var req UpdateColumnPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := ctrl.Repository.UpdateColumnPosition(id, req.Position); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column position updated suctrlessfully",
	})
}

// Request/Response structures
type KanbanColumnResponse struct {
	ID      string                 `json:"id"`
	Title   string                 `json:"title"`
	Order   int                    `json:"order"`
	Tickets []KanbanTicketResponse `json:"tickets"`
}

type KanbanTicketResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TicketNo    string   `json:"ticketNo"`
	Tags        []string `json:"tags"`
	Assignees   []string `json:"assignees"`
	Completed   bool     `json:"completed"`
	DueDate     string   `json:"due_date,omitempty"`
	Priority    string   `json:"priority,omitempty"`
}

type CreateColumnRequest struct {
	Name     string `json:"name" binding:"required"`
	Position int    `json:"position"`
}

type UpdateColumnRequest struct {
	Name     string `json:"name"`
	Position *int   `json:"position"`
}

type UpdateColumnPositionRequest struct {
	Position int `json:"position" binding:"required"`
}

type MoveTicketRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
}

type MoveTicketWithPositionRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
	Position int    `json:"position" binding:"required"`
}

type UpdateTicketPositionRequest struct {
	Position int `json:"position" binding:"required"`
}

type CreateAssignmentRequest struct {
	TicketID     string `json:"ticket_id" binding:"required"`
	UserID       string `json:"user_id" binding:"required"`
	UserFullName string `json:"user_full_name" binding:"required"`
}

type UpdateAssignmentRequest struct {
	UserFullName string `json:"user_full_name"`
}

type CreateTicketRequest struct {
	ColumnID    string `json:"column_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
}

type UpdateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
}
