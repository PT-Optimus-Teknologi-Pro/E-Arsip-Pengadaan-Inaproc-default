package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	tables := []string{"tender", "tender_selesai", "nontender", "nontender_selesai", "katalog", "paket"}
	for _, table := range tables {
		var count int64
		err := db.Table(table).Count(&count).Error
		if err != nil {
			fmt.Printf("Error counting %s: %v\n", table, err)
			continue
		}
		fmt.Printf("Table %s: %d records\n", table, count)

		// Get columns and types
		rows, err := db.Table(table).Limit(1).Rows()
		if err != nil {
			fmt.Printf("Error getting columns for %s: %v\n", table, err)
			continue
		}
		cols, _ := rows.Columns()
		types, _ := rows.ColumnTypes()
		fmt.Printf("Columns for %s:\n", table)
		for i, col := range cols {
			fmt.Printf("  %s: %s\n", col, types[i].DatabaseTypeName())
		}
		rows.Close()
	}

	var tenderMatch int64
	db.Table("tender").Joins("JOIN tender_selesai ts ON tender.kd_tender = ts.kd_tender").Count(&tenderMatch)
	fmt.Printf("Tender JOIN Tender Selesai: %d matches\n", tenderMatch)

	var nontenderMatch int64
	db.Table("nontender").Joins("JOIN nontender_selesai nts ON nontender.kd_nontender = nts.kd_nontender").Count(&nontenderMatch)
	fmt.Printf("Nontender JOIN Nontender Selesai: %d matches\n", nontenderMatch)

	var tenderMethods []string
	db.Table("tender").Joins("JOIN tender_selesai ts ON tender.kd_tender = ts.kd_tender").Where("tender.mtd_pemilihan IS NOT NULL AND tender.mtd_pemilihan != ''").Distinct("tender.mtd_pemilihan").Pluck("tender.mtd_pemilihan", &tenderMethods)
	fmt.Printf("Tender Methods: %v\n", tenderMethods)
}
