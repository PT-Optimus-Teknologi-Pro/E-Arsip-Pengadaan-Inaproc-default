package models

import (
	"gorm.io/gorm"
)

// HeroSlider model for the premium homepage hero section
type HeroSlider struct {
	gorm.Model
	Badge       string `json:"badge" form:"badge"`
	Title       string `json:"title" form:"title"`
	Subtitle    string `json:"subtitle" form:"subtitle"`
	Image       string `json:"image" form:"image"` // Path or URL to the background image
	Button1Text string `json:"button1_text" form:"button1_text"`
	Button1Link string `json:"button1_link" form:"button1_link"`
	Button2Text string `json:"button2_text" form:"button2_text"`
	Button2Link string `json:"button2_link" form:"button2_link"`
	Order       int    `json:"order" form:"order"` // For manual sorting
	IsActive    bool   `json:"is_active" form:"is_active" gorm:"default:true"`
}

func (HeroSlider) TableName() string {
	return "hero_sliders"
}

// SocialLink model for Dynamic Social Media icons in Footer
type SocialLink struct {
	gorm.Model
	Icon  string `json:"icon" form:"icon"` // Feather icon name or image path
	Label string `json:"label" form:"label"`
	URL   string `json:"url" form:"url"`
	Sort  int    `json:"sort" form:"sort"`
}

func (SocialLink) TableName() string {
	return "social_links"
}

// QuickLink model for Footer Column 2
type QuickLink struct {
	gorm.Model
	Label string `json:"label" form:"label"`
	URL   string `json:"url" form:"url"`
	Sort  int    `json:"sort" form:"sort"`
}

func (QuickLink) TableName() string {
	return "quick_links"
}

// ServiceLink model for Footer Column 3
type ServiceLink struct {
	gorm.Model
	Label string `json:"label" form:"label"`
	URL   string `json:"url" form:"url"`
	Sort  int    `json:"sort" form:"sort"`
}

func (ServiceLink) TableName() string {
	return "service_links"
}

// SiteSetting model for general site-wide configuration (Logos, Contact, SEO)
type SiteSetting struct {
	gorm.Model
	Key   string `json:"key" gorm:"uniqueIndex"` // e.g., "app_name", "logo_header", "footer_desc"
	Value string `json:"value"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}

// Helper funcs to fetch settings
func GetSetting(key string, fallback string) string {
	var setting SiteSetting
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		return fallback
	}
	return setting.Value
}

func SaveSetting(key string, value string) error {
	var setting SiteSetting
	err := db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		// Create new
		return db.Create(&SiteSetting{Key: key, Value: value}).Error
	}
	// Update existing
	setting.Value = value
	return db.Save(&setting).Error
}
