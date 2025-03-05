package models

import (
	"time"
)

type TipeTransaksi string

const (
	Tabung TipeTransaksi = "tabung"
	Tarik  TipeTransaksi = "tarik"
)

type Mutasi struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	NasabahID     uint          `json:"nasabah_id"`
	NoRekening    string        `json:"no_rekening"`
	SaldoAwal     int64         `json:"saldo_awal"`
	SaldoAkhir    int64         `json:"saldo_akhir"`
	NominalTabung *int64        `json:"nominal_tabung"`
	NominalTarik  *int64        `json:"nominal_tarik"`
	TipeTransaksi TipeTransaksi `json:"tipe_transaksi"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at" gorm:"index"`
}
