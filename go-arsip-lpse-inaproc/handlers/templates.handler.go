package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func EditTemplates(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/templates"
	id, _ := c.ParamsInt("id")
	if id != 0 {
		template := services.GetTemplates(uint(id))
		if template.ID == 0 {
			return c.SendStatus(404)
		}
		mp["template"] = template
		mp["url"] = "/templates/edit/" + utils.IntToString(id)
	}
	return c.Render("templates/form-templates", mp)
}

func CreateTemplates(c *fiber.Ctx) error {
	template := new(models.Templates)
	err := c.BodyParser(template)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Templates Gagal","/templates/edit")
	}
	template.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	template.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	err = services.CreateTemplates(*template)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Templates Gagal","/templates/edit")
	}
	return flashSuccess(c, "Tambah Templates Sukses","/templates")
}

func GetTemplates(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	templates := services.GetTemplates(uint(id))
	if templates.ID == 0 {
		return c.SendStatus(404)
	}
	mp["templates"] = templates
	return c.Render("templates/templates-detil", mp)
}

func GetAllTemplates(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("templates/templates", mp)
}

func UpdateTemplates(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	template := services.GetTemplates(uint(id))
	if template.ID == 0 {
		return flashError(c, "Template tidak ditemukan","/templates/edit/" + utils.IntToString(id))
	}
	err := c.BodyParser(&template)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Template Gagal", "/templates/edit/" + utils.IntToString(id))
	}
	err = services.SaveTemplates(template)
	return flashSuccess(c, "Edit Template Sukses", "/templates")
}

func DeleteTemplates(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	template := services.GetTemplates(uint(id))
	if template.ID == 0 {
		return flashError(c, "Templates tidak ditemukan","/templates")
	}
	err := services.DeleteTemplates(template)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Templates Gagal", "/templates")
	}
	return flashSuccess(c, "Hapus Templates Sukses","/templates")
}


func GetJsonTemplates(c *fiber.Ctx) error {
	return services.GetDataTableTemplates(c)
}

func GetTemplateSK(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := uint(103)
	pegawai := services.GetPegawai(id);
	mp["pegawai"] = pegawai
	return c.Render("templates/template-sk",mp)
}

func GetTemplateSKPokja(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("templates/template-sk-pokja",mp)
}
