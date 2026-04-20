package models

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
	FooterPhone       string `json:"footer_phone" gorm:"default:''"`
	FooterWorkHours   string `json:"footer_work_hours" gorm:"default:'Senin - Jumat: 08.00 - 16.00 WITA'"`
	FooterFacebook    string `json:"footer_facebook" gorm:"default:'#'"`
	FooterInstagram   string `json:"footer_instagram" gorm:"default:'#'"`
	FooterTwitter     string `json:"footer_twitter" gorm:"default:'#'"`
	LoadingSubtitle   string `json:"loading_subtitle" gorm:"default:'Portal Pengarsipan'"`
	// Identitas Instansi untuk Dokumen/Kop Surat
	DocInstansi      string `json:"doc_instansi" gorm:"default:'PEMERINTAH KOTA BANJARMASIN'"`
	DocSubInstansi   string `json:"doc_sub_instansi" gorm:"default:'S E K R E T A R I A T D A E R A H'"`
	DocAddress       string `json:"doc_address" gorm:"default:'J . RE. Martadinata No.1 Banjarmasin 70111'"`
	DocPhone         string `json:"doc_phone" gorm:"default:'(0511) 4368142 4368145'"`
	DocFax           string `json:"doc_fax" gorm:"default:'3353933'"`
	DocWebsite       string `json:"doc_website" gorm:"default:'http://www.banjarmasinkota.go.id/'"`
	DocEmail         string `json:"doc_email" gorm:"default:'admin@banjarmasinkota.go.id'"`
	DocPejabatNama   string `json:"doc_pejabat_nama" gorm:"default:'AHSAN BUDIMAN'"`
	DocPejabatJabata string `json:"doc_pejabat_jabata" gorm:"default:'SEKRETARIS DAERAH'"`
	DocPejabatNip    string `json:"doc_pejabat_nip" gorm:"default:''"`
	DocLogoPath      string `json:"doc_logo_path" gorm:"default:''"`
	DocRegion        string `json:"doc_region" gorm:"default:'KOTA BANJARMASIN'"`
	DocSignatureMode string `json:"doc_signature_mode" gorm:"default:'barcode'"`
	DocSignaturePath string `json:"doc_signature_path" gorm:"default:''"`
	FooterLogoWhite  bool   `json:"footer_logo_white" gorm:"default:false"`
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
	db.First(&s)
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
