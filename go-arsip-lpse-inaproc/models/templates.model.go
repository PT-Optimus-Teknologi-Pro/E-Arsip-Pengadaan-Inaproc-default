package models

import (
	"time"

	"gorm.io/gorm"
)

type Templates struct {
	gorm.Model
	Nama			string 			`form:"nama" json:"nama"`
	Content			string 			`form:"content" json:"context"`
	Variable		string			`form:"variable" json:"variable"`
	PeriodeAwal  	time.Time 		`json:"periode_awal"` // Tanggal penggunaan surat misalkan 1 jan 2025 sd 31 des 2025
	PeriodeAkhir 	time.Time 		`json:"periode_akhir"`
}

func GetTemplates(id uint) Templates {
	var template Templates
	db.First(&template, id)
	return template
}

func CreateTemplates(template Templates) error {
	return db.Create(&template).Error
}

func SaveTemplates(template Templates) error {
	return db.Save(&template).Error
}

func DeleteTemplates(template Templates) error {
	return db.Delete(&template).Error
}
