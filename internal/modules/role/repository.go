package role

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, data *Role) error
	FindAll(ctx context.Context) ([]Role, error)
	FindById(ctx context.Context, id uint) (*Role, error)
	FindByName(ctx context.Context, nama string) (*Role, error)
	Update(ctx context.Context, data *Role) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data *Role) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *repository) FindAll(ctx context.Context) ([]Role, error) {
	var data []Role
	err := r.db.WithContext(ctx).Find(&data).Error
	return data, err
}

func (r *repository) FindById(ctx context.Context, id uint) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &role, err
}

func (r *repository) FindByName(ctx context.Context, nama string) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).Where("LOWER(nama) = LOWER(?)", nama).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (r *repository) Update(ctx context.Context, data *Role) error {
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}
