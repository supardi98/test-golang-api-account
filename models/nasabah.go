package models

import (
	"time"
)

type Nasabah struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama       string     `json:"nama"`
	NIK        string     `json:"nik" gorm:"unique"`
	NoHP       string     `json:"no_hp" gorm:"unique"`
	NoRekening string     `json:"no_rekening"`
	Saldo      int64      `json:"saldo" gorm:"default=0,not null"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}
