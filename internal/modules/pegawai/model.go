package pegawai

import "time"

type Pegawai struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePegawaiDTO struct {
	Nama string `json:"nama" binding:"required"`
}

type UpdatePegawaiDTO struct {
	Nama string `json:"nama" binding:"required"`
}

type PegawaiQueryDTO struct {
	Page             int
	PageSize         int
	Keyword          string
	OrderByField     string
	OrderByDirection string
}
