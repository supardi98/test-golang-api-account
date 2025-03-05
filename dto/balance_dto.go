package dto

import "supardi98/service-account-api/models"

type InputPostTabung struct {
	NoRekening string `json:"no_rekening" validate:"required,number,len=10"`
	Nominal    int64  `json:"nominal" validate:"required,number,min=1"`
}

type ResponsePostTabung struct {
	Saldo int64 `json:"saldo"`
}

type InputPostTarik struct {
	NoRekening string `json:"no_rekening" validate:"required,number,len=10"`
	Nominal    int64  `json:"nominal" validate:"required,number,min=1"`
}

type ResponsePostTarik struct {
	Saldo int64 `json:"saldo"`
}

type ResponseGetSaldoByNoRekening struct {
	Saldo int64 `json:"saldo"`
}

type ResponseGetMutasiByNoRekening []models.Mutasi
