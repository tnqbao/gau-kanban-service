package controller

import (
	"github.com/tnqbao/gau-kanban-service/config"
	"github.com/tnqbao/gau-kanban-service/infra"
	"github.com/tnqbao/gau-kanban-service/repository"
)

type Controller struct {
	Config     *config.Config
	Infra      *infra.Infra
	Repository *repository.Repository
}

func NewController(config *config.Config, infra *infra.Infra, repo *repository.Repository) *Controller {
	return &Controller{
		Config:     config,
		Infra:      infra,
		Repository: repo,
	}
}
