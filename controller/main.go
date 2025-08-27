package controller

import (
	"github.com/tnqbao/gau-kanban-service/config"
	"github.com/tnqbao/gau-kanban-service/infra"
	"github.com/tnqbao/gau-kanban-service/provider"
	"github.com/tnqbao/gau-kanban-service/repository"
)

type Controller struct {
	Config     *config.Config
	Infra      *infra.Infra
	Repository *repository.Repository
	Provider   *provider.Provider
}

func NewController(config *config.Config, infra *infra.Infra, repo *repository.Repository) *Controller {
	provide := provider.InitProvider(config.EnvConfig)
	return &Controller{
		Config:     config,
		Infra:      infra,
		Repository: repo,
		Provider:   provide,
	}
}
