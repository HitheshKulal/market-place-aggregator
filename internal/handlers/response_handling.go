package handlers

import (
	"net/http"

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
