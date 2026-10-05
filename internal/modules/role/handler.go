package role

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

func (h *Handler) Create(c *gin.Context) {
	var req CreateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Nama tidak boleh kosong")
		return
	}

	result, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNamaExists) {
			helper.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		helper.Error(c, http.StatusInternalServerError, "Gagal Membuat data role")
		return
	}

	helper.Success(c, http.StatusCreated, "berhasil membuat data role", result)
}

func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.service.GetAllRole(c.Request.Context())
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil data list role")
		return
	}

	helper.Success(c, http.StatusOK, "success", result)
}

func (h *Handler) GetById(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
	}

	result, err := h.service.GetRoleById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil detail role",
		"data":    result,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req UpdateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama tidak boleh kosong"})
		return
	}

	result, err := h.service.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrNamaExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui data role",
		"data":    result,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus data role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus data role",
	})
}
