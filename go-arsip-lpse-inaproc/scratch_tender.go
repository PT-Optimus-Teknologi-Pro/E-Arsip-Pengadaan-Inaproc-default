package main

import (
	"arsip/models"
	"fmt"
)

func main() {
	db := models.GetDB()
	var _paket models.Paket
	// Target Paket ID 8 (sesuai screenshot)
	if err := db.First(&_paket, 8).Error; err != nil {
		fmt.Printf("Paket tidak ditemukan: %v\n", err)
		return
	}

	fmt.Printf("Paket %d: %s\n", _paket.ID, _paket.Nama)
	fmt.Printf("Metode saat ini: %d\n", _paket.Metode)

	// Ubah ke Metode 12 (Tender) agar dianggap Paket Pokja
	_paket.Metode = 12
	if err := db.Save(&_paket).Error; err != nil {
		fmt.Printf("Gagal update metode: %v\n", err)
		return
	}

	fmt.Printf("Berhasil mengubah Metode Paket menjadi 12 (Tender)\n")
}
