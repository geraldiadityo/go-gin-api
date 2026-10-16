package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/redis/go-redis/v9"
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
	*helper.BaseCache
}

func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repository{
		db:        db,
		BaseCache: helper.NewBaseCache(rdb, "role", 10*time.Minute),
	}
}

func (r *repository) Create(ctx context.Context, data *Role) error {
	err := r.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}

	_ = r.InvalidNamespace(ctx)
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]Role, error) {
	return helper.FetchWithCache(ctx, r.BaseCache, "all", func() ([]Role, error) {
		var data []Role
		err := r.db.WithContext(ctx).Find(&data).Error
		return data, err
	})
}

func (r *repository) FindById(ctx context.Context, id uint) (*Role, error) {
	subKey := fmt.Sprintf("id:%d", id)
	return helper.FetchWithCache(ctx, r.BaseCache, subKey, func() (*Role, error) {
		var role Role
		err := r.db.WithContext(ctx).First(&role, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		return &role, nil
	})
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
	err := r.db.WithContext(ctx).Save(data).Error
	if err != nil {
		return err
	}

	_ = r.InvalidNamespace(ctx)
	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&Role{}, id).Error

	if err != nil {
		return err
	}

	_ = r.InvalidNamespace(ctx)
	return nil
}
