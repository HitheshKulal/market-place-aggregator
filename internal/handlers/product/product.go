package product

import (
	"errors"
	"go-backend-api/internal/handlers"
	"go-backend-api/internal/services"
	"net/http"
	"strings"

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
	templateID := handlers.UintQuery(c, "mappingId", 0)

	if templateID == 0 {
		handlers.BadRequest(c, errors.New("invalid mappingId"))
		return
	}

	instances, err := h.services.Product.GetProductsByTemplate(c.Request.Context(), templateID)
	if err != nil {
		handlers.ErrorResponse(c, err)
	}

	handlers.SuccessResponse(c, instances)
}

func (h *Handler) UploadAndStoreProducts(c *gin.Context) {
	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidCSVFile(header.Filename) {
		handlers.UnprocessableEntityResponse(c, "Invalid file type. Only CSV files are allowed")
		return
	}

	// Parse and store products in one go
	result, err := h.services.Product.UploadAndStoreProducts(c, file, header.Filename)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":           "Products uploaded and stored successfully",
		"fileName":          result.FileName,
		"totalRows":         result.TotalRows,
		"successCount":      result.SuccessCount,
		"failedCount":       result.FailedCount,
		"discoveredColumns": result.DiscoveredColumns,
		"sampleProducts":    result.SampleProducts,
		"errors":            result.Errors,
	})
}

func isValidCSVFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".csv")
}
