package main

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=db-bungo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	type Satker struct {
		ID   uint
		Nama string
	}

	var satkers []Satker
	db.Raw("SELECT id, nama FROM satker WHERE nama ILIKE '%BADAN KEUANGAN%'").Scan(&satkers)

	for _, s := range satkers {
		fmt.Printf("ID: %d, Nama: %s\n", s.ID, s.Nama)
	}
}
