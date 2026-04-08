package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditUkpbj(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/ukpbj"
	id, _ := c.ParamsInt("id")
	if id != 0 {
		ukpbj := services.GetUkpbj(uint(id))
		if ukpbj.ID == 0 {
			return c.SendStatus(404)
		}
		mp["ukpbj"] = ukpbj
		mp["url"] = "/ukpbj/edit/" + utils.IntToString(id)
	}
	mp["adminUkpbj"] = models.GetAdminUKPBJ()
	return c.Render("ukpbj/form-ukpbj", mp)
}

func CreateUkpbj(c *fiber.Ctx) error {
	ukpbj := new(models.Ukpbj)
	err := c.BodyParser(ukpbj)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit UKPBJ Gagal","/ukpbj/edit")
	}
	log.Info("pegawai : ", ukpbj.PegId)
	ukpbj.TglDaftar, _ = time.Parse("2006-01-02", c.FormValue("tglDaftar"))
	err = services.CreateUkpbj(*ukpbj)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah/Edit UKPBJ Gagal", "/ukpbj/edit")
	}
	return flashSuccess(c, "Tambah/Edit UKPBJ Sukses","/ukpbj")
}

func GetAllUkpbj(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["count"] = models.GetCountUkpbj()
	return c.Render("ukpbj/ukpbj", mp)
}

func GetUkpbj(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	ukpbj := services.GetUkpbj(uint(id))
	if ukpbj.ID == 0 {
		return c.SendStatus(404)
	}
	mp := currentMap(c)
	mp["ukpbj"] = ukpbj
	return c.Render("ukpbj/ukpbj-detil", mp)
}

func UpdateUkpbj(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	ukpbj := services.GetUkpbj(uint(id))
	if ukpbj.ID == 0 {
		return flashError(c, "UKPBJ tidak ditemukan", "/ukpbj/edit/" + utils.IntToString(id))
	}
	var objUkpbj models.Ukpbj
	err := c.BodyParser(&objUkpbj)
	log.Info("ukpbj", objUkpbj)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit UKPBJ Gagal", "/ukpbj/edit/" + utils.IntToString(id))
	}
	ukpbj.Nama = objUkpbj.Nama
	ukpbj.Alamat = objUkpbj.Alamat
	ukpbj.Telepon = objUkpbj.Telepon
	ukpbj.Fax = objUkpbj.Fax
	ukpbj.IsActive = objUkpbj.IsActive
	ukpbj.PegId = objUkpbj.PegId
	services.SaveUkpbj(ukpbj)
	return flashSuccess(c, "Edit UKPBJ Sukses", "/ukpbj")
}

// delete user in db by ID
func DeleteUkpbj(c *fiber.Ctx) error {
	id,_ := c.ParamsInt("id")
	ukpbj := services.GetUkpbj(uint(id))
	if ukpbj.ID == 0 {
		return flashError(c, "UKPBJ tidak ditemukan","/ukpbj")
	}
	err := services.DeleteUkpbj(ukpbj)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus UKPBJ Gagal","/ukpbj")
	}
	return flashSuccess(c, "Hapus UKPBJ Sukses", "/ukpbj")
}

func GetJsonUkpbj(c *fiber.Ctx) error {
	return services.GetDataTableUkpbj(c)
}
