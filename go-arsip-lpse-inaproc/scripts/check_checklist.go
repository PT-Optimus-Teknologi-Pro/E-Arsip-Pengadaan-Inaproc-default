package main

import (
    "fmt"
    "log"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Checklist struct {
    gorm.Model
    Jenis       int `gorm:"not null" json:"jenis"`
    Metode      int `gorm:"not null,default:0" json:"metode"`
}

func (Checklist) TableName() string {
    return "checklist"
}

func mainCheckChecklist() {
    db, err := gorm.Open(sqlite.Open("arsip.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database")
    }

    var checks []Checklist
    db.Find(&checks)
    fmt.Printf("Total checklists: %d\n", len(checks))
    for _, c := range checks {
        fmt.Printf("ID: %d, Jenis: %d, Metode: %d\n", c.ID, c.Jenis, c.Metode)
    }
}
