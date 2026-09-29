package pegawai

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, data *Pegawai) error
	FindAll(ctx context.Context) ([]Pegawai, error)
	FindById(ctx context.Context, id uint) (*Pegawai, error)
	FindByName(ctx context.Context, nama string) (*Pegawai, error)
	Update(ctx context.Context, data *Pegawai) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data *Pegawai) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *repository) FindAll(ctx context.Context) ([]Pegawai, error) {
	var data []Pegawai
	err := r.db.WithContext(ctx).Find(&data).Error
	return data, err
}

func (r *repository) FindById(ctx context.Context, id uint) (*Pegawai, error) {
	var pegawai Pegawai
	err := r.db.WithContext(ctx).First(&pegawai, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pegawai, err
}

func (r *repository) FindByName(ctx context.Context, nama string) (*Pegawai, error) {
	var pegawai Pegawai
	err := r.db.WithContext(ctx).Where("LOWER(nama) = LOWER(?)", nama).First(&pegawai).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pegawai, err
}

func (r *repository) Update(ctx context.Context, data *Pegawai) error {
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Pegawai{}, id).Error
}
