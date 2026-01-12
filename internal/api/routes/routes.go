package routes

import (
	"go-backend-api/internal/handlers/health"
	service "go-backend-api/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func SetupRoutes(r *gin.Engine, services *service.Services) {
	healthHandler := health.NewHandler()

	r.GET("/health", healthHandler.Health)
}
