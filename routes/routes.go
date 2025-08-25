package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/controller"
)

func SetupRoutes(ctrl *controller.Controller) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v2/kanban")
	{
		// Column routes
		columns := api.Group("/columns")
		{
			columns.POST("", ctrl.CreateColumn)
			columns.GET("", ctrl.GetColumns)
			columns.PUT("/:id", ctrl.UpdateColumn)
			columns.DELETE("/:id", ctrl.DeleteColumn)
			columns.PUT("/:id/position", ctrl.UpdateColumnPosition)
		}

		// Ticket routes
		tickets := api.Group("/tickets")
		{
			tickets.POST("", ctrl.CreateTicket)
			tickets.GET("", ctrl.GetTickets)
			tickets.GET("/:id", ctrl.GetTicketByID)
			tickets.PUT("/:id", ctrl.UpdateTicket)
			tickets.DELETE("/:id", ctrl.DeleteTicket)

			// Position and movement operations
			tickets.PUT("/:id/position", ctrl.UpdateTicketPosition)
			tickets.PUT("/move", ctrl.MoveTicketToColumn)
			tickets.PUT("/move-with-position", ctrl.MoveTicketWithPosition)
		}

		// Assignment routes
		assignments := api.Group("/assignments")
		{
			assignments.POST("", ctrl.CreateAssignment)
			assignments.GET("/ticket/:ticket_id", ctrl.GetTicketAssignments)
			assignments.PUT("/:id", ctrl.UpdateAssignment)
			assignments.DELETE("/:id", ctrl.DeleteAssignment)
			assignments.DELETE("/user/:user_id", ctrl.DeleteAssignmentsByUserID)
		}

		// Checklist routes
		checklists := api.Group("/checklists")
		{
			checklists.POST("", ctrl.CreateChecklist)
			checklists.GET("/ticket/:ticketId", ctrl.GetChecklistsByTicketID)
			checklists.PUT("/:id", ctrl.UpdateChecklist)
			checklists.PUT("/:id/position", ctrl.UpdateChecklistPosition)
			checklists.DELETE("/:id", ctrl.DeleteChecklist)
		}

		// Kanban board view routes
		kanban := api.Group("/kanban")
		{
			kanban.GET("/board", ctrl.GetKanbanBoard)
		}
	}
	return r
}
