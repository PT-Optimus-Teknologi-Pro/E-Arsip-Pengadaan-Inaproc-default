package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetKajiUlang(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	if paket.Status < 4 {
		paket.Status = 4
		services.SavePaket(paket)
	}
	mp["paket"] = paket
	mp["allowJawab"] = !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool)
	return c.Render("reviu/kajiulang", mp)
}

func PublishKajiUlang(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := mp["id"].(uint)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	payload := new(models.KajiUlang)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if mp["isPokja"].(bool) {
		payload.PegId = userid
		payload.PntId = paket.PntId
	}
	if mp["isPP"].(bool) {
		payload.PpId = userid
	}
	if mp["isPPK"].(bool) {
		payload.PpkId = userid
	}
	payload.PktID = id
	if err := services.SimpanKajiUlang(c, userid, payload) ; err != nil {
		log.Error(err)
		return flashError(c, "Kirim Pertanyaan Gagal","/kajiulang/"+utils.UintToString(id))
	}
	if paket.Status < 4 { // set status paket kaji ulang
		paket.Status = 4
		services.SavePaket(paket)
	}
	return flashSuccess(c, "Kirim Pertanyaan Sukses","/kajiulang/"+utils.UintToString(id))
}

func GetJsonKajiUlang(c *fiber.Ctx)  error {
	id := utils.StringToUint(c.Params("id"))
	data := services.GetKajiUlangPaket(id)
	return c.JSON(data)
}

func GetKajiUlangList(c *fiber.Ctx)  error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	data := services.GetKajiUlangPaket(id)
	mp["datas"] = data
	mp["allowJawab"] = !mp["isAdmin"].(bool) && !mp["isUkpbj"].(bool)
	return c.Render("reviu/kajiulang-list", mp)
}

func GetJawabKajiUlang (c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	data := services.GetKajiUlang(id)
	mp["data"] = data
	return c.Render("reviu/form-jawab", mp)
}

func JawabKajiUlang(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := mp["id"].(uint)
	id := utils.StringToUint(c.Params("id"))
	payload := new(models.KajiUlang)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	pertanyaan := services.GetKajiUlang(id)
	paket := pertanyaan.Paket()
	if mp["isPokja"].(bool) {
		payload.PegId = userid
		payload.PntId = paket.PntId
	}
	if mp["isPP"].(bool) {
		payload.PpId = userid
	}
	if mp["isPPK"].(bool) {
		payload.PpkId = userid
	}
	payload.PktID = paket.ID
	payload.ParentId = id
	if err := services.SimpanKajiUlang(c, userid, payload) ; err != nil {
		log.Error(err)
		return flashError(c, "Kirim Penjelasan Gagal",	"/kajiulang/"+utils.UintToString(paket.ID))
	}
	return flashSuccess(c, "Kirim Penjelasan Sukses","/kajiulang/"+utils.UintToString(paket.ID))
}

func PreviewBA(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	data := services.GetKajiUlangPaket(id)
	
	// Fetch foto rapat dari dok_paket
	fotoRapat := models.GetDokPaketJenis(id, models.FOTO_RAPAT)
	
	mp["datas"] = data
	mp["paket"] = services.GetPaket(id)
	mp["fotoRapat"] = fotoRapat
	return c.Render("ba/ba-kajiulang", mp)
}

func UploadFotoRapat(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := mp["id"].(uint)
	id := utils.StringToUint(c.Params("id"))
	if id == 0 {
		id = utils.StringToUint(c.FormValue("pkt_id"))
	}
	if id == 0 {
		return flashError(c, "ID paket tidak ditemukan", "/paket")
	}

	if err := services.SimpanFotoRapatKajiUlang(c, id, userid); err != nil {
		log.Error(err)
		return flashError(c, "Upload foto gagal", "/kajiulang/berita-acara/"+utils.UintToString(id))
	}

	return flashSuccess(c, "Upload foto berhasil", "/kajiulang/berita-acara/"+utils.UintToString(id))
}
