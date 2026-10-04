package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := UserCreateDTO{
		Username:  req.Username,
		Password:  req.Password,
		PegawaiID: req.PegawaiID,
		RoleID:    req.RoleID,
	}

	result, err := h.service.CreateUser(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrUsernameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat User baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat data user",
		"data":    result,
	})
}
