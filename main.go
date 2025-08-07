package main

import (
	"log"

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

	router := routes.SetupRouter(cfg)
	router.Run(":8080")
}
