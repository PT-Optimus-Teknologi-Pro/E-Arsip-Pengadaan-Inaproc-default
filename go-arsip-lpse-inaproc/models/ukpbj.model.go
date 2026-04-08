package models

import (
	"time"
	"gorm.io/gorm"
)

type Ukpbj struct {
	gorm.Model
	AgcId     uint		`form:"agc_id" json:"agc_id"`
	PegId     uint		`form:"peg_id" json:"peg_id"`
	UpdatedBy string	`form:"update_by" json:"updated_by"`
	Nama      string    `form:"nama" json:"nama"`
	Alamat    string    `form:"alamat" json:"alamat"`
	Telepon   string    `form:"telepon" json:"telepon"`
	Fax       string    `form:"fax" json:"fax"`
	TglDaftar time.Time `form:"tgl_daftar" json:"tgl_daftar"`
	IsActive  bool      `form:"is_active" json:"is_active"`
}

func (Ukpbj) TableName() string {
	return "ukpbj"
}

func (u Ukpbj) GetTglDaftar() string {
	return u.TglDaftar.Format("02-01-2006")
}

func (u Ukpbj) GetAdmnin() Pegawai {
	return GetPegawai(u.PegId)
}


func GetUkpbj(id uint) Ukpbj {
	var ukpbj Ukpbj
	db.First(&ukpbj, id)
	return ukpbj
}

func GetUkpbjAktif() Ukpbj {
	var ukpbj Ukpbj
	db.First(&ukpbj, "is_active=?", true)
	return ukpbj
}

func CreateUkpbj(ukpbj Ukpbj) error {
	return 	db.Create(&ukpbj).Error
}

func SaveUkpbj(ukpbj Ukpbj) error {
	return 	db.Save(&ukpbj).Error
}

func DeleteUkpbj(ukpbj Ukpbj) error {
	return db.Delete(&ukpbj, ukpbj.ID).Error
}

func GetCountUkpbj() int64 {
	var result int64
	db.Model(&Ukpbj{}).Count(&result)
	return result
}
