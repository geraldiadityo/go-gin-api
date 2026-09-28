package pegawai

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, data *Pegawai) error
	FindAll(ctx context.Context) ([]Pegawai, error)
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
