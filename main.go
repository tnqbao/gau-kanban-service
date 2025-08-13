package main

import (
	"log"

	"github.com/tnqbao/gau-kanban-service/controller"
	"github.com/tnqbao/gau-kanban-service/infra"
	"github.com/tnqbao/gau-kanban-service/repository"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-kanban-service/config"
	"github.com/tnqbao/gau-kanban-service/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	inf := infra.InitInfra(cfg)
	if inf == nil {
		log.Fatal("Failed to initialize infrastructure")
	}

	repo := repository.InitRepository(inf)
	if repo == nil {
		log.Fatal("Failed to initialize repository")
	}

	ctrl := controller.NewController(cfg, inf, repo)
	if ctrl == nil {
		log.Fatal("Failed to initialize controller")
	}

	router := routes.SetupRoutes(ctrl)
	if router == nil {
		log.Fatal("Failed to set up router")
	}
	err = router.Run(":8080")
	if err != nil {
		return
	}
}
