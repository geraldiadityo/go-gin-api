package users

import (
	"context"
	"errors"

	"github.com/geraldiadityo/go-backend/internal/helper"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, data *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindById(ctx context.Context, id uint) (*User, error)
	FindAllWithOptions(ctx context.Context, opts helper.QueryOptions) ([]User, error)
	CountAll(ctx context.Context, conditions []helper.Condition) (int64, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, id uint) error
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

func (r *repository) FindAllWithOptions(ctx context.Context, opts helper.QueryOptions) ([]User, error) {
	var data []User
	query := r.db.WithContext(ctx).Model(&User{}).Preload("Pegawai").Preload("Role")
	for _, cond := range opts.Where {
		query = query.Where(cond.Query, cond.Args...)
	}

	if opts.OrderBy != "" {
		query = query.Order(opts.OrderBy)
	}

	if opts.Take > 0 {
		query = query.Limit(opts.Take)
	}

	if opts.Skip > 0 {
		query = query.Offset(opts.Skip)
	}

	err := query.Find(&data).Error

	return data, err
}

func (r *repository) CountAll(ctx context.Context, conditions []helper.Condition) (int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&User{})

	for _, cond := range conditions {
		query = query.Where(cond.Query, cond.Args...)
	}

	err := query.Count(&total).Error
	return total, err
}

func (r *repository) FindById(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Preload("Pegawai").Preload("Role").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, data *User) error {
	// Gunakan Omit("Pegawai", "Role") agar GORM tidak me-reset 
	// PegawaiID dan RoleID ke data struct relasi yang di-preload sebelumnya.
	err := r.db.WithContext(ctx).Omit("Pegawai", "Role").Save(data).Error

	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Preload("Pegawai").Preload("Role").First(data, data.ID).Error
	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}
