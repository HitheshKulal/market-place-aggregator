package mapping

import (
	"go-backend-api/internal/api/requests"

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
	products, err := h.services.Mapping.Index(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
	}

	handlers.SuccessResponse(c, products)
}

func (h *Handler) Store(c *gin.Context) {
	var req requests.CreateMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.UnprocessableEntityResponse(c, err.Error())
		return
	}

	mapping, err := h.services.Mapping.Store(c, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
	}

	handlers.SuccessResponse(c, mapping)
}
