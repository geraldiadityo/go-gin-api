package users

import (
	"context"
	"errors"

	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
	"github.com/geraldiadityo/go-backend/internal/modules/role"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists = errors.New("Username sudah terdaftar")
)

type Service interface {
	CreateUser(ctx context.Context, req UserCreateDTO) (*User, error)
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

func (s *service) CreateUser(ctx context.Context, req UserCreateDTO) (*User, error) {
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

	return newUser, nil
}
