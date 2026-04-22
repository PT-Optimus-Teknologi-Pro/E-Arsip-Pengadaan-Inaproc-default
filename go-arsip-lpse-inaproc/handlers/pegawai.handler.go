package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditPegawai(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/pegawai"
	id, _ := c.ParamsInt("id")
	if id > 0 {
		user := services.GetPegawai(uint(id))
		if user.ID == 0 {
			return c.SendStatus(404)
		}
		mp["pegawai"] = user
		mp["url"] = "/pegawai/edit/" + utils.IntToString(id)
	}
	return c.Render("pegawai/form-pegawai", mp)
}

func CreatePegawai(c *fiber.Ctx) error {
	user := new(models.Pegawai)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Pegawai Gagal","/pegawai/edit")
	}
	user.PegMasaBerlaku, _ = time.Parse("2006-01-02", c.FormValue("masa_berlaku"))
	user.PegIsactive = models.APPROVED
	user.PegStatus = models.APPROVED
	err = services.CreatePegawai(user)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit Pegawai Gagal", "/pegawai/edit")
	}
	return flashSuccess(c, "Tambah/Edit Pegawai Sukses","/pegawai")
}

func GetAllPegawai(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("pegawai/pegawai", mp)
}

func GetPegawai(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return c.SendStatus(404)
	}
	documents := services.GetDocumentPegawai(uint(id))
	mp["pegawai"] = user
	mp["documents"] = documents
	return c.Render("pegawai/pegawai-detil", mp)
}

// update a user in db
func UpdatePegawai(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return flashError(c, "Pegawai tidak ditemukan","/pegawai/edit/" + utils.IntToString(id))
	}
	err := c.BodyParser(&user)
	if err != nil {
		return flashError(c, "Edit Pegawai Gagal","/pegawai/edit/" + utils.IntToString(id))
	}
	user.PegIsactive = utils.StringToInt(c.FormValue("peg_isactive", "0"))
	services.UpdatePegawai(user)
	return flashSuccess(c, "Edit Pegawai Sukses","/pegawai")
}

// delete user in db by ID
func DeletePegawai(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	err := services.DeletePegawai(uint(id))
	if err != nil {
		return flashError(c,  err.Error(),"/pegawai")
	}
	return flashSuccess(c, "Hapus Pegawai Sukses","/pegawai")
}

func GetJsonPegawai(c *fiber.Ctx) error {
	sess := getSession(c)
	usrgroup := sess.Get("group")
	return services.GetDataTablePegawai(c, usrgroup.(string))
}
