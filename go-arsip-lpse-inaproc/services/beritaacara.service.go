package services

import (
	"fmt"
	"time"
	"arsip/models"
)

func GetAllReviu() []models.Reviu {
	return models.GetAllReviu()
}


func SimpanBaNego(id uint, nomor string, tanggal time.Time) error {
	paket := GetPaket(id)
	if paket.Status < 4 {
		return fmt.Errorf("Berita Acara Gagal tersimpan dikarenakan status paket belum proses pengadaan")
	}

	return nil
}

func SimpanBaPenetapan(id uint, nomor string, tanggal time.Time) error {
	paket := GetPaket(id)
	if paket.Status < 4 {
		return fmt.Errorf("Berita Acara Gagal tersimpan dikarenakan status paket belum proses pengadaan")
	}

	return nil
}
