package users

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
	"github.com/geraldiadityo/go-backend/internal/modules/role"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists = errors.New("Username sudah terdaftar")
	ErrNotFound       = errors.New("user is not found")
)

type Service interface {
	CreateUser(ctx context.Context, req UserCreateDTO) (*UserResponse, error)
	GetAllUser(ctx context.Context, query UserQueryDTO) ([]UserResponse, helper.Meta, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetById(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, id uint, req UserUpdateRequest) (*UserResponse, error)
	GetUserById(ctx context.Context, id uint) (*UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	ChangeStatus(ctx context.Context, id uint) (*UserResponse, error)
}

type service struct {
	repo           Repository
	pegawaiService pegawai.Service
	roleService    role.Service
}

func NewService(repo Repository, ps pegawai.Service, rs role.Service) Service {
	return &service{
		repo:           repo,
		pegawaiService: ps,
		roleService:    rs,
	}
}

func (s *service) toUserResponse(data User) UserResponse {
	return UserResponse{
		ID:       data.ID,
		Username: data.Username,
		Pegawai:  data.Pegawai,
		Role:     data.Role,
		Status:   data.Status,
	}
}

func (s *service) CreateUser(ctx context.Context, req UserCreateDTO) (*UserResponse, error) {
	pegawaiData, err := s.pegawaiService.GetPegawaiById(ctx, req.PegawaiID)
	if err != nil {
		return nil, errors.New("Pegawai tidak ditemukan")
	}

	roleData, err := s.roleService.GetRoleById(ctx, req.RoleID)
	if err != nil {
		return nil, errors.New("Role tidak ditemukan")
	}

	existing, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		PegawaiID: pegawaiData.ID,
		RoleID:    roleData.ID,
		Status:    true,
	}

	err = s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*newUser)

	return &response, nil
}

func (s *service) GetAllUser(ctx context.Context, query UserQueryDTO) ([]UserResponse, helper.Meta, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = 10
	}

	skip := (query.Page - 1) * query.PageSize
	take := query.PageSize

	orderBy := "id desc"

	if query.OrderByField != "" {
		dir := "asc"
		if strings.ToLower(query.OrderByDirection) == "desc" || query.OrderByDirection == "-1" {
			dir = "desc"
		}
		orderBy = fmt.Sprintf("%s %s", query.OrderByField, dir)
	}

	var conditions []helper.Condition
	if strings.TrimSpace(query.keyword) != "" {
		conditions = append(conditions, helper.Condition{
			Query: "username ILIKE ?",
			Args:  []interface{}{"%" + query.keyword + "%"},
		})
	}

	opts := helper.QueryOptions{
		Where:   conditions,
		OrderBy: orderBy,
		Take:    take,
		Skip:    skip,
	}

	listData, err := s.repo.FindAllWithOptions(ctx, opts)

	if err != nil {
		return nil, helper.Meta{}, err
	}

	totalItem, err := s.repo.CountAll(ctx, opts.Where)

	if err != nil {
		return nil, helper.Meta{}, err
	}

	if len(listData) == 0 {
		return []UserResponse{}, helper.Meta{
			TotalItem:   0,
			TotalPage:   0,
			CurrentPage: query.Page,
		}, nil
	}

	totalPage := int(math.Ceil(float64(totalItem) / float64(query.PageSize)))
	meta := helper.Meta{
		TotalItem:   totalItem,
		TotalPage:   totalPage,
		CurrentPage: query.Page,
	}

	var reponses []UserResponse
	for _, user := range listData {
		reponses = append(reponses, s.toUserResponse(user))
	}

	return reponses, meta, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.repo.FindByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	return user, err
}

func (s *service) GetById(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	return user, err
}

func (s *service) UpdateUser(ctx context.Context, id uint, req UserUpdateRequest) (*UserResponse, error) {
	user, err := s.GetById(ctx, id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// 1. Update Username (jika dikirim & berbeda)
	if req.Username != "" && user.Username != req.Username {
		existing, err := s.repo.FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrUsernameExists
		}
		user.Username = req.Username
	}

	// 2. Update Pegawai (jika dikirim)
	if req.PegawaiID != 0 {
		pegawai, err := s.pegawaiService.GetPegawaiById(ctx, req.PegawaiID)
		if err != nil {
			return nil, errors.New("Pegawai tidak di temukan")
		}
		user.PegawaiID = pegawai.ID
	}

	// 3. Update Role (jika dikirim)
	if req.RoleID != 0 {
		role, err := s.roleService.GetRoleById(ctx, req.RoleID)
		if err != nil {
			return nil, errors.New("Role tidak ditemukan")
		}
		user.RoleID = role.ID
	}

	// 4. Update Password (jika dikirim)
	if req.Password != "" {
		if req.OldPassword == "" {
			return nil, errors.New("password lama harus di isi jika ingin merubah password")
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
		if err != nil {
			return nil, errors.New("password lama salah")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.Password = string(hashedPassword)
	}

	// 5. Simpan perubahan (WAJIB tangani error dari repo)
	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*user)

	return &response, nil
}

func (s *service) GetUserById(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.GetById(ctx, id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	response := s.toUserResponse(*user)
	return &response, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	_, err := s.GetById(ctx, id)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	err = s.repo.Delete(ctx, id)
	return err
}

func (s *service) ChangeStatus(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.GetById(ctx, id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// Ubah status menjadi kebalikannya
	user.Status = !user.Status

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*user)
	return &response, nil
}
