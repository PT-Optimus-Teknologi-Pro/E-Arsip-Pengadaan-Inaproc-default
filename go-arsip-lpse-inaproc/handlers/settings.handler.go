package handlers

import (
	"arsip/cache"
	"arsip/models"
	"arsip/services"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const settingsUploadDir = "./public/uploads/settings"

func GetLogoSettings(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	settings := services.GetSettings()
	mp["settings"] = settings
	return c.Render("settings/logo-settings", mp)
}

func UpdateLogoSettings(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}

	if err := os.MkdirAll(settingsUploadDir, 0755); err != nil {
		log.Error("mkdir settings upload dir: ", err)
		return flashError(c, "Gagal membuat direktori upload", "/settings/logo")
	}

	settings := services.GetSettings()

	// Handle upload logo navbar
	logoFile, err := c.FormFile("logo")
	if err == nil && logoFile != nil {
		ext := strings.ToLower(filepath.Ext(logoFile.Filename))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".svg" || ext == ".webp" {
			filename := fmt.Sprintf("logo_%d%s", time.Now().Unix(), ext)
			dst := filepath.Join(settingsUploadDir, filename)
			if err := c.SaveFile(logoFile, dst); err != nil {
				log.Error("Gagal simpan logo: ", err)
			} else {
				settings.LogoPath = "/uploads/settings/" + filename
			}
		}
	}

	// Handle upload favicon
	faviconFile, err := c.FormFile("favicon")
	if err == nil && faviconFile != nil {
		ext := strings.ToLower(filepath.Ext(faviconFile.Filename))
		if ext == ".ico" || ext == ".png" || ext == ".svg" {
			filename := fmt.Sprintf("favicon_%d%s", time.Now().Unix(), ext)
			dst := filepath.Join(settingsUploadDir, filename)
			if err := c.SaveFile(faviconFile, dst); err != nil {
				log.Error("Gagal simpan favicon: ", err)
			} else {
				settings.FaviconPath = "/uploads/settings/" + filename
			}
		}
	}

	if appTitle := c.FormValue("app_title"); appTitle != "" {
		settings.AppTitle = appTitle
	}
	if loadingSubtitle := c.FormValue("loading_subtitle"); loadingSubtitle != "" {
		settings.LoadingSubtitle = loadingSubtitle
	}

	if err := services.SaveSettings(&settings); err != nil {
		log.Error("save settings: ", err)
		return flashError(c, "Gagal menyimpan pengaturan", "/settings/logo")
	}

	cache.Delete("app_settings")
	return flashSuccess(c, "Pengaturan berhasil disimpan", "/settings/logo")
}

func GetHeroSettings(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	settings := services.GetSettings()
	mp["settings"] = settings
	mp["sliders"] = services.GetHeroSliders()
	return c.Render("settings/hero-settings", mp)
}

func UpdateHeroSettings(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool) {
		return Forbiden(c)
	}

	settings := services.GetSettings()

	if heroBadge := c.FormValue("hero_badge"); heroBadge != "" {
		settings.HeroBadge = heroBadge
	}
	if heroTitle := c.FormValue("hero_title"); heroTitle != "" {
		settings.HeroTitle = heroTitle
	}
	if heroSubtitle := c.FormValue("hero_subtitle"); heroSubtitle != "" {
		settings.HeroSubtitle = heroSubtitle
	}

	if err := services.SaveSettings(&settings); err != nil {
		log.Error("save hero settings: ", err)
		return flashError(c, "Gagal menyimpan pengaturan hero", "/settings/hero")
	}

	cache.Delete("app_settings")
	return flashSuccess(c, "Pengaturan hero berhasil disimpan", "/settings/hero")
}

func UploadHeroSlider(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool) {
		return Forbiden(c)
	}

	if err := os.MkdirAll(settingsUploadDir, 0755); err != nil {
		return flashError(c, "Gagal membuat direktori upload", "/settings/hero")
	}

	file, err := c.FormFile("slider_image")
	if err != nil {
		return flashError(c, "File tidak ditemukan", "/settings/hero")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("slider_%d%s", time.Now().UnixNano(), ext)
	dst := fmt.Sprintf("%s/%s", settingsUploadDir, filename)

	if err := c.SaveFile(file, dst); err != nil {
		return flashError(c, "Gagal menyimpan file", "/settings/hero")
	}

	slider := models.HeroSlider{
		Image: "/uploads/settings/" + filename,
	}
	if err := services.SaveHeroSlider(&slider); err != nil {
		return flashError(c, "Gagal menyimpan data slider", "/settings/hero")
	}

	return flashSuccess(c, "Slider berhasil diunggah", "/settings/hero")
}

func DeleteHeroSlider(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}

	id, _ := c.ParamsInt("id")
	services.DeleteHeroSlider(uint(id))
	cache.Delete("hero_sliders")
	return flashSuccess(c, "Slider berhasil dihapus", "/settings/hero")
}

// ===== FOOTER SETTINGS =====

func GetFooterSettings(c *fiber.Ctx) error {
	sess := getSession(c)
	mp := currentMap(c)
	mp["settings"] = services.GetSettings()
	mp["name"] = sess.Get("name")
	mp["group"] = sess.Get("group")

	mp["socialLinks"] = services.GetFooterSocialLinks()
	mp["quickLinks"] = services.GetFooterQuickLinks()
	mp["servicesLinks"] = services.GetFooterServices()

	return c.Render("settings/footer-settings", mp)
}

func UpdateFooterSettings(c *fiber.Ctx) error {
	settings := services.GetSettings()
	settings.FooterDescription = c.FormValue("footer_description")
	settings.FooterEmail = c.FormValue("footer_email")
	settings.FooterAddress = c.FormValue("footer_address")
	settings.FooterFacebook = c.FormValue("footer_facebook")
	settings.FooterInstagram = c.FormValue("footer_instagram")
	settings.FooterTwitter = c.FormValue("footer_twitter")

	if err := services.SaveSettings(&settings); err != nil {
		return flashError(c, "Gagal menyimpan pengaturan footer", "/settings/footer")
	}

	cache.Delete("app_settings")
	return flashSuccess(c, "Pengaturan footer berhasil disimpan", "/settings/footer")
}

// ─── Footer Social Links Handlers ───────────────────────────────────────────

func CreateFooterSocialLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}

	icon := c.FormValue("icon")

	file, err := c.FormFile("icon_file")
	if err == nil && file != nil {
		if err := os.MkdirAll(settingsUploadDir, 0755); err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".svg" {
				filename := fmt.Sprintf("soc_%d%s", time.Now().UnixNano(), ext)
				dst := fmt.Sprintf("%s/%s", settingsUploadDir, filename)
				if err := c.SaveFile(file, dst); err == nil {
					icon = "/uploads/settings/" + filename
				}
			}
		}
	}

	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	err = services.CreateFooterSocialLink(&models.FooterSocialLink{
		Icon:  icon,
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	if err != nil {
		return flashError(c, "Gagal menambah tautan sosial", "/settings/footer")
	}
	cache.Delete("social_links")
	return flashSuccess(c, "Tautan sosial berhasil ditambahkan", "/settings/footer")
}

func UpdateFooterSocialLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")

	icon := c.FormValue("icon")

	file, err := c.FormFile("icon_file")
	if err == nil && file != nil {
		if err := os.MkdirAll(settingsUploadDir, 0755); err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".svg" {
				filename := fmt.Sprintf("soc_%d%s", time.Now().UnixNano(), ext)
				dst := fmt.Sprintf("%s/%s", settingsUploadDir, filename)
				if err := c.SaveFile(file, dst); err == nil {
					icon = "/uploads/settings/" + filename
				}
			}
		}
	}

	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	services.UpdateFooterSocialLink(uint(id), &models.FooterSocialLink{
		Icon:  icon,
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	cache.Delete("social_links")
	return flashSuccess(c, "Tautan sosial berhasil diperbarui", "/settings/footer")
}

func DeleteFooterSocialLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")
	services.DeleteFooterSocialLink(uint(id))
	cache.Delete("social_links")
	return flashSuccess(c, "Tautan sosial berhasil dihapus", "/settings/footer")
}

// ─── Footer Quick Links Handlers ────────────────────────────────────────────

func CreateFooterQuickLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	err := services.CreateFooterQuickLink(&models.FooterQuickLink{
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	if err != nil {
		return flashError(c, "Gagal menambah tautan cepat", "/settings/footer")
	}
	cache.Delete("quick_links")
	return flashSuccess(c, "Tautan cepat berhasil ditambahkan", "/settings/footer")
}

func UpdateFooterQuickLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")
	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	services.UpdateFooterQuickLink(uint(id), &models.FooterQuickLink{
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	cache.Delete("quick_links")
	return flashSuccess(c, "Tautan cepat berhasil diperbarui", "/settings/footer")
}

func DeleteFooterQuickLink(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")
	services.DeleteFooterQuickLink(uint(id))
	cache.Delete("quick_links")
	return flashSuccess(c, "Tautan cepat berhasil dihapus", "/settings/footer")
}

// ─── Footer Services Handlers ───────────────────────────────────────────────

func CreateFooterService(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	err := services.CreateFooterService(&models.FooterService{
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	if err != nil {
		return flashError(c, "Gagal menambah layanan", "/settings/footer")
	}
	cache.Delete("services_links")
	return flashSuccess(c, "Layanan berhasil ditambahkan", "/settings/footer")
}

func UpdateFooterService(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")
	sort, _ := c.ParamsInt("sort", 0)
	if formSort := c.FormValue("sort"); formSort != "" {
		fmt.Sscanf(formSort, "%d", &sort)
	}

	services.UpdateFooterService(uint(id), &models.FooterService{
		Label: c.FormValue("label"),
		URL:   c.FormValue("url"),
		Sort:  sort,
	})
	cache.Delete("services_links")
	return flashSuccess(c, "Layanan berhasil diperbarui", "/settings/footer")
}

func DeleteFooterService(c *fiber.Ctx) error {
	if !currentMap(c)["isAdmin"].(bool) && !currentMap(c)["isUkpbj"].(bool) {
		return Forbiden(c)
	}
	id, _ := c.ParamsInt("id")
	services.DeleteFooterService(uint(id))
	cache.Delete("services_links")
	return flashSuccess(c, "Layanan berhasil dihapus", "/settings/footer")
}
