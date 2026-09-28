package pegawai

import "context"

type Service interface {
	CreatePegawai(ctx context.Context, req CreatePegawaiDTO) (Pegawai, error)
	GetAllPegawai(ctx context.Context) ([]Pegawai, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePegawai(ctx context.Context, req CreatePegawaiDTO) (Pegawai, error) {
	pegawai := Pegawai{
		Nama: req.Nama,
	}

	err := s.repo.Create(ctx, &pegawai)

	return pegawai, err
}

func (s *service) GetAllPegawai(ctx context.Context) ([]Pegawai, error) {
	return s.repo.FindAll(ctx)
}
