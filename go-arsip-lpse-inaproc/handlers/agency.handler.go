package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditAgency(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/agency"
	id, _ := c.ParamsInt("id")
	if id != 0 {
		agency := services.GetAgency(uint(id))
		if agency.ID == 0 {
			return c.SendStatus(404)
		}
		mp["agency"] = agency
		mp["url"] = "/agency/" + utils.IntToString(id)
	}
	return c.Render("agency/form-agency", mp)
}

func CreateAgency(c *fiber.Ctx) error {
	agency := new(models.Agency)
	err := c.BodyParser(agency)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Agency Gagal", "/agency/edit")
	}
	agency.AgcTglDaftar, _ = time.Parse("2006-01-02", c.FormValue("tglDaftar"))
	err = services.CreateAgency(*agency)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Agency Gagal", "/agency/edit")
	}
	// Return the created agency
	return flashSuccess(c, "Tambah Agency Sukses", "/agency")
}

// Get All Users from db
func GetAllAgency(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("agency/agency", mp)
}

// GetSingleUser from db
func GetAgency(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	agency := services.GetAgency(uint(id))
	mp := currentMap(c)
	if agency.ID == 0 {
		return c.SendStatus(404)
	}
	mp["agency"] = agency
	return c.Render("agency/agency-detil", mp)
}

func UpdateAgency(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	agency := services.GetAgency(uint(id))
	if agency.ID == 0 {
		return flashError(c, "Agency tidak ditemukan", "/agency/edit/"+utils.IntToString(id))
	}
	err := c.BodyParser(&agency)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Agency Gagal", "/agency/edit/"+utils.IntToString(id))
	}
	err = services.SaveAgency(agency);
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Agency Gagal", "/agency/edit/"+utils.IntToString(id))
	}
	return flashSuccess(c, "Edit Agency Sukses","/agency")
}

// delete user in db by ID
func DeleteAgency(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	agency := services.GetAgency(uint(id))
	if agency.ID == 0 {
		return flashError(c, "Agency tidak ditemukan","/agency")
	}
	err := services.DeleteAgency(agency)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Agency Gagal","/agency")
	}
	return flashSuccess(c, "Hapus Agency Sukses","/agency")
}

func GetJsonAgency(c *fiber.Ctx) error {
	return services.GetDataTableAgency(c)
}
