package pegawai

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/geraldiadityo/go-backend/internal/helper"
)

var (
	ErrNamaExists = errors.New("nama pegawai sudah terdaftar")
	ErrNotFound   = errors.New("data pegawai tidak di temukan")
)

type Service interface {
	CreatePegawai(ctx context.Context, req CreatePegawaiDTO) (Pegawai, error)
	GetAllPegawai(ctx context.Context, query PegawaiQueryDTO) ([]Pegawai, helper.Meta, error)
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

func (s *service) GetAllPegawai(ctx context.Context, query PegawaiQueryDTO) ([]Pegawai, helper.Meta, error) {
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
	if strings.TrimSpace(query.Keyword) != "" {
		conditions = append(conditions, helper.Condition{
			Query: "nama ILIKE ?",
			Args:  []interface{}{"%" + query.Keyword + "%"},
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
		return []Pegawai{}, helper.Meta{
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

	return listData, meta, nil
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
