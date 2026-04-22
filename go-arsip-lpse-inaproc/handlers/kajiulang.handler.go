package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"strconv"
	"time"

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
	userid := utils.InterfaceToUint(mp["id"])
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
	if paket.Status < 4 { // set status paket reviu dokumen
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
	userid := utils.InterfaceToUint(mp["id"])
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
	userid := utils.InterfaceToUint(mp["id"])
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

func GetManageBAReviu(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}

	// 1. Get existing BA Metadata or init new
	var ba models.BeritaAcara
	models.GetDB().Where("pkt_id = ? AND jenis = 'REVIU'", id).First(&ba)
	if ba.ID == 0 {
		// Auto-fill from Paket/SIRUP as requested
		ba.PktId = id
		ba.Jenis = "REVIU"
		ba.SubKeg = paket.Rup().Nama // Assuming SIRUP name is SubKeg context
		ba.Pengadaan = paket.Nama
	}

	// 2. Get master review list and saved results
	reviuMaster := services.GetAllReviu()
	var reviuResults []models.ReviuPaket
	models.GetDB().Where("pkt_id = ?", id).Find(&reviuResults)

	// Map results for easier access in template
	resMap := make(map[uint]models.ReviuPaket)
	for _, r := range reviuResults {
		resMap[r.RevId] = r
	}

	// 3. Fetch dynamic BA templates for default content
	tplDasar := models.GetTemplateByVariable("reviu_dasar")
	tplPembahasan := models.GetTemplateByVariable("reviu_pembahasan")
	tplKesimpulan := models.GetTemplateByVariable("reviu_kesimpulan")

	mp["ba"] = ba
	mp["paket"] = paket
	mp["reviuMaster"] = reviuMaster
	mp["reviuResults"] = resMap
	mp["tplDasar"] = tplDasar.Content
	mp["tplPembahasan"] = tplPembahasan.Content
	mp["tplKesimpulan"] = tplKesimpulan.Content

	// Helper for date formatting in form inputs
	if ba.Tanggal.Valid {
		mp["baTanggalStr"] = ba.Tanggal.Time.Format("2006-01-02")
	}

	return c.Render("reviu/ba-form", mp)
}

func SubmitBAReviu(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	
	// 1. Save BA Metadata
	var ba models.BeritaAcara
	models.GetDB().Where("pkt_id = ? AND jenis = 'REVIU'", id).First(&ba)
	if err := c.BodyParser(&ba); err != nil {
		return flashError(c, "Gagal memproses data BA: "+err.Error(), c.Get("Referer"))
	}
	ba.PktId = id
	ba.Jenis = "REVIU"
	
	// Handle Date manual parsing if needed
	tglStr := c.FormValue("tanggal")
	if tglStr != "" {
		t, _ := time.Parse("2006-01-02", tglStr)
		ba.Tanggal.Time = t
		ba.Tanggal.Valid = true
	}

	if err := models.GetDB().Save(&ba).Error; err != nil {
		return flashError(c, "Gagal menyimpan BA: "+err.Error(), c.Get("Referer"))
	}

	// 2. Save Checklist Results
	reviuMaster := services.GetAllReviu()
	for _, m := range reviuMaster {
		statusVal := c.FormValue("status_" + utils.UintToString(m.ID))
		ketVal := c.FormValue("ket_" + utils.UintToString(m.ID))
		
		status, _ := strconv.Atoi(statusVal)
		
		var res models.ReviuPaket
		models.GetDB().Where("pkt_id = ? AND rev_id = ?", id, m.ID).First(&res)
		res.PktId = id
		res.RevId = m.ID
		res.Status = status
		res.Keterangan = ketVal
		res.PegId = userid
		models.GetDB().Save(&res)
	}

	return flashSuccess(c, "Berita Acara Berhasil Disimpan", "/kajiulang/"+utils.UintToString(id))
}

func SignBAReviu(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	signatureData := c.FormValue("signature")

	if signatureData == "" {
		return flashError(c, "Tanda tangan tidak boleh kosong", c.Get("Referer"))
	}

	// Call signature saving service (reusing or adapting SaveTTD logic)
	// We'll save it as a specific document type or link it to the BA Signer
	err := services.SaveSignatureBA(c, id, userid, signatureData)
	if err != nil {
		return flashError(c, "Gagal menyimpan tanda tangan: "+err.Error(), c.Get("Referer"))
	}

	return flashSuccess(c, "Tanda Tangan Berhasil Disimpan", c.Get("Referer"))
}
