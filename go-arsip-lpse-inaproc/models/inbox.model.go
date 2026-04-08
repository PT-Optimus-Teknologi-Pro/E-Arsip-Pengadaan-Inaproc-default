package models

import (
	"gorm.io/gorm"
	"time"
)

type Inbox struct {
	gorm.Model
	EnqueueDate time.Time
	Subject     string
	Content     string
	Status      string
	PegId       uint
}

func (Inbox) TableName() string {
	return "inbox"
}

func GetInbox(id uint) Inbox {
	var inbox Inbox
	db.First(&inbox, id)
	return inbox
}

func SaveInbox(inbox *Inbox) error  {
	return db.Save(inbox).Error
}

func DeleteInbox(inbox *Inbox) error  {
	return db.Delete(inbox).Error
}

type HakAkses struct {
	gorm.Model
	PegId        uint
	PeriodeAwal  time.Time
	PeriodeAkhir time.Time
	Usrgroup     string // PPK , PP, Pokja, Kepala UKPBJ
}
