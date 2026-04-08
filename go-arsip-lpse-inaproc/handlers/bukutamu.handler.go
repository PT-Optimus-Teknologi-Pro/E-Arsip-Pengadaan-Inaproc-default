package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func BukuTamu(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["userPengadaan"] = services.GetPegawaiPengadaan()
	mp["userNonPengadaan"] = services.GetPegawaiNonPengadaan()
	return c.Render("bukutamu/bukutamu", mp)
}

func BukuTamuView(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	buku := services.GetBukuTamu(uint(id))
	if buku.ID == 0 {
		return c.SendStatus(404)
	}
	mp["buku"] = buku
	mp["pegawaiList"] = services.GetPegawaiPengadaan()
	return c.Render("bukutamu/bukutamu-view", mp)
}

func SubmitBukuTamu(c *fiber.Ctx) error {
	buku := new(models.BukuTamu)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(buku)
	if err != nil {
		log.Error(err)
		return flashError(c, "Pendaftaran Buku Tamu Gagal","/bukutamu")
	}
	// mp := currentMap(c)
	// userid := mp["id"].(uint)
	id, err := services.SaveDocBukuTamu(c, buku, "gambar")
	if err != nil {
		log.Error(err)
		return flashError(c, "Pendaftaran Buku Tamu Gagal","/bukutamu")
	}
	return flashSuccessWithData(c, "Pendaftaran buku tamu Sukses", utils.UintToString(id), "/bukutamu/success")
}

func BukutamuKonfirmasi(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["id_tamu"] = c.Params("id")
	return c.Render("bukutamu/bukutamu-konfirmasi", mp)
}

func BukuTamuList(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("bukutamu/bukutamu-list", mp)
}

func GetJsonBukuTamu(c *fiber.Ctx) error {
	mp := currentMap(c)
	isUkpbj := mp["isUkpbj"].(bool)
	return services.GetDataTableBukuTamu(c, isUkpbj)
}

func BukuTamuUpdate(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	status := utils.StringToInt(c.FormValue("status"))
	buku := services.GetBukuTamu(uint(id))
	if buku.ID == 0 {
		return c.SendStatus(404)
	}
	buku.Status = status
	buku.PegId = mp["id"].(uint)
	if buku.Status == 1 {
		buku.TglProses = time.Now()
	} else if buku.Status == 2 {
		buku.TglSelesai = time.Now()
	}
	_, err := services.SaveBukuTamu(&buku)
	if err != nil {
		log.Error(err)
		return flashError(c, "Update buku tamu Gagal","/bukutamu/"+utils.IntToString(id))
	}
	return flashSuccess(c, "Update buku tamu Sukses", "/bukutamu/"+utils.IntToString(id))
}
func SubmitProsesBukuTamu(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	peg_id := utils.StringToInt(c.FormValue("peg_id"))
	catatan := c.FormValue("catatan")
	tindaklanjut := c.FormValue("tindaklanjut")
	buku := services.GetBukuTamu(uint(id))
	if buku.ID == 0 {
		return c.SendStatus(404)
	}
	buku.Status = 1
	buku.PegId = uint(peg_id)
	buku.Catatan = catatan
	buku.TindakLanjut = tindaklanjut
	if buku.Status == 1 {
		buku.TglProses = time.Now()
	}
	_, err := services.SaveBukuTamu(&buku)
	if err != nil {
		log.Error(err)
		return flashError(c, "Update buku tamu Gagal","/bukutamu/"+utils.IntToString(id))
	}
	return flashSuccess(c, "Update buku tamu Sukses", "/bukutamu/"+utils.IntToString(id))
}
