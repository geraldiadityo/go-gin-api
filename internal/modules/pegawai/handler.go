package pegawai

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal Menyimpan data pegawai",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "berhasil membuat pegawai",
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
