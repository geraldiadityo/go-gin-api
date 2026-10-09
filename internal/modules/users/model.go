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

type UserUpdateRequest struct {
	Username    string `json:"username"`
	PegawaiID   uint   `json:"pegawaiId"`
	RoleID      uint   `json:"roleId"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	OldPassword string `json:"old_password"`
}

type UserCreateDTO struct {
	Username  string
	Password  string
	PegawaiID uint
	RoleID    uint
}

type UserQueryDTO struct {
	Page             int
	PageSize         int
	keyword          string
	OrderByField     string
	OrderByDirection string
}

type UserResponse struct {
	ID       uint            `json:"id"`
	Username string          `json:"username"`
	Pegawai  pegawai.Pegawai `json:"pegawai"`
	Role     role.Role       `json:"role"`
	Status   bool            `json:"status"`
}
