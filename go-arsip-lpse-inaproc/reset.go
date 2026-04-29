package main

import (
	"arsip/config"
	"arsip/models"
	"arsip/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ResetApp() {
	db, err := gorm.Open(postgres.Open(config.GetDbUrl()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var users []models.Pegawai
	db.Find(&users)

	fmt.Println("Daftar User:")
	for _, u := range users {
		fmt.Printf("User: %s, Role: %s, ID: %d\n", u.PegNamauser, u.Usrgroup, u.ID)
	}

	// Update password admin to 123456
	var admin models.Pegawai
	db.First(&admin, "peg_namauser = ?", "ADMIN")
	if admin.ID != 0 {
		admin.Passw = utils.HashPassword("123456")
		db.Save(&admin)
		fmt.Println("Password ADMIN telah direset menjadi: 123456")
	} else {
		fmt.Println("User ADMIN tidak ditemukan")
	}
}
