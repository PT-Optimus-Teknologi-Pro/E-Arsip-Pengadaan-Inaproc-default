package models

import (
	"gorm.io/gorm"
	"time"
)

type DokTemplate struct {
	gorm.Model
	Jenis        string    `gorm:"not null"` // Penunjukan Pejabat Pengadaan ,Penunjukan  Pokja Pemilihan, Berita Acara Negosiasi dll
	Keterangan   string    `json:"keterangan"`
	PeriodeAwal  time.Time `json:"periode_awal"` // Tanggal penggunaan surat misalkan 1 jan 2025 sd 31 des 2025
	PeriodeAkhir time.Time `json:"periode_akhir"`
	DokId        uint      `json:"dok_id"` // master doc template id, biasany docx
}

func (c DokTemplate) Dokumen() Document {
	var dokumen Document
	if c.DokId > 0 {
		db.First(&dokumen, c.DokId)
	}
	return dokumen
}

func (c DokTemplate) IsInChecklist(list []ChecklistDok) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if c.ID == v.DokId {
			return true;
		}
	}
	return false;
}

func GetAllDokTemplate() []DokTemplate {
	var doktemplates []DokTemplate
	db.Find(&doktemplates)
	return doktemplates
}

func GetDokTemplate(id uint) DokTemplate {
	var res DokTemplate
	db.First(&res, id)
	return res
}

func SaveDocTemplate(template DokTemplate) error {
	return db.Save(&template).Error
}

func DeleteDocTemplate(template DokTemplate) error {
	return db.Delete(&template).Error
}
