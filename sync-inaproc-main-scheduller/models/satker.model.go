package models

import (
	"log/slog"

	"golang.org/x/exp/maps"
)

type Satker struct {
	KdSatker		uint		`gorm:"primaryKey;autoIncrement:false" json:"kd_satker"`
	NamaSatker		string 		`gorm:"nama_satker" json:"nama_satker"`
	Alamat 			string 		`gorm:"alamat" json:"alamat"`
    Fax				string 		`gorm:"fax" json:"fax"`
    JenisKlpd		string 		`gorm:"jenis_klpd" json:"jenis_klpd"`
    JenisSatker		string 	 	`gorm:"jenis_satker" json:"jenis_satker"`
    KdKlpd			string 		`gorm:"kd_klpd" json:"kd_klpd"`
    KdSatkerStr		string 		`gorm:"kd_satker_str" json:"kd_satker_str"`
    KetSatker		string 		`gorm:"ket_satker" json:"ket_satker"`
    KodeEseleon		string 		`gorm:"kode_eselon" json:"kode_eselon"`
    Kodepos			string 		`gorm:"kodepos" json:"kodepos"`
    NamaKlpd		string 		`gorm:"nama_klpd" json:"nama_klpd"`
    StatusSatker	string 		`gorm:"status_satker" json:"status_satker"`
    TahunAktif		string 		`gorm:"tahun_aktif" json:"tahun_aktif"`
    Telepon			string 		`gorm:"telepon" json:"telepon"`
}

func (c Satker) TableName() string {
	return "inaproc_satker"
}

func SaveSatker(datas *map[uint]Satker) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		err := DB.Save(values).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
