package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/controller"
	"github.com/tnqbao/gau-kanban-service/middlewares"
)

func SetupRoutes(ctrl *controller.Controller) *gin.Engine {
	router := gin.Default()

	middleware, err := middlewares.NewMiddlewares(ctrl)
	if err != nil {
		panic("Failed to initialize middlewares: " + err.Error())
	}
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.AuthMiddleware)

	// API version 2
	v2 := router.Group("/api/v2")
	{
		// Kanban routes
		kanban := v2.Group("/kanban")
		{
			// Board routes
			boards := kanban.Group("/boards")
			{
				boards.POST("", ctrl.CreateBoard)
				boards.GET("", ctrl.GetBoards)
				boards.GET("/:id", ctrl.GetBoardByID)
				boards.GET("/:id/search", ctrl.SearchTickets)
				boards.GET("/:id/filter", ctrl.FilterTickets)
				boards.PATCH("/:id", ctrl.UpdateBoard)
				boards.DELETE("/:id", ctrl.DeleteBoard)
				boards.PATCH("/:id/archive", ctrl.ArchiveBoard)
			}

			// Column routes
			columns := kanban.Group("/columns")
			{
				columns.POST("", ctrl.CreateColumn)
				columns.GET("/:id", ctrl.GetColumnByID)
				columns.PATCH("/:id", ctrl.UpdateColumn)
				columns.DELETE("/:id", ctrl.DeleteColumn)
				columns.PATCH("/:id/reorder", ctrl.ReorderColumn)
			}

			// Ticket routes
			tickets := kanban.Group("/tickets")
			{
				tickets.POST("", ctrl.CreateTicket)
				tickets.GET("/:id", ctrl.GetTicketByID)
				tickets.PATCH("/:id", ctrl.UpdateTicket)
				tickets.PATCH("/:id/move", ctrl.MoveTicket)
				tickets.DELETE("/:id", ctrl.DeleteTicket)
			}

			// User routes
			users := kanban.Group("/users")
			{
				users.POST("", ctrl.CreateUser)
			}

			// Member routes
			members := kanban.Group("/members")
			{
				members.POST("", ctrl.AddMemberToBoard)
			}

			// Label routes
			labels := kanban.Group("/labels")
			{
				labels.POST("", ctrl.CreateLabel)
				labels.POST("/assign", ctrl.CreateLabelTicket)
			}

			// Assignee routes
			assignees := kanban.Group("/assignees")
			{
				assignees.POST("", ctrl.CreateAssignee)
				assignees.DELETE("", ctrl.DeleteAssignee)
			}

			// Checklist routes
			checklists := kanban.Group("/checklists")
			{
				checklists.POST("", ctrl.CreateChecklist)
				checklists.PATCH("/:id", ctrl.UpdateChecklist)
				checklists.DELETE("/:id", ctrl.DeleteChecklist)
			}
		}
	}

	return router
}
