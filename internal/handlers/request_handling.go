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
