package models

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type Panitia struct {
	gorm.Model
	Nama  string `gorm:"not null" json:"nama"`
	Tahun int    `gorm:"not null" json:"tahun"`
}

func (obj Panitia) AnggotaList() []Pegawai {
	var anggota []Pegawai
	db.Find(&anggota, "id in (SELECT peg_id FROM anggota_panitia WHERE pnt_id = ? and deleted_at IS NULL) AND peg_isactive=1 AND peg_status IN (1,2)", obj.ID)
	return anggota
}

func (obj Panitia) Anggota() string {
	anggota := obj.AnggotaList()
	if len(anggota) == 0 {
		return ""
	}
	var res string
	for _, v := range anggota {
		if len(res) > 1 {
			res += ","
		}
		res += v.PegNama
	}
	return res
}

func (obj Panitia) AnggotaById(pegId uint) Pegawai {
	var anggota Pegawai
	db.First(&anggota, "id in (SELECT peg_id FROM anggota_panitia WHERE pnt_id = ? AND peg_id=?) and deleted_at IS NULL", obj.ID, pegId)
	return anggota
}

func (obj Panitia) IsAnggota(pegId uint) bool {
	var count int64
	db.Model(&AnggotaPanitia{}).Where("pnt_id=? AND peg_id=? and deleted_at IS NULL", obj.ID, pegId).Count(&count)
	return count > 0
}

func (obj Panitia) SaveAnggotaPokja(panitia PanitiaDTO) error {
	var anggotas []AnggotaPanitia
	for _, v := range panitia.Anggota {
		if !obj.IsAnggota(v)  {
			anggotas = append(anggotas, AnggotaPanitia{
				PegId: v,
				PntId: obj.ID,
			})
		}
	}
	if len(anggotas) > 0 {
		err := db.Save(&anggotas).Error
		if err != nil {
			log.Error(err)
			return errors.New("Gagal Simpan anggota pokja")
		}
	}
	err := db.Where("pnt_id=? AND peg_id NOT IN ? and deleted_at IS NULL", obj.ID, panitia.Anggota).Delete(&AnggotaPanitia{}).Error
	if err != nil {
		log.Error(err)
		return errors.New("Gagal Hapus anggota pokja")
	}
	return nil
}

func GetPanitias() []Panitia {
	var res []Panitia
	db.Find(&res)
	return res
}

func GetPanitia(id uint) Panitia {
	var res Panitia
	db.First(&res, id)
	return res
}

func SavePanitia(panitia *Panitia) error {
	return db.Save(&panitia).Error
}

func DeletePanitia(panitia *Panitia) error {
	return db.Delete(&panitia).Error
}

type AnggotaPanitia struct {
	gorm.Model
	PegId uint `gorm:"not null"`
	PntId uint `gorm:"not null"`
}

func (c AnggotaPanitia) Pegawai() Pegawai {
	var pegawai Pegawai
	db.First(&pegawai, c.PegId)
	return pegawai
}

func GetAnggotaPokja(id uint) []AnggotaPanitia {
	var anggotas []AnggotaPanitia
	db.Find(&anggotas, "pnt_id = ?", id)
	return anggotas
}

func DeleteAnggotaPokja(pntId uint) error {
	return db.Delete(&AnggotaPanitia{}, "pnt_id=?", pntId).Error
}

func IsPegawaiInPanitia(pegId uint, pntId uint) bool {
	var count int64
	db.Model(&AnggotaPanitia{}).Where("pnt_id=? AND peg_id=? and deleted_at IS NULL", pntId, pegId).Count(&count)
	return count > 0
}

type PanitiaDTO struct {
	ID      uint   `form:"id" json:"id"`
	Nama    string `form:"nama" json:"nama"`
	Tahun   int    `form:"tahun" json:"tahun"`
	Anggota []uint `form:"anggota" json:"anggota"`
}
