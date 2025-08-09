package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// GetKanbanBoard trả về toàn bộ kanban board với format như yêu cầu
func (ctrl *Controller) GetKanbanBoard(c *gin.Context) {
	columns, err := ctrl.Repository.GetAllColumnWithFullTicketDetails()
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
