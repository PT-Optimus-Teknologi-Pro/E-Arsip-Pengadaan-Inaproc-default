package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TenderArsiparis struct {
	KdTender       uint    `json:"kd_tender" gorm:"column:kd_tender"`
	KdRup          string  `json:"kd_rup" gorm:"column:kd_rup"`
	NamaPaket      string  `json:"nama_paket" gorm:"column:nama_paket"`
	MtdPemilihan   string  `json:"mtd_pemilihan" gorm:"column:mtd_pemilihan"`
	JenisPengadaan string  `json:"jenis_pengadaan" gorm:"column:jenis_pengadaan"`
	Pagu           float64 `json:"pagu" gorm:"column:pagu"`
	NilaiKontrak   float64 `json:"nilai_kontrak" gorm:"column:nilai_kontrak"`
}

type NontenderArsiparis struct {
	KdNontender    string  `json:"kd_nontender" gorm:"column:kd_nontender"`
	KdRup          string  `json:"kd_rup" gorm:"column:kd_rup"`
	NamaPaket      string  `json:"nama_paket" gorm:"column:nama_paket"`
	MtdPemilihan   string  `json:"mtd_pemilihan" gorm:"column:mtd_pemilihan"`
	JenisPengadaan string  `json:"jenis_pengadaan" gorm:"column:jenis_pengadaan"`
	Pagu           float64 `json:"pagu" gorm:"column:pagu"`
	NilaiKontrak   float64 `json:"nilai_kontrak" gorm:"column:nilai_kontrak"`
}


func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Test Tender
	orm := db.Table("tender").
		Select("tender.kd_tender, tender.kd_rup, tender.nama_paket, tender.mtd_pemilihan, tender.jenis_pengadaan, tender.pagu, ts.nilai_kontrak").
		Joins("JOIN tender_selesai ts ON tender.kd_tender = ts.kd_tender")

	var total int64
	err = orm.Session(&gorm.Session{}).Count(&total).Error
	fmt.Println("Tender Total:", total, "Err:", err)

	var datas []TenderArsiparis
	err = orm.Limit(10).Offset(0).Find(&datas).Error
	fmt.Println("Tender Data Count:", len(datas), "Err:", err)

	// Test NonTender
	sub1 := db.Table("nontender").
		Select("nontender.kd_nontender::text as kd_nontender, nontender.kd_rup, nontender.nama_paket, nontender.mtd_pemilihan, nontender.jenis_pengadaan, nontender.pagu, nts.nilai_kontrak").
		Joins("JOIN nontender_selesai nts ON nontender.kd_nontender = nts.kd_nontender")

	sub2 := db.Table("katalog").
		Select("katalog.order_id as kd_nontender, katalog.rup_code as kd_rup, katalog.rup_name as nama_paket, 'E-Purchasing' as mtd_pemilihan, 'Barang/Jasa' as jenis_pengadaan, katalog.total as pagu, katalog.total as nilai_kontrak")

	orm2 := db.Table("(?) AS comb", db.Raw("? UNION ALL ?", sub1, sub2))
	
	var total2 int64
	err = orm2.Session(&gorm.Session{}).Count(&total2).Error
	fmt.Println("Nontender Total:", total2, "Err:", err)

	var datas2 []NontenderArsiparis
	err = orm2.Limit(10).Offset(0).Find(&datas2).Error
	fmt.Println("Nontender Data Count:", len(datas2), "Err:", err)
}
