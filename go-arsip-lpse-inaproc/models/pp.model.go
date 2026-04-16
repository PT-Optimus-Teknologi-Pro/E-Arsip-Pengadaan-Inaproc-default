package models

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type PejabatPengadaan struct {
	gorm.Model
	Groups 			string 		`json:"groups"`
	Tahun			int 		`json:"tahun"`
	NoSk			string 		`json:"no_sk"`
	TglSk			time.Time	`json:"tgl_sk"`
	TempatSk		string 		`json:"tempat_sk"`
	PeriodeAwal		time.Time	`json:"periode_awal"`
	PeriodeAkhir	time.Time	`json:"periode_akhir"`
	UkpbjId			uint 		`json:"ukpbj_id"`
}

func (PejabatPengadaan) TableName() string {
	return "pejabat_pengadaan"
}

func (obj PejabatPengadaan) Pegawai() []Pegawai {
	var res []Pegawai
	db.Find(&res, "id IN (SELECT peg_id FROM pejabat_pengadaan_pegawai WHERE pp_id =?)", obj.ID)
	return res
}

func (obj PejabatPengadaan) IsPegawaiInGroup(pegId uint) bool {
	var count int64
	db.Model(&PejabatPengadaanPegawai{}).Where("pp_id=? AND peg_id=?", obj.ID, pegId).Count(&count)
	return count > 0
}

func (obj PejabatPengadaan) IsSatkerInGroup(satkerId uint) bool {
	var count int64
	db.Model(&PejabatPengadaanSatker{}).Where("pp_id=? AND satker_id=?", obj.ID, satkerId).Count(&count)
	return count > 0
}

func (obj PejabatPengadaan) Satker() []APISatkerSirup {
	var satkers []APISatkerSirup
	db.Select([]string{"id", "nama", "id_kldi"}).Find(&satkers, "id IN (SELECT satker_id FROM pejabat_pengadaan_satker WHERE pp_id =?)", obj.ID)
	return satkers
}

func (obj PejabatPengadaan) SavePejabatPengadaanSatker(satkers *[]PejabatPengadaanSatker) error {
	if len(*satkers) > 0 {
		err := db.Save(satkers).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
		var ids []uint
		for _, v := range *satkers {
			ids = append(ids, v.SatkerId)
		}
		err = db.Where("pp_id = ? AND satker_id NOT IN (?)", obj.ID, ids).Delete(&PejabatPengadaanSatker{}).Error
		if err != nil {
			log.Error(err)
			return errors.New("Gagal Hapus satker di Pejabat Pengadaan")
		}
	}
	return nil
}

func (obj PejabatPengadaan) SavePejabatPengadaanPegawai(dto []uint) error {
	if len(dto) > 0 {
		var pegawais []PejabatPengadaanPegawai
		for _, v := range dto {
			objPegawai := PejabatPengadaanPegawai {
				PpId: obj.ID,
				PegId: v,
			}
			pegawais = append(pegawais, objPegawai)
		}
		err := db.Save(&pegawais).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = db.Where("pp_id = ? AND peg_id NOT IN ?", obj.ID, dto).Delete(&PejabatPengadaanPegawai{}).Error
		if err != nil {
			log.Error(err)
			return errors.New("Gagal Hapus Pegawai di Pejabat Pengadaan")
		}
	}
	return nil
}

func GetPejabatPengadaan(id uint) PejabatPengadaan {
	var res PejabatPengadaan
	db.First(&res, id)
	return res
}

func DeletePejabatPengadaan(obj *PejabatPengadaan) error {
	return db.Delete(obj).Error
}

func SavePejabatPengadaan(obj *PejabatPengadaan) error {
	return db.Save(obj).Error
}


type PejabatPengadaanSatker struct {
	PpId		uint 		`gorm:"primaryKey"`
	SatkerId	uint		`gorm:"primaryKey"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

func (PejabatPengadaanSatker) TableName() string {
	return "pejabat_pengadaan_satker"
}

func DeletePejabatPengadaanSatker(ppId uint) error {
	return db.Delete(&PejabatPengadaanSatker{}, "pp_id=?", ppId).Error
}

func SavePejabatPengadaanSatker(obj *PejabatPengadaanSatker) error {
	return db.Save(obj).Error
}

type PejabatPengadaanPegawai struct {
	PpId		uint 		`gorm:"primaryKey"`
	PegId		uint		`gorm:"primaryKey"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

func (PejabatPengadaanPegawai) TableName() string {
	return "pejabat_pengadaan_pegawai"
}

func DeletePejabatPengadaanPegawai(ppId uint) error {
	return db.Delete(&PejabatPengadaanPegawai{}, "pp_id=?", ppId).Error
}

type PejabatPengadaanDTO struct {
	ID 				uint 		`json:"id" form:"ID"`
	Groups			string 		`json:"groups" form:"groups"`
	Tahun			int 		`json:"tahun" form:"tahun"`
	NoSk			string 		`json:"no_sk" form:"no_sk"`
	TglSk			time.Time	`json:"tgl_sk" form:"tgl_sk"`
	TempatSk		string 		`json:"tempat_sk" form:"tempat_sk"`
	PeriodeAwal		time.Time	`json:"periode_awal" form:"periode_awal"`
	PeriodeAkhir	time.Time	`json:"periode_akhir" form:"periode_akhir"`
	Satker 			[]uint      `json:"satker" form:"satker"`
	Pegawai 		[]uint 		`json:"pegawai" form:"pegawai"`
}

func GetPejabatPengadaanSatker(satkerId uint) PejabatPengadaanSatker {
	var res PejabatPengadaanSatker
	db.Find(&res, "satker_id=?", satkerId)
	return res
}
