package bootstrap

import (
	"go-backend-api/internal/api/routes"
	"go-backend-api/internal/config"
	"go-backend-api/internal/logger"
	"go-backend-api/internal/services"

	"github.com/gin-gonic/gin"
)

func InitRouter(services *services.Services, cfg *config.Config) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use()

	// Setup routes
	routes.SetupRoutes(router, services)
	logger.Info("Router initialized successfully")
	return router
}
