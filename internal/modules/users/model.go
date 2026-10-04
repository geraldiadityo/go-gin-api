package users

import (
	"time"

	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
	"github.com/geraldiadityo/go-backend/internal/modules/role"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	PegawaiID uint      `gorm:"not null;column:pegawaiId" json:"pegawaiId"`
	RoleID    uint      `gorm:"not null;column:roleId" json:"roleId"`
	Status    bool      `gorm:"default:true" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`

	Pegawai pegawai.Pegawai `gorm:"foreignKey:PegawaiID" json:"pegawai"`
	Role    role.Role       `gorm:"foreignKey:RoleID" json:"role"`
}

type UserCreateRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	PegawaiID       uint   `json:"pegawaiId" binding:"required"`
	RoleID          uint   `json:"roleId" binding:"required"`
}

type UserCreateDTO struct {
	Username  string
	Password  string
	PegawaiID uint
	RoleID    uint
}
