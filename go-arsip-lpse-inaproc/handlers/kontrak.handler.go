package handlers

import (
	"arsip/services"
	"arsip/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HasilPekerjaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	mp["paket"] = paket
	switch paket.Metode {
	case 8:
		realisasi := paket.GetNontender().GetRealisasi()
		if len(realisasi) > 0 {
			log.Info("realisi ", realisasi[0])
			mp["realisasi"] = realisasi[0]
		}
	case 9:
		mp["purchase"] = paket.GetPurchase()
	default:
		realisasi := paket.GetTender().GetRealisasi()
		if len(realisasi) > 0 {
			log.Info("realisi ", realisasi[0])
			mp["realisasi"] = realisasi[0]
		}
	}
	return c.Render("paket/hasil-pekerjaan", mp)
}

func SimpanHasilPekerjaan(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	status := utils.StringToUint(c.FormValue("status"))
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	paket.Status = int(status)
	err := services.SavePaket(paket)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Hasil Pekerjaan Gagal", "/hasil-pekerjaan/"+utils.UintToString(id))
	}
	return flashSuccess(c, "Simpan Hasil Pekerjaan Berhasil", "/hasil-pekerjaan/"+utils.UintToString(id))
}

func KontrakPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	mp["paket"] = paket
	switch paket.Metode {
	case 8:
		realisasi := paket.GetNontender().GetRealisasi()
		if len(realisasi) > 0 {
			log.Info("realisi ", realisasi[0])
			mp["realisasi"] = realisasi[0]
		}
	case 9:
		mp["purchase"] = paket.GetPurchase()
	default:
		realisasi := paket.GetTender().GetRealisasi()
		if len(realisasi) > 0 {
			log.Info("realisi ", realisasi[0])
			mp["realisasi"] = realisasi[0]
		}
	}
	return c.Render("paket/kontrak", mp)
}

func SimpanKontrakPaket(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	status := utils.StringToUint(c.FormValue("status"))
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	paket.Status = int(status)
	err := services.SavePaket(paket)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Kontrak Gagal", "/kontrak/"+utils.UintToString(id))
	}
	return flashSuccess(c, "Simpan Kontrak Berhasil", "/kontrak/"+utils.UintToString(id))
}

func SimpanDokKontrak(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanDokKontrak(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Dokumen Pendukung pengadaan Gagal", "/kontrak/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan  Dokumen Pendukung pengadaan Sukses", "/kontrak/"+utils.UintToString(id))
}

func SimpanDokHasilPekerjaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := mp["id"].(uint)
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanDokHasilPekerjaan(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Dokumen Pendukung pengadaan Gagal", "/hasil-pekerjaan/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan  Dokumen Pendukung pengadaan Sukses", "/hasil-pekerjaan/"+utils.UintToString(id))
}

func DownloadKontrak(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id") // id paket
	paket := services.GetPaket(uint(id))
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	log.Info("create zip file");
	var files []string
	for _,v := range paket.DokKontrakList() {
		files = append(files, v.Document().Filepath)
	}
	zipFile,  err := utils.CreateZip(files, "dok-kontrak.zip");
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}

func DownloadPekerjaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id") // id paket
	paket := services.GetPaket(uint(id))
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	log.Info("create zip file");
	var files []string
	for _,v := range paket.DokHasilList() {
		files = append(files, v.Document().Filepath)
	}
	zipFile,  err := utils.CreateZip(files, "dok-hasil-pekerjaan.zip");
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}
