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

func (h *Handler) Update(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, "validated error")
		return
	}

	result, err := h.service.UpdateUser(c.Request.Context(), id, req)

	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "berhasil update data pegawai", result)
}

func (h *Handler) GetById(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	data, err := h.service.GetUserById(c.Request.Context(), id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil mengambil data user", data)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil menghapus data user", true)
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	id, err := helper.ParseID(c)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	result, err := h.service.ChangeStatus(c.Request.Context(), id)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil merubah status user", result)
}
