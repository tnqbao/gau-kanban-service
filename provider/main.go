package provider

import (
	"github.com/tnqbao/gau-kanban-service/config"
)

type Provider struct {
	LoggerProvider               *LoggerProvider
	AuthorizationServiceProvider *AuthorizationServiceProvider
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig) *Provider {
	loggerProvider := NewLoggerProvider()
	authorizationServiceProvider := NewAuthorizationServiceProvider(cfg)
	provider = &Provider{
		LoggerProvider:               loggerProvider,
		AuthorizationServiceProvider: authorizationServiceProvider,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
