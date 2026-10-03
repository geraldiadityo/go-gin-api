package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context) (uint, error) {
	paramID := c.Param("id")
	id, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
