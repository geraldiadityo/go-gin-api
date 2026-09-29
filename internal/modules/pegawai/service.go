package pegawai

import (
	"context"
	"errors"
)

var (
	ErrNamaExists = errors.New("nama pegawai sudah terdaftar")
	ErrNotFound   = errors.New("data pegawai tidak di temukan")
)

type Service interface {
	CreatePegawai(ctx context.Context, req CreatePegawaiDTO) (Pegawai, error)
	GetAllPegawai(ctx context.Context) ([]Pegawai, error)
	GetPegawaiById(ctx context.Context, id uint) (Pegawai, error)
	UpdatePegawai(ctx context.Context, id uint, req UpdatePegawaiDTO) (Pegawai, error)
	DeletePegawai(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePegawai(ctx context.Context, req CreatePegawaiDTO) (Pegawai, error) {
	existing, err := s.repo.FindByName(ctx, req.Nama)
	if err != nil {
		return Pegawai{}, err
	}

	if existing != nil {
		return Pegawai{}, ErrNamaExists
	}

	pegawai := Pegawai{
		Nama: req.Nama,
	}

	err = s.repo.Create(ctx, &pegawai)
	return pegawai, err
}

func (s *service) GetAllPegawai(ctx context.Context) ([]Pegawai, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetPegawaiById(ctx context.Context, id uint) (Pegawai, error) {
	pegawai, err := s.repo.FindById(ctx, id)
	if err != nil {
		return Pegawai{}, err
	}

	if pegawai == nil {
		return Pegawai{}, ErrNotFound
	}
	return *pegawai, nil
}

func (s *service) UpdatePegawai(ctx context.Context, id uint, req UpdatePegawaiDTO) (Pegawai, error) {
	pegawai, err := s.GetPegawaiById(ctx, id)
	if err != nil {
		return Pegawai{}, err
	}

	if pegawai.Nama != req.Nama {
		existing, err := s.repo.FindByName(ctx, req.Nama)
		if err != nil {
			return Pegawai{}, err
		}
		if existing != nil && existing.ID != id {
			return Pegawai{}, ErrNamaExists
		}
	}

	pegawai.Nama = req.Nama

	err = s.repo.Update(ctx, &pegawai)
	return pegawai, err
}

func (s *service) DeletePegawai(ctx context.Context, id uint) error {
	_, err := s.GetPegawaiById(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
