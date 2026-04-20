package models

import (
	"arsip/config"
	"arsip/utils"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	var err error
	db, err = gorm.Open(postgres.Open(config.GetDbUrl()), &gorm.Config{
		CreateBatchSize: 1000, // Sets the default batch size
	})
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to Connect to database...")
	}
	db.AutoMigrate(&Inbox{}, &HakAkses{}, &Checklist{}, &ChecklistDok{}, &DokTemplate{}, &Provinsi{}, &Kabupaten{}, &BukuTamu{}, &Feedback{}, &FeedbackKategori{}, &Templates{}, &Document{})
	db.AutoMigrate(&Agency{}, &Ukpbj{}, &Pegawai{}, &SatkerSirup{}, &StrukturAnggaran{}, &PaketSirup{}, &RupSwakelola{}, &SwakelolaSirup{}, &Panitia{}, &AnggotaPanitia{},
		&PejabatPengadaan{}, &PejabatPengadaanSatker{}, &PejabatPengadaanPegawai{},
		&Tender{}, &Nontender{}, &Pencatatan{}, &Swakelola{})
	db.AutoMigrate(&Anggaran{}, &Paket{}, &PaketAnggaran{}, &PaketSatker{}, &PaketLokasi{}, &ChecklistPaket{}, &ChecklistPaketHistory{}, &DokPaket{},
		&BeritaAcara{}, &Reviu{}, &ReviuPaket{}, &KajiUlang{}, &DokPersiapan{}, &PerubahanData{}, &PersetujuanDokPersiapan{}, &PaketPPk{})
	db.AutoMigrate(&Itkp{}, &AppSettings{}, &HeroSlider{}, &FooterSocialLink{}, &FooterQuickLink{}, &FooterService{}, &DokumenTercetak{})

	fmt.Println("Connected to database...")
}

type Date time.Time

func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
    	*ct = Date(time.Time{})
     	return
  	}
	nt, err := time.Parse("2006-01-02", s)
	*ct = Date(nt)
	return
}

func (date *Date) Scan(value interface{}) (err error) {
	if value == nil {
		*date = Date(time.Time{})
		return nil
	}
	*date = Date(value.(time.Time))
	return nil
}

func (date Date) Value() (driver.Value, error) {
	t := time.Time(date)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format("2006-01-02"), nil
}

func (date Date) Format() string {
	t := time.Time(date)
	if t.IsZero() {
		return ""
	}
	return t.Format("02-01-2006")
}

type Datetime time.Time

func (ct *Datetime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
    	*ct = Datetime(time.Time{})
     	return
  	}
	s = strings.ReplaceAll(s, "Z", "")
	s = strings.ReplaceAll(s, "T", " ")
	nt, err := time.Parse("2006-01-02 15:04:05", s)
	*ct = Datetime(nt)
	return
}

func (date *Datetime) Scan(value interface{}) (err error) {
	if value == nil {
		*date = Datetime(time.Time{})
		return nil
	}
	*date = Datetime(value.(time.Time))
	return nil
}

func (date Datetime) Value() (driver.Value, error) {
	t := time.Time(date)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format("2006-01-02 15:04:05"), nil
}

func (date Datetime) Format() string {
	t := time.Time(date)
	if t.IsZero() {
		return ""
	}
	return t.Format("02-01-2006 15:04")
}

type StringUint uint

func (st *StringUint) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringUint(v)
	case float64:
		*st = StringUint(int(v))
	case string:
		*st = StringUint(utils.StringToUint(v))
	}
	return nil
}

func Count(query string, values ...interface{}) int64 {
	var total int64
	db.Raw(query, values...).Count(&total)
	return total
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
