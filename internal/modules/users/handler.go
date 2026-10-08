package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geraldiadityo/go-backend/internal/helper"
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

func (h *Handler) GetAll(c *gin.Context) {
	var query UserQueryDTO
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query.Page = page
	query.PageSize = pageSize
	query.keyword = c.Query("keyword")
	query.OrderByField = c.Query("orderByField")
	query.OrderByDirection = c.Query("orderByDirection")

	data, meta, err := h.service.GetAllUser(c.Request.Context(), query)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal Mengambil list user")
		return
	}

	helper.SuccessWithMeta(c, http.StatusOK, "succes", data, meta)
}
