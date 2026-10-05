package pegawai

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
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePegawaiDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "ID Tidak valid")
		return
	}

	result, err := h.service.CreatePegawai(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNamaExists) {
			helper.Error(c, http.StatusNotFound, err.Error())
			return
		}
		helper.Error(c, http.StatusInternalServerError, "Gagal membuat data pegawai")
		return
	}

	helper.Success(c, http.StatusCreated, "berhasil membuat data pegawai baru", result)
}

func (h *Handler) GetAll(c *gin.Context) {
	var query PegawaiQueryDTO
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query.Page = page
	query.PageSize = pageSize
	query.Keyword = c.Query("keyword")
	query.OrderByField = c.Query("orderByField")
	query.OrderByDirection = c.Query("orderByDirection")

	data, meta, err := h.service.GetAllPegawai(c.Request.Context(), query)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil data list pegawai")
		return
	}

	helper.SuccessWithMeta(c, http.StatusOK, "success", data, meta)
}

func (h *Handler) GetById(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	result, err := h.service.GetPegawaiById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			helper.Error(c, http.StatusNotFound, err.Error())
			return
		}
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "success", result)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req UpdatePegawaiDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Nama tidak boleh kosong")
		return
	}

	result, err := h.service.UpdatePegawai(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			helper.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrNamaExists) {
			helper.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		helper.Error(c, http.StatusInternalServerError, "Gagal Update data pegawai")
		return
	}

	helper.Success(c, http.StatusOK, "berhasil mengupdate data pegawai", result)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := h.service.DeletePegawai(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			helper.Error(c, http.StatusNotFound, err.Error())
			return
		}
		helper.Error(c, http.StatusInternalServerError, "Gagal menghapus data pegawai")
	}

	helper.Success(c, http.StatusOK, "Berhasil menghapus data pegawai", true)
}

// helper
// func parseID(c *gin.Context) (uint, error) {
// 	paramID := c.Param("id")
// 	id, err := strconv.ParseUint(paramID, 10, 32)
// 	if err != nil {
// 		return 0, nil
// 	}

// 	return uint(id), nil
// }
