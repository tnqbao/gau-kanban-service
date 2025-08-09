package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/controller"
)

func SetupRouter(ctrl *controller.Controller) *gin.Engine {

	// Initialize repository manager

	// Initialize kanban controller
	//kanbanCtrl := controller.NewColumnController(repoManager)

	r := gin.Default()

	apiRoutes := r.Group("/api/v2/kanban/")
	{
		apiRoutes.GET("/", ctrl.CheckHealth)

		// Kanban Board endpoints
		apiRoutes.GET("/board", ctrl.GetKanbanBoard)
		apiRoutes.GET("/tag-colors", ctrl.GetTagColors)

		// Column management
		apiRoutes.POST("/columns", ctrl.CreateColumn)
		apiRoutes.GET("/columns", ctrl.GetColumns)
		apiRoutes.PUT("/columns/:id", ctrl.UpdateColumn)
		apiRoutes.DELETE("/columns/:id", ctrl.DeleteColumn)
		apiRoutes.PATCH("/columns/:id/position", ctrl.UpdateColumnPosition)

		// Ticket management
		apiRoutes.POST("/tickets", ctrl.CreateTicket)
		apiRoutes.GET("/tickets", ctrl.GetTickets)
		apiRoutes.GET("/tickets/:id", ctrl.GetTicket)
		apiRoutes.PUT("/tickets/:id", ctrl.UpdateTicket)
		apiRoutes.DELETE("/tickets/:id", ctrl.DeleteTicket)
		apiRoutes.PATCH("/tickets/move", ctrl.MoveTicketToColumn)
		apiRoutes.PATCH("/tickets/move-with-position", ctrl.MoveTicketWithPosition)
		apiRoutes.PATCH("/tickets/:id/position", ctrl.UpdateTicketPosition)

		// Assignment management
		apiRoutes.POST("/assignments", ctrl.CreateAssignment)
		apiRoutes.PUT("/assignments/:id", ctrl.UpdateAssignment)
		apiRoutes.DELETE("/assignments/:id", ctrl.DeleteAssignment)
		apiRoutes.DELETE("/users/:user_id/assignments", ctrl.DeleteAssignmentsByUserID)
		//apiRoutes.GET("/tickets/:ticket_id/assignments", ctrl.GetTicketAssignments)
	}
	return r
}
