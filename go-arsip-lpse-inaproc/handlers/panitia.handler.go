package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditPanitia(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/pokja"
	id, _ := c.ParamsInt("id")
	mp["anggotas"] = services.GetListAnggotaPokja()
	if id != 0 {
		panitia := services.GetPanitia(uint(id))
		if panitia.ID == 0 {
			return c.SendStatus(404)
		}
		mp["panitia"] = panitia
		mp["url"] = "/pokja/" + utils.IntToString(id)
	} else {
		mp["panitia"] = models.Panitia{Tahun: time.Now().Year()}
	}
	return c.Render("panitia/form-panitia", mp)
}

func CreatePanitia(c *fiber.Ctx) error {
	panitia := new(models.PanitiaDTO)
	err := c.BodyParser(panitia)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error() ,"/pokja/edit")
	}
	err = services.CreatePanitia(*panitia)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), "/pokja/edit")
	}
	return flashSuccess(c, "Tambah Panitia Sukses","/pokja")
}

// Get All Users from db
func GetAllPanitia(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("panitia/panitia", mp)
}

// GetSingleUser from db
func GetPanitia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	panitia := services.GetPanitia(uint(id))
	if panitia.ID == 0 {
		return c.SendStatus(404)
	}
	anggotas := services.GetAnggotaPokja(uint(id))
	mp["panitia"] = panitia
	mp["anggotas"] = anggotas
	return c.Render("panitia/panitia-detil", mp)
}

// update a user in db
func UpdatePanitia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var panitia models.PanitiaDTO
	err := c.BodyParser(&panitia)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Panitia Gagal", "/pokja/edit/" + utils.IntToString(id))
	}
	log.Info("DTO : ", panitia.Anggota)
	err = services.UpdatePanitia(uint(id), panitia)
	if err != nil {
		return flashError(c, err.Error(), "/pokja/edit/" + utils.IntToString(id))
	}
	return flashSuccess(c, "Edit Panitia Sukses","/pokja")
}

// delete user in db by ID
func DeletePanitia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	err := services.DeletePanitia(uint(id))
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Panitia Gagal","/pokja")
	}
	return flashSuccess(c, "Hapus Panitia Sukses","/pokja")
}

func GetJsonPanitia(c *fiber.Ctx) error {
	return services.GetDataTablePanitia(c)
}
