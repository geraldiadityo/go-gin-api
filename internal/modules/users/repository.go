package users

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, data *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data *User) error {
	err := r.db.WithContext(ctx).Create(data).Error

	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Preload("Pegawai").Preload("Role").First(data, data.ID).Error
	return nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("LOWER(username) = LOWER(?)", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}
