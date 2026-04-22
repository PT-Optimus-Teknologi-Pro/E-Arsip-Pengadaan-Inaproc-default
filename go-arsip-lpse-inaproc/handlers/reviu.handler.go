package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditReviu(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/reviu"
	id, _ := c.ParamsInt("id")
	if id != 0 {
		reviu := services.GetReviu(uint(id))
		if reviu.ID == 0 {
			return c.SendStatus(404)
		}
		mp["reviu"] = reviu
		mp["url"] = "/reviu/" + utils.IntToString(id)
	}
	mp["bidanglist"] = models.GetAllReviuBidang()
	return c.Render("reviu/form-reviu", mp)
}

func CreateReviu(c *fiber.Ctx) error {
	reviu := new(models.Reviu)
	err := c.BodyParser(reviu)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Reviu Gagal","/reviu/edit")
	}
	err = services.CreateReviu(*reviu)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Reviu Gagal","/reviu/edit")
	}
	return flashSuccess(c, "Tambah Reviu Sukses","/reviu")
}

// Get All Users from db
func GetAllReviu(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("reviu/reviu", mp)
}

// GetSingleUser from db
func GetReviu(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	reviu := services.GetReviu(uint(id))
	if reviu.ID == 0 {
		return c.SendStatus(404)
	}
	mp["reviu"] = reviu
	return c.Render("reviu/reviu-detil", mp)
}

// update a user in db
func UpdateReviu(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var formreviu models.Reviu
	err := c.BodyParser(&formreviu)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Reviu Gagal", "/reviu/edit/" + utils.IntToString(id))
	}
	err = services.UpdateReviu(uint(id), formreviu)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Reviu Gagal","/reviu/edit/" + utils.IntToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Edit Reviu Sukses","/reviu")
}

// delete user in db by ID
func DeleteReviu(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	err := services.DeleteReviu(uint(id))
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Reviu Gagal","/reviu")
	}
	return flashSuccess(c, "Hapus Reviu Sukses","/reviu")
}

func GetJsonReviu(c *fiber.Ctx) error {
	return services.GetDataTableReviu(c)
}

// NEW HANDLERS

func GetReviuPaketList(c *fiber.Ctx) error {
    mp := currentMap(c)
    // Filter for packages in Review Document status (Status 4)
    var pakets []models.Paket
    models.GetDB().Where("status = 4").Find(&pakets)
    mp["pakets"] = pakets
    return c.Render("reviu/reviu-paket", mp)
}

func GetAllReviuBidang(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["bidangs"] = models.GetAllReviuBidang()
	return c.Render("reviu/bidang-list", mp)
}

func EditReviuBidang(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	var bidang models.ReviuBidang
	if id != 0 {
		bidang = models.GetReviuBidang(uint(id))
	}
	mp["bidang"] = bidang
	return c.Render("reviu/bidang-form", mp)
}

func SaveReviuBidang(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var bidang models.ReviuBidang
	if id != 0 {
		bidang = models.GetReviuBidang(uint(id))
	}
	if err := c.BodyParser(&bidang); err != nil {
		return flashError(c, "Gagal memproses data", "/reviu/bidang")
	}
	if err := models.SaveReviuBidang(&bidang); err != nil {
		return flashError(c, "Gagal menyimpan data", "/reviu/bidang")
	}
	return flashSuccess(c, "Berhasil menyimpan Kategori", "/reviu/bidang")
}

func DeleteReviuBidang(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	bidang := models.GetReviuBidang(id)
	if err := models.DeleteReviuBidang(&bidang); err != nil {
		return flashError(c, "Gagal menghapus data", "/reviu/bidang")
	}
	return flashSuccess(c, "Berhasil menghapus Kategori", "/reviu/bidang")
}
