package models

import (
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	gorm.Model
	UpdatedBy    string    `form:"update_by" json:"updated_by"`
	AgcNama      string    `form:"agc_nama" json:"agc_nama"`
	AgcAlamat    string    `form:"agc_alamat" json:"agc_alamat"`
	AgcTelepon   string    `form:"agc_telepon" json:"agc_telepon"`
	AgcFax       string    `form:"agc_fax" json:"agc_fax"`
	AgcWebsite   string    `form:"agc_website" json:"agc_website"`
	AgcTglDaftar time.Time `form:"agc_tgl_daftar" json:"agc_tgl_daftar"`
	InstansiId   string    `form:"instansi_id" json:"instansi_id"`
}

func (Agency) TableName() string {
	return "agency"
}

func (u Agency) GetTglDaftar() string {
	return u.AgcTglDaftar.Format("02-01-2006")
}


func GetAgency(id uint) Agency {
	var agency Agency
	db.First(&agency, id)
	return agency
}

func CreateAgency(agency Agency) error {
	return db.Create(&agency).Error
}

func SaveAgency(agency Agency) error {
	return db.Save(&agency).Error
}

func DeleteAgency(agency Agency) error {
	return db.Delete(&agency).Error
}

func GetCountAgency() int64 {
	var result int64
	db.Model(&Agency{}).Count(&result)
	return result
}
