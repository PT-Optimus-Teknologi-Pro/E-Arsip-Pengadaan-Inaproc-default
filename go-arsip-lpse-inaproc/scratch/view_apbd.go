package main

import (
	"arsip/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(config.GetDbUrl()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	type Result struct {
		NamaSatker      string
		TotalBelanja    float64
		BelanjaPengadaan float64
		PaguPenyedia    float64
	}

	var res []Result
	// Mencari data untuk Dinas Kesehatan atau SKPD lainnya yang ada di gambar
	query := `
		SELECT s.nama as nama_satker, a.total_belanja, a.belanja_pengadaan, p.pagu as pagu_penyedia
		FROM struktur_anggaran a
		LEFT JOIN satker s ON a.id_satker = s.id
		LEFT JOIN (
			SELECT id_satker, sum(pagu) as pagu 
			FROM paket_sirup 
			WHERE tahun = 2026 AND paket_aktif='TRUE' AND paket_terumumkan='TRUE' 
			GROUP BY id_satker
		) p ON a.id_satker = p.id_satker
		WHERE a.tahun_anggaran = 2026
		ORDER BY a.total_belanja DESC
		LIMIT 10
	`
	db.Raw(query).Scan(&res)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-35s | %-15s | %-15s\n", "Nama Satker", "Total APBD", "Pagu Terumumkan")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, r := range res {
		fmt.Printf("%-35s | Rp %13.0f | Rp %13.0f\n", r.NamaSatker, r.TotalBelanja, r.PaguPenyedia)
	}
	fmt.Println("--------------------------------------------------------------------------------")
}
