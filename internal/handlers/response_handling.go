package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:    true,
		Data:       data,
		StatusCode: http.StatusOK,
	})
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusInternalServerError,
	})
}

func UnprocessableEntityResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success:    false,
		Error:      message,
		StatusCode: http.StatusUnprocessableEntity,
	})
}

func NoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, Response{
		Success:    false,
		StatusCode: http.StatusNoContent,
	})
}

func NotFoundResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusNotFound,
	})
}

func ForbiddenResponse(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusForbidden,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:    true,
		Data:       data,
		StatusCode: http.StatusOK,
	})
}

func AcceptedResponse(c *gin.Context) {
	c.JSON(http.StatusAccepted, nil)
}

func ConflictResponse(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusConflict,
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusBadRequest,
	})
}

func NotImplementedResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotImplemented, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusBadRequest,
	})
}

func UnauthorizedResponse(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Response{
		Success:    false,
		Error:      err.Error(),
		StatusCode: http.StatusUnauthorized,
	})
}

func BytesResponse(c *gin.Context, filename string, contentType string, data []byte) {
	// Auto-detect content type if not provided
	if contentType == "" {
		// 1. Try by file extension
		contentType = mime.TypeByExtension(filepath.Ext(filename))

		// 2. Fallback: detect from data
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}
	}

	contentDisposition := c.DefaultQuery("response-content-disposition", "attachment")

	if contentDisposition != "inline" && contentDisposition != "attachment" {
		contentDisposition = "attachment"
	}

	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=%s", contentDisposition, filename))
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(data)
	c.Abort()
}

func SuccessXMLResponse(c *gin.Context, data interface{}) {
	// Get Accept header from request
	acceptHeader := c.GetHeader("Accept")

	// Default to XML, return JSON only if explicitly requested
	switch {
	case strings.Contains(strings.ToLower(acceptHeader), "application/json"):
		c.JSON(http.StatusOK, Response{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusOK,
		})
	default:
		// Default to XML for all other cases (including no header, */*, application/xml, etc.)
		c.XML(http.StatusOK, Response{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusOK,
		})
	}
}

func ErrorXMLResponse(c *gin.Context, statusCode int, errorMessage interface{}) {
	// Get Accept header from the request
	acceptHeader := c.GetHeader("Accept")

	// Check Accept header and respond accordingly (JSON or XML)
	switch {
	case strings.Contains(strings.ToLower(acceptHeader), "application/json"):
		// If the client expects JSON, return the error response in JSON format
		c.JSON(statusCode, errorMessage)
	default:
		// Default to XML (for application/xml or no header specified)
		c.XML(statusCode, errorMessage)
	}
}
