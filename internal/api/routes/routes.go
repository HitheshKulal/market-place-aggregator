package routes

import (
	"go-backend-api/internal/handlers/health"
	"go-backend-api/internal/handlers/mapping"
	"go-backend-api/internal/handlers/product"
	"go-backend-api/internal/handlers/template"

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
	productHandler := product.NewHandler(services)
	templateHandler := template.NewHandler(services)
	mappingHandler := mapping.NewHandler(services)

	r.GET("/health", healthHandler.Health)

	public := r.Group("/api/v1")
	products := public.Group("products")
	{
		products.GET("/:id", productHandler.Index)
	}

	template := public.Group("templates")
	{
		template.GET("/", templateHandler.Index)
		template.POST("/", templateHandler.Store)
	}

	mapping := public.Group("mappings")
	{
		mapping.GET("/", mappingHandler.Index)
		mapping.POST("/", mappingHandler.Store)
	}

}
