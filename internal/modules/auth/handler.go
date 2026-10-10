package auth

import (
	"errors"
	"net/http"

	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrInactiveUser) {
			helper.Error(c, http.StatusForbidden, err.Error())
			return
		}

		helper.Error(c, http.StatusInternalServerError, "Gagal melakukan login")
		return
	}

	helper.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrInactiveUser) {
			helper.Error(c, http.StatusForbidden, err.Error())
			return
		}

		helper.Error(c, http.StatusInternalServerError, "Gagal melakukan refresh token")
		return
	}

	helper.Success(c, http.StatusOK, "Refresh token berhasil", result)
}

func (h *Handler) Logout(c *gin.Context) {
	// Mengambil user_id dari hasil AuthMiddleware
	userIDVal, exists := c.Get("user_id")
	if !exists {
		helper.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		helper.Error(c, http.StatusInternalServerError, "Invalid user ID in token")
		return
	}

	err := h.service.Logout(c.Request.Context(), userID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal melakukan logout")
		return
	}

	helper.Success[any](c, http.StatusOK, "Logout berhasil", nil)
}
