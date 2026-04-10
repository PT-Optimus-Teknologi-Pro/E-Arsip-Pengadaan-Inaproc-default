package services

import "arsip/models"

// GetSettings mengambil pengaturan aplikasi
func GetSettings() models.AppSettings {
	return models.GetSettings()
}

// SaveSettings menyimpan pengaturan aplikasi
func SaveSettings(s *models.AppSettings) error {
	return models.SaveSettings(s)
}

// ─── HeroSlider ────────────────────────────────────────────────────────────────

func GetHeroSliders() []models.HeroSlider     { return models.GetHeroSliders() }
func SaveHeroSlider(s *models.HeroSlider) error { return models.CreateHeroSlider(s) }
func DeleteHeroSlider(id uint) error          { return models.DeleteHeroSlider(id) }

// ─── FooterSocialLink ─────────────────────────────────────────────────────────

func GetFooterSocialLinks() []models.FooterSocialLink { return models.GetFooterSocialLinks() }
func CreateFooterSocialLink(v *models.FooterSocialLink) error {
	return models.CreateFooterSocialLink(v)
}
func UpdateFooterSocialLink(id uint, v *models.FooterSocialLink) error {
	return models.UpdateFooterSocialLink(id, v)
}
func DeleteFooterSocialLink(id uint) error { return models.DeleteFooterSocialLink(id) }

// ─── FooterQuickLink ──────────────────────────────────────────────────────────

func GetFooterQuickLinks() []models.FooterQuickLink { return models.GetFooterQuickLinks() }
func CreateFooterQuickLink(v *models.FooterQuickLink) error {
	return models.CreateFooterQuickLink(v)
}
func UpdateFooterQuickLink(id uint, v *models.FooterQuickLink) error {
	return models.UpdateFooterQuickLink(id, v)
}
func DeleteFooterQuickLink(id uint) error { return models.DeleteFooterQuickLink(id) }

// ─── FooterService ────────────────────────────────────────────────────────────

func GetFooterServices() []models.FooterService         { return models.GetFooterServices() }
func CreateFooterService(v *models.FooterService) error { return models.CreateFooterService(v) }
func UpdateFooterService(id uint, v *models.FooterService) error {
	return models.UpdateFooterService(id, v)
}
func DeleteFooterService(id uint) error { return models.DeleteFooterService(id) }
