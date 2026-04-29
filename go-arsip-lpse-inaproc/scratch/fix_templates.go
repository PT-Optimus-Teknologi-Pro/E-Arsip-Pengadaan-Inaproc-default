package main

import (
	"arsip/models"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func mainFixTemplates() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := models.GetDB()
	var pps []models.PejabatPengadaan
	db.Find(&pps)

	fmt.Printf("Total PejabatPengadaan records: %d\n", len(pps))
	for _, p := range pps {
		fmt.Printf("ID: %d, NoSk: %s\n", p.ID, p.NoSk)
	}
}
