package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DokumenTercetak struct {
	ID                    uint      `gorm:"primarykey"`
	PaketID               uint      `gorm:"index"`
	JenisDokumen          string    // e.g. "BA_KAJIULANG", "SK_PENUNJUKAN_POKJA"
	NomorSurat            string
	TentangSurat          string
	TahunSurat            string
	TempatPenetapan       string
	TanggalPenetapan      time.Time `gorm:"type:date"`
	NomorKeputusanSekda   string
	TanggalTerbitKeputusan time.Time `gorm:"type:date"`
	PembuatPegawaiID      uint      // ID of the user who printed this
	Md5Hash               string    `gorm:"index;unique"`
	UrlValidasi           string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// TableName overrides the table name
func (DokumenTercetak) TableName() string {
	return "dokumen_tercetak"
}

// GenerateMD5 creates the unique hash before creating the record
func (d *DokumenTercetak) BeforeCreate(tx *gorm.DB) (err error) {
	if d.Md5Hash == "" {
		// Unique string: PaketID + JenisDokumen + Timestamp
		raw := fmt.Sprintf("%d-%s-%d", d.PaketID, d.JenisDokumen, time.Now().UnixNano())
		hash := md5.Sum([]byte(raw))
		d.Md5Hash = hex.EncodeToString(hash[:])
	}
	return
}
