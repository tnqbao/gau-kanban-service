package middlewares

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/config"
)

type CORSConfig struct {
	AllowMethods     []string `json:"allowMethods"`
	AllowHeaders     []string `json:"allowHeaders"`
	ExposeHeaders    []string `json:"exposeHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
	MaxAge           int      `json:"maxAge"`
}

func CORSMiddleware(config *config.EnvConfig) gin.HandlerFunc {
	corsBytes, err := os.ReadFile("config/cors.json")
	if err != nil {
		panic("Failed to read CORS config: " + err.Error())
	}

	var corsConfig CORSConfig
	if err := json.Unmarshal(corsBytes, &corsConfig); err != nil {
		panic("Failed to parse CORS config: " + err.Error())
	}

	//domains := config.CORS.AllowDomains
	//domainList := strings.Split(domains, ",")

	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: false,
		MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
	})
}
