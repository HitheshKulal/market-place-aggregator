package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func UintParam(c *gin.Context, param string) (uint, error) {
	value, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func UintQuery(c *gin.Context, param string, defaultValue uint) uint {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		return defaultValue
	}

	return uint(value)
}

func StringQuery(c *gin.Context, param string, defaultValue string) string {
	value := c.Query(param)
	if value == "" {
		return defaultValue
	}
	return value
}
