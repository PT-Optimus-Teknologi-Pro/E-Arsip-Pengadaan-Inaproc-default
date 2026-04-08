package models

import (
	"database/sql/driver"
	"log/slog"
	"strings"
	"sync-inaproc/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const MAX_DATA 	= 1000

var DB *gorm.DB

func init() {
  var err error
  DB, err = gorm.Open(postgres.Open(config.GetDbUrl()), &gorm.Config{
  	CreateBatchSize: 1000,
  })
  if err != nil {
    panic("failed to connect database")
    } else {
   		slog.Info("connected to database")
    }
    err = DB.AutoMigrate(&Satker{}, &Program{}, &RupAnggaran{}, &Rup{}, &RupSwakelola{}, &RupSwakelolaAnggaran{},
    			&Katalog{}, &Penyedia{}, &KatalogArchive{}, &PenyediaArchive{},
       			&Nontender{}, &JadwalNontender{}, &NontenderSelesai{},
    			&Pencatatan{}, &PencatatanRealisasi{}, &Swakelola{}, &SwakelolaRealisasi{},
       			&Tender{}, &Jadwal{} , &Peserta{}, &TenderSelesai{},
    			&Kontrak{}, &KontrakNontender{})
    if err != nil {
        slog.Error("Error sync table ", "error", err)
    }
}


type Date time.Time

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(date).Format("2006-01-02")), nil
}

func (c *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*c = Date(time.Time{}) // Set to zero time
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*c = Date(t)
	return nil
}

type CustomeDate time.Time

func (date CustomeDate) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}
func (date CustomeDate) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(date).Format("2006-01-02 15:04:05.000000Z07:00")), nil
}

func (c *CustomeDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*c = CustomeDate(time.Time{}) // Set to zero time
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05.000000Z07:00", s)
	if err != nil {
		return err
	}
	*c = CustomeDate(t)
	return nil
}
