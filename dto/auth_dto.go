package dto

type InputPostDaftar struct {
	Nama string `json:"nama" validate:"required"`
	NIK  string `json:"nik" validate:"required,number,len=16"`
	NoHP string `json:"no_hp" validate:"required,number,min=10,max=14"`
}

type ResponsePostDaftar struct {
	NoRekening string `json:"no_rekening"`
}
