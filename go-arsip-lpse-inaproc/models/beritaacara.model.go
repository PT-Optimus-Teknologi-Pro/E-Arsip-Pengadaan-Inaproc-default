package models

import (
	"database/sql"
	"gorm.io/gorm"
)

var BidangList = []string {
	"Anggaran", "Rencana Umum Pengadaan (RUP)",  "Spesifikasi", "Gambar (Barang)/ Gambar Kerja (Konstruksi)",
    "KAK", "Metode pengadaan", "HPS", "Rancangan Kontrak", "Analisa Pasar",
    "Jaminan penawaran untuk A KOMSTELAR dengan nilai total HPS > 10M",
    "Sertifikat garansi/kartu jaminan/garansi purnajual untuk pengadaan barang",
    "Sertifikat/dokumen dalam rangka pengadaan barang impor untuk pengadaan barang",
}

type BeritaAcara struct {
	gorm.Model
	PktId     uint         `gorm:"pkt_id" json:"pkt_id"` // Link to Paket
	Nomor     string       `gorm:"nomor"`
	Jenis     string       `form:"jenis"`
	Hari      string       `form:"hari"`      // NEW
	Tanggal   sql.NullTime `form:"tanggal"`
	Tempat    string       `form:"tempat"`    // NEW
	Waktu     string       `form:"waktu"`     // NEW
	SubKeg    string       `form:"sub_keg"`   // NEW (Sub Kegiatan)
	Pengadaan string       `form:"pengadaan"` // NEW (Nama Pengadaan)
	Uraian    string       `form:"uraian"`
	DokId     uint         `form:"dok_id"`
}

func (BeritaAcara) TableName() string {
	return "berita_acara"
}

func GetBeritaAcara(id uint) BeritaAcara {
	var result BeritaAcara
	db.First(&result, id)
	return result
}

type Reviu struct {
	gorm.Model
	Bidang 		string		`json:"bidang"`
	Content		string 		`json:"content"`
	Opsi1		string 		`json:"opsi1"`
	Opsi2		string 		`json:"opsi2"`
}

func (Reviu) TableName() string {
	return "reviu"
}

func GetReviu(id uint) Reviu {
	var reviu Reviu
	db.First(&reviu, id)
	return reviu
}

func SaveReviu(reviu *Reviu) error {
	return db.Save(reviu).Error
}

func DeleteReviu(reviu *Reviu) error {
	return db.Delete(reviu).Error
}

func GetAllReviu() []Reviu {
	var results []Reviu
	db.Find(&results)
	return results
}

type ReviuPaket struct {
	gorm.Model
	PktId			uint 		`json:"pkt_id"`
	RevId			uint		`json:"rev_id"`
	Status          int         `json:"status"` // 0: Kosong, 1: Sesuai/Tersedia, 2: Tidak Sesuai/Tidak Tersedia
	Keterangan		string		`json:"Keteranga"`
	CatatanKhusus	string		`json:"catatan_khusus"`
	PegId			uint 		`json:"peg_id"`
}

func (ReviuPaket) TableName() string {
	return "reviu_paket"
}
