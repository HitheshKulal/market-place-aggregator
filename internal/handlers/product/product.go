package product

import (
	"go-backend-api/internal/handlers"
	"go-backend-api/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Index(c *gin.Context) {
	instances, err := h.services.Product.GetProductsByTemplate(c.Request.Context(), 1)
	if err != nil {
		handlers.ErrorResponse(c, err)
	}

	handlers.SuccessResponse(c, instances)
}
