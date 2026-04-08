package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateDocTemplate(c *fiber.Ctx) error {
	doktemplate := new(models.DokTemplate)
	err := c.BodyParser(doktemplate)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Dokumen template Gagal","/doc-template/edit")
	}
	doktemplate.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	doktemplate.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	mp := currentMap(c)
	userid := mp["id"].(uint)
	err = services.SaveDocTemplate(c, *doktemplate, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Dokumen template Gagal","/doc-template/edit")
	}
	return flashSuccess(c, "Tambah/Edit  Dokumen template Sukses","/doc-template")
}

// Get All Users from db
func GetAllDocTemplate(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("dok-template/dok-template", mp)
}

// GetSingleUser from db
func GetDocTemplate(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	doktemplate := services.GetDokTemplate(id)
	if doktemplate.ID == 0 {
		return c.SendStatus(404)
	}
	mp := currentMap(c)
	mp["doktemplate"] = doktemplate
	mp["file"] = doktemplate.Dokumen()
	return c.Render("dok-template/dok-template-view", mp)
}

// update a user in db
func UpdateDocTemplate(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	template := services.GetDokTemplate(id)
	if template.ID == 0 {
		return flashError(c, "Dokumen template tidak ditemukan", "/doc-template/edit/" + utils.UintToString(id))
	}
	err := c.BodyParser(&template)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Dokumen Template Gagal", "/doc-template/edit/" + utils.UintToString(id))
	}
	template.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	template.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	mp := currentMap(c)
	userid := mp["id"].(uint)
	err = services.SaveDocTemplate(c, template, userid)
	// Return the updated user
	return flashSuccess(c, "Edit Dokumen Template Sukses","/doc-template")
}

// delete user in db by ID
func DeleteDocTemplate(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	template := services.GetDokTemplate(id)
	if template.ID == 0 {
		return flashError(c, "Dokumen template tidak ditemukan", "/doc-template")
	}
	err := services.DeleteDocTemplate(template)
	if err != nil {
		return flashError(c, "Hapus Dokumen template Gagal", "/doc-template")
	}
	return flashSuccess(c, "Hapus Dokumen template Sukses", "/doc-template")
}

func GetJsonDocTemplate(c *fiber.Ctx) error {
	return services.GetDataTableDocTemplate(c)
}

func EditDocTemplate(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/doc-template"
	id := utils.StringToUint(c.Params("id"))
	if id != 0 {
		template := services.GetDokTemplate(id)
		if template.ID == 0 {
			return c.SendStatus(404)
		}
		mp["doktemplate"] = template
		mp["file"] = template.Dokumen()
		mp["url"] = "/doc-template/" + utils.UintToString(id)
	}
	return c.Render("dok-template/form-dok-template", mp)
}
