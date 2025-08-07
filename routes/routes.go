package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/config"
	"github.com/tnqbao/gau-kanban-service/controller"
	"github.com/tnqbao/gau-kanban-service/infra"
	"github.com/tnqbao/gau-kanban-service/repository"
)

func SetupRouter(config *config.Config) *gin.Engine {
	inf := infra.InitInfra(config)
	ctrl := controller.NewController(config, inf)

	// Initialize repository manager
	repoManager := repository.NewRepositoryManager(inf.Postgres.DB)

	// Initialize kanban controller
	kanbanCtrl := controller.NewColumnController(repoManager)

	r := gin.Default()

	apiRoutes := r.Group("/api/v2/kanban/")
	{
		apiRoutes.GET("/", ctrl.CheckHealth)

		// Kanban Board endpoints
		apiRoutes.GET("/board", kanbanCtrl.GetKanbanBoard)
		apiRoutes.GET("/tag-colors", kanbanCtrl.GetTagColors)

		// Column management
		apiRoutes.POST("/columns", kanbanCtrl.CreateColumn)
		apiRoutes.GET("/columns", kanbanCtrl.GetColumns)
		apiRoutes.PUT("/columns/:id", kanbanCtrl.UpdateColumn)
		apiRoutes.DELETE("/columns/:id", kanbanCtrl.DeleteColumn)
		apiRoutes.PATCH("/columns/:id/position", kanbanCtrl.UpdateColumnPosition)

		// Ticket management
		apiRoutes.POST("/tickets", kanbanCtrl.CreateTicket)
		apiRoutes.PUT("/tickets/:id", kanbanCtrl.UpdateTicket)
		apiRoutes.DELETE("/tickets/:id", kanbanCtrl.DeleteTicket)
		apiRoutes.PATCH("/tickets/move", kanbanCtrl.MoveTicketToColumn)
	}
	return r
}
