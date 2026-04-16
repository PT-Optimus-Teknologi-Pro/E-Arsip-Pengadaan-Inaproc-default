package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditPerubahanData(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/perubahan-data"
	id, _ := c.ParamsInt("id")
	if id > 0 {
		perubahan := services.GetPerubahanData(uint(id))
		if perubahan.ID == 0 {
			return c.SendStatus(404)
		}
		mp["perubahan"] = perubahan
		mp["url"] = "/perubahan-data/edit/" + utils.IntToString(id)
	}
	return c.Render("perubahan-data/form-perubahan-data", mp)
}

func CreatePerubahanData(c *fiber.Ctx) error {
	perubahan := new(models.PerubahanData)
	err := c.BodyParser(perubahan)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Perubahan Data Gagal","/perubahan-data/edit")
	}
	userssion := getUserSession(c);
	perubahan.PegId = userssion.Id
	err = services.CreatePerubahanData(c, perubahan, "file")
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Perubahan Data Gagal", "/perubahan-data/edit")
	}
	// Return the created agency
	return flashSuccess(c, "Tambah/Edit Perubahan Data Sukses", "/perubahan-data")
}

// Get All Users from db
func GetAllPerubahanData(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	user := models.GetPegawai(userid)
	mp["allow"] = user.IsApprove() && user.IsAktif()
	return c.Render("perubahan-data/perubahan-data", mp)
}

// GetSingleUser from db
func GetPerubahanData(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	perubahan := services.GetPerubahanData(uint(id))
	if perubahan.ID == 0 {
		return c.SendStatus(404)
	}
	mp["perubahan"] = perubahan
	return c.Render("perubahan-data/perubahan-data-detil", mp)
}

// update a user in db
func UpdatePerubahanData(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	perubahan := services.GetPerubahanData(uint(id))
	if perubahan.ID == 0 {
		return flashError(c, "Perubahan data tidak ditemukan","/perubahan-data/edit/" + utils.IntToString(id))
	}
	err := c.BodyParser(&perubahan)
	if err != nil {
		return flashError(c, "Edit Perubahan Gagal", "/perubahan-data/edit/" + utils.IntToString(id))
	}
	services.UpdatePerubahanData(perubahan)
	// Return the updated user
	return flashSuccess(c, "Edit Perubahan Sukses", "/perubahan-data")
}

// delete user in db by ID
func DeletePerubahanData(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	err := services.DeletePerubahanData(uint(id))
	if err != nil {
		return flashError(c, err.Error(), "/perubahan-data")
	}
	return flashSuccess(c, "Hapus Perubahan Sukses","/perubahan-data")
}

func GetJsonPerubahanData(c *fiber.Ctx) error {
	usrsession := getUserSession(c);
	return services.GetDataTablePerubahanData(c, usrsession)
}
