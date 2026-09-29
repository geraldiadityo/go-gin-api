package pegawai

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePegawaiDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama tidak boleh kosong",
		})
		return
	}

	result, err := h.service.CreatePegawai(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNamaExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan data pegawai",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat data pegawai",
		"data":    result,
	})
}

func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.service.GetAllPegawai(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal Mengambil data pegawai",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data pegawai",
		"data":    result,
	})
}

func (h *Handler) GetById(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := h.service.GetPegawaiById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil detail pegawai",
		"data":    result,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req UpdatePegawaiDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama tidak boleh kosong"})
		return
	}

	result, err := h.service.UpdatePegawai(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrNamaExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pegawai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui data pegawai",
		"data":    result,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.service.DeletePegawai(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data pegawai"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus data pegawai",
	})
}

// helper
func parseID(c *gin.Context) (uint, error) {
	paramID := c.Param("id")
	id, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		return 0, nil
	}

	return uint(id), nil
}
