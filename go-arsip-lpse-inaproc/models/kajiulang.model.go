package models

import "gorm.io/gorm"

type KajiUlang struct {
	gorm.Model
	ParentId	uint 		`json:"parent_id" form:"parent_id"`
	PktID		uint 		`json:"pkt_id" form:"pkt_id"`
	PegId		uint 		`json:"peg_id" form:"peg_id"`
	PpkId		uint 		`json:"ppk_id" form:"ppk_id"`
	PntId		uint 		`json:"pnt_id" form:"pnt_id"`
	PpId		uint 		`json:"pp_id" form:"pp_id"`
	Role 		string 		`json:"role" form:"role"`
	Dokumen		string		`json:"dokumen" form:"dokumen"`
	Uraian		string 		`json:"uraian" form:"uraian"`
	DokId		uint		`json:"dok_id" form:"dok_id"`
}

func (s KajiUlang) TableName() string {
	return "kaji_ulang"
}

func (s KajiUlang) Paket() Paket {
	var res Paket
	db.First(&res, s.PktID)
	return res
}

func (s KajiUlang) Parent() KajiUlang {
	var res KajiUlang
	db.First(&res, s.ParentId)
	return res
}

func (s KajiUlang) Ppk() Pegawai {
	var res Pegawai
	db.First(&res, s.PpkId)
	return res
}

func (s KajiUlang) Pp() Pegawai {
	var res Pegawai
	db.First(&res, s.PpId)
	return res
}

func (s KajiUlang) Panitia() Panitia {
	var res Panitia
	db.First(&res, s.PntId)
	return res
}

func (s KajiUlang) DokUpload() Document {
	var res Document
	db.First(&res, s.DokId)
	return res
}

func (s KajiUlang) Pegawai() Pegawai {
	var res Pegawai
	db.First(&res, s.PegId)
	return res
}

func (s KajiUlang) Pengirim() string {
	if s.PpkId != 0 {
		return s.Ppk().PegNama
	}
	if s.PpId != 0 {
		return s.Pp().PegNama
	}
	if s.PntId != 0 {
		return s.Panitia().Nama
	}
	return ""
}

func (s KajiUlang) Penjelasan() []KajiUlang {
	var res []KajiUlang
	db.Find(&res, "parent_id=?", s.ID)
	return res
}

func (s KajiUlang) IsAllowJawab(group string, id uint) bool {
	switch group {
	case ADMIN, UKPBJ:
		return false
	case PPK:
		return s.PpkId != id
	case PP:
		return s.PpId != id
	case POKJA:
		return s.PegId != id
	}
	return false
}

func GetKajiUlangPaket(id uint) []KajiUlang {
	var res []KajiUlang
	db.Find(&res, "pkt_id=? and parent_id = 0", id)
	return res
}

func GetKajiUlang(id uint) KajiUlang {
	var res KajiUlang
	db.First(&res, id)
	return res
}

func SaveKajiUlang(kajiulang *KajiUlang) error {
	return db.Save(kajiulang).Error
}
