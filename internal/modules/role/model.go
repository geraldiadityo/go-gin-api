package role

import "time"

type Role struct {
	ID        uint      `gorm:"privateKey" json:"id"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRoleDTO struct {
	Nama string `json:"nama" binding:"required"`
}

type UpdateRoleDTO struct {
	Nama string `json:"nama" binding:"required"`
}
