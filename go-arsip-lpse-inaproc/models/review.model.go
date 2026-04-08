package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Tahun      int
	Jenis      string
	Metode     string
	Nomor      string
	Bidang     string
	KajiUlang  string
	Kesimpulan string // pilihan : pakai – tidak pakai, sesuai – tidak sesuai, tersedia – tidak tersedia dll
}
