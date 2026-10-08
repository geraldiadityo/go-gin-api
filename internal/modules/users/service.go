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
)

type Service interface {
	CreateUser(ctx context.Context, req UserCreateDTO) (*UserResponse, error)
	GetAllUser(ctx context.Context, query UserQueryDTO) ([]UserResponse, helper.Meta, error)
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
