package role

import (
	"context"
	"errors"
)

var (
	ErrNamaExists = errors.New("Nama role sudah terdaftar")
	ErrNotFound   = errors.New("Data role tidak di temukan")
)

type Service interface {
	CreateRole(ctx context.Context, req CreateRoleDTO) (Role, error)
	GetAllRole(ctx context.Context) ([]Role, error)
	GetRoleById(ctx context.Context, id uint) (Role, error)
	UpdateRole(ctx context.Context, id uint, req UpdateRoleDTO) (Role, error)
	DeleteRole(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRole(ctx context.Context, req CreateRoleDTO) (Role, error) {
	existing, err := s.repo.FindByName(ctx, req.Nama)
	if err != nil {
		return Role{}, err
	}

	if existing != nil {
		return Role{}, ErrNamaExists
	}

	role := Role{
		Nama: req.Nama,
	}

	err = s.repo.Create(ctx, &role)
	return role, err
}

func (s *service) GetAllRole(ctx context.Context) ([]Role, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetRoleById(ctx context.Context, id uint) (Role, error) {
	role, err := s.repo.FindById(ctx, id)
	if err != nil {
		return Role{}, err
	}

	if role == nil {
		return Role{}, ErrNotFound
	}

	return *role, nil
}

func (s *service) UpdateRole(ctx context.Context, id uint, req UpdateRoleDTO) (Role, error) {
	role, err := s.GetRoleById(ctx, id)

	if err != nil {
		return Role{}, err
	}

	if role.Nama != req.Nama {
		existing, err := s.repo.FindByName(ctx, req.Nama)
		if err != nil {
			return Role{}, err
		}

		if existing != nil && existing.ID != id {
			return Role{}, ErrNamaExists
		}
	}

	role.Nama = req.Nama
	err = s.repo.Update(ctx, &role)
	return role, err
}

func (s *service) DeleteRole(ctx context.Context, id uint) error {
	_, err := s.GetRoleById(ctx, id)

	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
