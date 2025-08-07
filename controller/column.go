package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/repository"
	"github.com/tnqbao/gau-kanban-service/utils"
)

type ColumnController struct {
	r repository.RepositoryInterface
}

func NewColumnController(r repository.RepositoryInterface) *ColumnController {
	return &ColumnController{
		r: r,
	}
}

// GetKanbanBoard trả về toàn bộ kanban board với format như yêu cầu
func (cc *ColumnController) GetKanbanBoard(c *gin.Context) {
	columns, err := cc.r.GetRepositories().Column.GetAllWithFullTicketDetails()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	// Convert to format like initialColumns
	kanbanColumns := make([]KanbanColumnResponse, len(columns))
	for i, col := range columns {
		tickets := make([]KanbanTicketResponse, len(col.Tickets))
		for j, ticket := range col.Tickets {
			// Convert labels to tags format
			tags := make([]string, len(ticket.Labels))
			for k, label := range ticket.Labels {
				tags[k] = label.Name
			}

			// Convert assignees to string array
			assignees := make([]string, len(ticket.Assignees))
			for k, assignee := range ticket.Assignees {
				assignees[k] = assignee.UserID
			}

			tickets[j] = KanbanTicketResponse{
				ID:          ticket.ID,
				Title:       ticket.Title,
				Description: ticket.Description,
				TicketNo:    ticket.TicketID,
				Tags:        tags,
				Assignees:   assignees,
				Completed:   ticket.Completed,
				DueDate:     ticket.DueDate,
				Priority:    ticket.Priority,
			}
		}

		kanbanColumns[i] = KanbanColumnResponse{
			ID:      col.ID,
			Title:   col.Title,
			Order:   col.Order,
			Tickets: tickets,
		}
	}

	utils.JSON200(c, gin.H{
		"data": kanbanColumns,
	})
}

// CreateColumn tạo column mới
func (cc *ColumnController) CreateColumn(c *gin.Context) {
	var req CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	column := &entity.Column{
		Name:     req.Name,
		Position: req.Position,
	}

	if err := cc.r.GetRepositories().Column.Create(column); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column created successfully",
		"data":    column,
	})
}

// UpdateColumn cập nhật thông tin column
func (cc *ColumnController) UpdateColumn(c *gin.Context) {
	id := c.Param("id")
	var req UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	column, err := cc.r.GetRepositories().Column.GetByID(id)
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

	if err := cc.r.GetRepositories().Column.Update(column); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column updated successfully",
		"data":    column,
	})
}

// DeleteColumn xóa column
func (cc *ColumnController) DeleteColumn(c *gin.Context) {
	id := c.Param("id")

	if err := cc.r.GetRepositories().Column.Delete(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column deleted successfully",
	})
}

// GetColumns lấy danh sách tất cả columns
func (cc *ColumnController) GetColumns(c *gin.Context) {
	columns, err := cc.r.GetRepositories().Column.GetAll()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"data": columns,
	})
}

// UpdateColumnPosition thay đổi vị trí column
func (cc *ColumnController) UpdateColumnPosition(c *gin.Context) {
	id := c.Param("id")
	var req UpdateColumnPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := cc.r.GetRepositories().Column.UpdatePosition(id, req.Position); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Column position updated successfully",
	})
}

// MoveTicketToColumn di chuyển ticket sang column khác
func (cc *ColumnController) MoveTicketToColumn(c *gin.Context) {
	var req MoveTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	if err := cc.r.GetRepositories().Ticket.MoveToColumn(req.TicketID, req.ColumnID); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket moved successfully",
	})
}

// CreateTicket tạo ticket mới trong column
func (cc *ColumnController) CreateTicket(c *gin.Context) {
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

	if err := cc.r.GetRepositories().Ticket.Create(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket created successfully",
		"data":    ticket,
	})
}

// UpdateTicket cập nhật thông tin ticket
func (cc *ColumnController) UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ticket, err := cc.r.GetRepositories().Ticket.GetByID(id)
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

	if err := cc.r.GetRepositories().Ticket.Update(ticket); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket updated successfully",
		"data":    ticket,
	})
}

// DeleteTicket xóa ticket
func (cc *ColumnController) DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	if err := cc.r.GetRepositories().Ticket.Delete(id); err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Ticket deleted successfully",
	})
}

// GetTagColors trả về mapping màu sắc cho các tag
func (cc *ColumnController) GetTagColors(c *gin.Context) {
	labels, err := cc.r.GetRepositories().Label.GetAll()
	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	tagColors := make(map[string]string)
	for _, label := range labels {
		tagColors[label.Name] = label.Color
	}

	utils.JSON200(c, gin.H{
		"data": tagColors,
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
