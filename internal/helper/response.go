package helper

import (
	"github.com/gin-gonic/gin"
)

// Definisikan ApiResponse dengan Generic Type [T any]
type ApiResponse[T any] struct {
	Message     string  `json:"message,omitempty"`
	Data        T       `json:"data,omitempty"`
	Error       string  `json:"error,omitempty"`
	TotalItem   *int    `json:"totalItem,omitempty"`
	TotalPage   *int    `json:"totalPage,omitempty"`
	CurrentPage *int    `json:"currentPage,omitempty"`
	Token       *string `json:"token,omitempty"`
}

// Helper untuk HTTP Success (2xx)
func Success[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, ApiResponse[T]{
		Message: message,
		Data:    data,
	})
}

// Helper untuk HTTP Error (4xx, 5xx)
func Error(c *gin.Context, statusCode int, errMessage string) {
	c.JSON(statusCode, ApiResponse[any]{
		Error: errMessage,
	})
}

func SuccessWithMeta[T any](c *gin.Context, statusCode int, message string, data T, meta Meta) {
	totalItemInt := int(meta.TotalItem)
	c.JSON(statusCode, ApiResponse[T]{
		Message:     message,
		Data:        data,
		TotalItem:   &totalItemInt,
		TotalPage:   &meta.TotalPage,
		CurrentPage: &meta.CurrentPage,
	})
}
