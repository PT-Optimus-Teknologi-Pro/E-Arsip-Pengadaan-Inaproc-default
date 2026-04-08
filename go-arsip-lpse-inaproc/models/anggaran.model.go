package models

import "gorm.io/gorm"

type Anggaran struct {
	gorm.Model
	KodeRekening string  `json:"kode_rekening" form:"kode_rekening"`
	Nilai        float64 `json:"nilai" form:"nilai"`
	Tahun        int64   `json:"tahun" form:"tahun"`
	Uraian       string  `json:"uraian" form:"uraian"`
	Sumber       string  `json:"sumber" form:"sumber"`
	StkId        string  `json:"stk_id" form:"stk_id"`
}
