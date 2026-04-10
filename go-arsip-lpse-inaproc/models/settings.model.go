package models

import "fmt"

// AppSettings menyimpan pengaturan aplikasi global (singleton — hanya 1 row)
type AppSettings struct {
	ID                uint   `gorm:"primaryKey;autoIncrement"`
	LogoPath          string `json:"logo_path" gorm:"default:''"`
	FaviconPath       string `json:"favicon_path" gorm:"default:''"`
	AppTitle          string `json:"app_title" gorm:"default:'e-Arsip Pengadaan'"`
	HeroBadge         string `json:"hero_badge" gorm:"default:'Portal Transparansi Pengadaan Daerah'"`
	HeroTitle         string `json:"hero_title" gorm:"default:'Sistem Informasi Arsip dan Monitoring Pengadaan'"`
	HeroSubtitle      string `json:"hero_subtitle" gorm:"default:'Aplikasi terintegrasi untuk pencatatan, pengarsipan, dan pelaporan progres pengadaan barang dan jasa pemerintah daerah berbasis data RUP.'"`
	FooterDescription string `json:"footer_description" gorm:"default:'Portal Informasi Monitoring Pengadaan Daerah. Menyediakan transparansi data rencana dan realisasi pengadaan barang/jasa secara real-time untuk efektivitas tata kelola pemerintahan.'"`
	FooterEmail       string `json:"footer_email" gorm:"default:'admin@lpse.example.go.id'"`
	FooterAddress     string `json:"footer_address" gorm:"default:'Gedung Sekretariat Daerah, Lantai 2. Bagian Pengadaan Barang dan Jasa.'"`
	FooterFacebook    string `json:"footer_facebook" gorm:"default:'#'"`
	FooterInstagram   string `json:"footer_instagram" gorm:"default:'#'"`
	FooterTwitter     string `json:"footer_twitter" gorm:"default:'#'"`
	LoadingSubtitle   string `json:"loading_subtitle" gorm:"default:'Portal Pengarsipan'"`
}

func (AppSettings) TableName() string {
	return "app_settings"
}

// HeroSlider untuk background hero yang bisa di-slide
type HeroSlider struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Image string `json:"image"`
}

func (HeroSlider) TableName() string {
	return "hero_sliders"
}

// FooterSocialLink — Social media links yang bisa dikelola lewat admin
type FooterSocialLink struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Icon  string `json:"icon"`  // feather icon name atau path gambar
	Label string `json:"label"` // Nama platform
	URL   string `json:"url"`
	Sort  int    `json:"sort" gorm:"default:0"`
}

func (FooterSocialLink) TableName() string { return "footer_social_links" }

// FooterQuickLink — Tautan Cepat di footer
type FooterQuickLink struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Sort  int    `json:"sort" gorm:"default:0"`
}

func (FooterQuickLink) TableName() string { return "footer_quick_links" }

// FooterService — Layanan Kami di footer
type FooterService struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Sort  int    `json:"sort" gorm:"default:0"`
}

func (FooterService) TableName() string { return "footer_services" }

// ─── Settings CRUD ───────────────────────────────────────────────────────────

func GetSettings() AppSettings {
	var s AppSettings
	var count int64
	db.Model(&AppSettings{}).Count(&count)
	db.First(&s)
	fmt.Printf("GetSettings: Count=%d, ID=%d, Title=%s\n", count, s.ID, s.AppTitle)
	if s.ID == 0 {
		s = AppSettings{ID: 1}
		db.Create(&s)
	}
	return s
}

func SaveSettings(s *AppSettings) error {
	return db.Save(s).Error
}

// ─── HeroSlider CRUD ─────────────────────────────────────────────────────────

func GetHeroSliders() []HeroSlider {
	var res []HeroSlider
	db.Find(&res)
	return res
}

func CreateHeroSlider(s *HeroSlider) error {
	return db.Create(s).Error
}

func DeleteHeroSlider(id uint) error {
	return db.Delete(&HeroSlider{}, id).Error
}

// ─── FooterSocialLink CRUD ───────────────────────────────────────────────────

func GetFooterSocialLinks() []FooterSocialLink {
	var res []FooterSocialLink
	db.Order("sort asc, id asc").Find(&res)
	return res
}

func CreateFooterSocialLink(v *FooterSocialLink) error { return db.Create(v).Error }

func UpdateFooterSocialLink(id uint, v *FooterSocialLink) error {
	return db.Model(&FooterSocialLink{}).Where("id = ?", id).Updates(v).Error
}

func DeleteFooterSocialLink(id uint) error {
	return db.Delete(&FooterSocialLink{}, id).Error
}

// ─── FooterQuickLink CRUD ────────────────────────────────────────────────────

func GetFooterQuickLinks() []FooterQuickLink {
	var res []FooterQuickLink
	db.Order("sort asc, id asc").Find(&res)
	return res
}

func CreateFooterQuickLink(v *FooterQuickLink) error { return db.Create(v).Error }

func UpdateFooterQuickLink(id uint, v *FooterQuickLink) error {
	return db.Model(&FooterQuickLink{}).Where("id = ?", id).Updates(v).Error
}

func DeleteFooterQuickLink(id uint) error {
	return db.Delete(&FooterQuickLink{}, id).Error
}

// ─── FooterService CRUD ──────────────────────────────────────────────────────

func GetFooterServices() []FooterService {
	var res []FooterService
	db.Order("sort asc, id asc").Find(&res)
	return res
}

func CreateFooterService(v *FooterService) error { return db.Create(v).Error }

func UpdateFooterService(id uint, v *FooterService) error {
	return db.Model(&FooterService{}).Where("id = ?", id).Updates(v).Error
}

func DeleteFooterService(id uint) error {
	return db.Delete(&FooterService{}, id).Error
}
