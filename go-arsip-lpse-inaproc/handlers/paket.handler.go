package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateManualPaketForm(c *fiber.Ctx) error {
	mp := currentMap(c)
	usersession := getUserSession(c)
	pegawai := services.GetPegawai(usersession.Id)
	if !usersession.IsArsiparis() && !usersession.IsPpk() && !usersession.IsAdmin() && !usersession.IsUkpbj() && !usersession.IsPokja() && !usersession.IsPp() {
		return Forbiden(c)
	}
	if !pegawai.IsApprove() {
		return Forbiden(c)
	}
	mp["satkers"] = services.GetAllSatker()
	mp["tahuns"] = services.GetTahunRupList()
	// Pass caller role context so the view can render correct back button
	if usersession.IsPpk() {
		mp["callerRole"] = "ppk"
	} else if usersession.IsArsiparis() {
		mp["callerRole"] = "arsiparis"
	} else if usersession.IsUkpbj() {
		mp["callerRole"] = "ukpbj"
	} else if usersession.IsPokja() {
		mp["callerRole"] = "pokja"
	} else if usersession.IsPp() {
		mp["callerRole"] = "pp"
	} else {
		mp["callerRole"] = "admin"
	}
	return c.Render("paket/form-manual", mp)
}

func SaveManualPaket(c *fiber.Ctx) error {
	usersession := getUserSession(c)
	pegawai := services.GetPegawai(usersession.Id)
	if !usersession.IsArsiparis() && !usersession.IsPpk() && !usersession.IsAdmin() && !usersession.IsUkpbj() && !usersession.IsPokja() && !usersession.IsPp() {
		return Forbiden(c)
	}
	if !pegawai.IsApprove() {
		return Forbiden(c)
	}
	
	keterangan := c.FormValue("keterangan")
	nama := c.FormValue("nama")
	if nama == "" {
		// Use first line or first 100 chars of keterangan as nama
		nama = keterangan
		if len(nama) > 100 {
			nama = nama[:100] + "..."
		}
		if nama == "" {
			nama = "Dokumen Privat"
		}
	}
	
	// Default values for removed fields
	tahun := time.Now().Year()
	pagu := 0.0
	hps := 0.0
	metodeStr := "-"
	jenis := "-"
	satkerId := 0

	// Auto-assign ppkId: if creator is PPK, use their own ID
	var ppkId uint
	if usersession.IsPpk() {
		ppkId = usersession.Id
	}

	paketId, err := services.CreateManualPaket(c, usersession.Id, ppkId, nama, tahun, pagu, hps, 0, metodeStr, jenis, uint(satkerId), keterangan)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal simpan paket manual: "+err.Error(), "/paket/create-manual")
	}

	// Redirect back to Dokumen Privat for all relevant roles
	if usersession.IsPpk() || usersession.IsUkpbj() || usersession.IsPokja() || usersession.IsPp() {
		return flashSuccess(c, "Paket manual berhasil dibuat", "/"+strings.ToLower(usersession.Role)+"/dokumen-privat")
	}
	return flashSuccess(c, "Paket manual berhasil dibuat", "/paket/"+utils.UintToString(paketId))
}

func EditManualPaketForm(c *fiber.Ctx) error {
	mp := currentMap(c)
	usersession := getUserSession(c)
	pegawai := services.GetPegawai(usersession.Id)
	if !usersession.IsArsiparis() && !usersession.IsPpk() && !usersession.IsAdmin() && !usersession.IsUkpbj() && !usersession.IsPokja() && !usersession.IsPp() {
		return Forbiden(c)
	}
	if !pegawai.IsApprove() {
		return Forbiden(c)
	}
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	// PPK can only edit their own manual packages
	if usersession.IsPpk() && paket.PpkId != usersession.Id {
		return Forbiden(c)
	}

	mp["paket"] = paket
	mp["satkers"] = services.GetAllSatker()
	// Pass existing evidence docs so the form can show them
	mp["listBukti"] = models.GetDokPaketJenisList(paket.ID, "Bukti Manual")
	if usersession.IsPpk() {
		mp["callerRole"] = "ppk"
	} else if usersession.IsArsiparis() {
		mp["callerRole"] = "arsiparis"
	} else if usersession.IsUkpbj() {
		mp["callerRole"] = "ukpbj"
	} else if usersession.IsPokja() {
		mp["callerRole"] = "pokja"
	} else if usersession.IsPp() {
		mp["callerRole"] = "pp"
	} else {
		mp["callerRole"] = "admin"
	}
	return c.Render("paket/form-manual-edit", mp)
}

func UpdateManualPaket(c *fiber.Ctx) error {
	usersession := getUserSession(c)
	if !usersession.IsArsiparis() && !usersession.IsPpk() && !usersession.IsAdmin() && !usersession.IsUkpbj() && !usersession.IsPokja() && !usersession.IsPp() {
		return Forbiden(c)
	}

	id := utils.StringToUint(c.Params("id"))
	// PPK can only edit their own manual packages
	if usersession.IsPpk() {
		paket := services.GetPaket(id)
		if paket.PpkId != usersession.Id {
			return Forbiden(c)
		}
	}

	keterangan := c.FormValue("keterangan")
	nama := c.FormValue("nama")
	if nama == "" {
		nama = keterangan
		if len(nama) > 100 {
			nama = nama[:100] + "..."
		}
		if nama == "" {
			nama = "Dokumen Privat"
		}
	}
	
	// Default values for removed fields
	tahun := time.Now().Year()
	pagu := 0.0
	hps := 0.0
	metodeStr := "-"
	jenis := "-"
	satkerId := 0

	err := services.UpdateManualPaket(c, id, usersession.Id, nama, tahun, pagu, hps, 0, metodeStr, jenis, uint(satkerId), keterangan)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal update paket manual: "+err.Error(), "/paket/edit-manual/"+utils.UintToString(id))
	}

	if usersession.IsPpk() || usersession.IsUkpbj() || usersession.IsPokja() || usersession.IsPp() {
		return flashSuccess(c, "Paket manual berhasil diperbarui", "/paket/"+utils.UintToString(id))
	}
	return flashSuccess(c, "Paket manual berhasil diperbarui", "/paket/"+utils.UintToString(id))
}

func DownloadBuktiZip(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}

	listBukti := models.GetDokPaketJenisList(id, "Bukti Manual")
	if len(listBukti) == 0 {
		return flashError(c, "Tidak ada dokumen bukti untuk didownload", c.Get("Referer"))
	}

	// Just download the single file if there's only one
	if len(listBukti) == 1 {
		doc := listBukti[0].Document()
		return c.Download(doc.Filepath, doc.Filename)
	}

	var files []string
	for _, b := range listBukti {
		files = append(files, b.Document().Filepath)
	}

	// Use package name for the zip filename
	safeName := strings.ReplaceAll(paket.Nama, " ", "_")
	// Simple regex replacement for other special chars if needed, or just use strings.Map
	safeName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, safeName)

	zipName := "bukti-" + safeName + "-" + utils.UintToString(id) + ".zip"
	zipFile, err := utils.CreateZip(files, zipName)
	if err != nil {
		log.Error("Error creating zip: ", err)
		return flashError(c, "Gagal membuat file zip: "+err.Error(), c.Get("Referer"))
	}

	return c.Download(zipFile, zipName)
}

func EditPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/paket"
	id := utils.StringToUint(c.Params("id"))
	if id != 0 {
		paket := services.GetPaket(id)
		if paket.ID == 0 {
			return c.SendStatus(404)
		}
		mp["paket"] = paket
		mp["url"] = "/paket/" + utils.UintToString(id)
		mp["metodelist"] = models.GetAllMetodePengadaan()
		mp["jenislist"] = models.GetAllJenisPengadaan()
		// Get all satkers for the current package year (or current year)
		tahun := paket.Tahun
		if tahun == 0 {
			tahun = time.Now().Year()
		}
		mp["satkers"] = models.GetAllSatkerSirup(tahun)
		return c.Render("paket/form-paket", mp)
	} else {
		now := time.Now().Year()
		tahun := c.QueryInt("tahun", now)
		satker := c.Query("satker")
		mp["tahun"] = tahun
		mp["tahunlist"] = services.GetTahunRupList()
		mp["satker"] = satker
		mp["metode"] = c.Query("metode")
		mp["jenis"] = c.Query("jenis")
		mp["jenislist"] = models.GetAllJenisPengadaan()
		mp["metodelist"] = models.GetAllMetodePengadaan();
		return c.Render("paket/rencana-paket", mp)
	}
}

func CreatePaket(c *fiber.Ctx) error {
	rupid := utils.StringToUint(c.FormValue("id"))
	log.Info("create paket from rup ", rupid)
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	pegawai := services.GetPegawai(userid)
	if !pegawai.IsApprove() {
		return Forbiden(c)
	}
	paketId, err := services.CreatePaket(rupid, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Paket Gagal: "+err.Error(), "/paket/edit")
	}
	return flashSuccess(c, "Tambah Paket Sukses", "/paket/"+strconv.Itoa(int(paketId)))
}

// Get All Users from db
func GetAllPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	usersession := getUserSession(c)
	user := usersession.Pegawai()
	mp["allowBuatPaket"] = user.IsApprove() && (usersession.IsPpk() || usersession.IsAdmin()) && user.IsAktif()

	if isArsiparis, ok := mp["isArsiparis"].(bool); ok && isArsiparis {
		var tenderMethods []string
		var tenderJenis []string
		models.GetDB().Table("tender").Joins("JOIN tender_selesai ts ON tender.kd_tender = ts.kd_tender").Where("tender.mtd_pemilihan IS NOT NULL AND tender.mtd_pemilihan != ''").Distinct("tender.mtd_pemilihan").Pluck("tender.mtd_pemilihan", &tenderMethods)
		models.GetDB().Table("tender").Joins("JOIN tender_selesai ts ON tender.kd_tender = ts.kd_tender").Where("tender.jenis_pengadaan IS NOT NULL AND tender.jenis_pengadaan != ''").Distinct("tender.jenis_pengadaan").Pluck("tender.jenis_pengadaan", &tenderJenis)

		var nonTenderMethods []string
		var nonTenderJenis []string
		models.GetDB().Table("nontender").Joins("JOIN nontender_selesai ts ON nontender.kd_nontender = ts.kd_nontender").Where("nontender.mtd_pemilihan IS NOT NULL AND nontender.mtd_pemilihan != ''").Distinct("nontender.mtd_pemilihan").Pluck("nontender.mtd_pemilihan", &nonTenderMethods)
		models.GetDB().Table("nontender").Joins("JOIN nontender_selesai ts ON nontender.kd_nontender = ts.kd_nontender").Where("nontender.jenis_pengadaan IS NOT NULL AND nontender.jenis_pengadaan != ''").Distinct("nontender.jenis_pengadaan").Pluck("nontender.jenis_pengadaan", &nonTenderJenis)

		// E-Purchasing data is pulled from katalog table via UNION in datatable model, so we must add it as an option.
		nonTenderMethods = append(nonTenderMethods, "E-Purchasing")
		nonTenderJenis = append(nonTenderJenis, "Barang/Jasa")

		mp["tenderMethods"] = tenderMethods
		mp["tenderJenis"] = tenderJenis
		mp["nonTenderMethods"] = nonTenderMethods
		mp["nonTenderJenis"] = nonTenderJenis
	}

	return c.Render("paket/paket", mp)
}

// GetSingleUser from db
func GetPaket(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	if paket.IsDraft() && paket.CountChecklist() == 0 {
		paket.GeneratePersyaratan()
	}
	mp["paket"] = paket
	mp["anggaran"] = paket.Rup().AnggaranLabel()
	mp["allowAjukan"] = paket.IsAllowAjukan()
	mp["panitias"] = services.GetPanitias()
	mp["pps"] = services.GetPPs(paket.SatkerId)
	mp["ppks"] = services.GetPPKs()
	mp["metodelist"] = models.GetAllMetodePengadaan()
	mp["jenislist"] = models.GetAllJenisPengadaan()
	mp["prosesOnlyPpk"] = paket.IsOnlyPpk()
	mp["bukti"] = models.GetDokPaketJenis(paket.ID, "Bukti Manual")
	// For manual packages, also pass the bukti document details directly
	if paket.RupId == 0 {
		mp["listBukti"] = models.GetDokPaketJenisList(paket.ID, "Bukti Manual")
	}
	
	// Dynamic back URL for manual packages
	backUrl := "/paket"
	if paket.RupId == 0 {
		group := ""
		if mp["group"] != nil {
			group = strings.ToLower(utils.InterfaceToString(mp["group"]))
		}
		if group == "ukpbj" || group == "ppk" || group == "pokja" || group == "pp" || group == "arsiparis" {
			backUrl = "/" + group + "/dokumen-privat"
		}
	}
	mp["backUrl"] = backUrl

	// Manual/private packages get their own dedicated detail view
	if paket.RupId == 0 {
		return c.Render("paket/dokumen-privat-detil", mp)
	}
	return c.Render("paket/paket-detil", mp)
}

// update a user in db
func UpdatePaket(c *fiber.Ctx) error {
	var formPaket models.Paket
	
	// get id params
	id := utils.StringToUint(c.Params("id"))
	err := c.BodyParser(&formPaket)
	if err != nil {
		log.Error(err)
		res := fmt.Sprintf("%d", id)
		return flashError(c, "Edit Paket Gagal", "/paket/" + res)
	}

	// Manual parsing for formatted currency strings
	paguStr := c.FormValue("pagu_str")
	if paguStr != "" {
		paguStr = strings.ReplaceAll(paguStr, ".", "")
		f, _ := strconv.ParseFloat(paguStr, 64)
		formPaket.Pagu = f
	}

	hpsStr := c.FormValue("hps_str")
	if hpsStr != "" {
		hpsStr = strings.ReplaceAll(hpsStr, ".", "")
		f, _ := strconv.ParseFloat(hpsStr, 64)
		formPaket.Hps = f
	}

	// Update the rest of the fields that BodyParser might miss or we want to be explicit about
	formPaket.Metode = utils.StringToInt(c.FormValue("metode"))
	formPaket.KgrId = utils.StringToInt(c.FormValue("kgr_id"))
	formPaket.SatkerId = uint(utils.StringToInt(c.FormValue("satker_id")))

	services.SimpanPaket(id, formPaket)
	// Return the updated user
	return flashSuccess(c, "Edit Paket Sukses", "/paket")
}

func UpdateMetodePaket(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	metode := utils.StringToInt(c.FormValue("metode"))
	
	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	
	paket.Metode = metode
	err := models.SavePaket(&paket)
	if err != nil {
		return flashError(c, "Koreksi Gagal: "+err.Error(), "/paket/"+utils.UintToString(id))
	}
	
	// Re-generate checklist if needed
	_ = paket.GeneratePersyaratan()
	
	return flashSuccess(c, "Metode berhasil dikoreksi", "/paket/"+utils.UintToString(id))
}

// delete user in db by ID
func DeletePaket(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)

	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = c.Get("Referer", "/paket")
	}

	if paket.ID == 0 {
		return flashError(c, "Paket tidak ditemukan", redirect)
	}

	mp := currentMap(c)

	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}

	err := services.HapusPaket(id)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), redirect)
	}
	return flashSuccess(c, "Hapus Paket Sukses", redirect)
}

func GetJsonPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	isPPk := mp["isPPK"].(bool)
	isUkpbj := mp["isUkpbj"].(bool)
	isPp  := mp["isPP"].(bool)
	isPokja := mp["isPokja"].(bool)
	isArsiparis := mp["isArsiparis"].(bool)
	return services.GetDataTablePaket(c, userid, isPPk, isUkpbj, isPokja, isPp, isArsiparis)
}

func GetJsonTenderArsiparis(c *fiber.Ctx) error {
	mp := currentMap(c)
	isArsiparis := mp["isArsiparis"].(bool)
	if !isArsiparis {
		return Forbiden(c)
	}
	return models.GetDataTableTenderArsiparis(c)
}

func GetDetailTenderArsiparis(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	isArsiparis := mp["isArsiparis"].(bool)
	if !isArsiparis {
		return Forbiden(c)
	}

	var tender models.Tender
	models.GetDB().First(&tender, id)
	if tender.KdTender == 0 {
		return c.Status(404).SendString("Tender not found")
	}

	var tenderSelesai models.TenderSelesai
	models.GetDB().First(&tenderSelesai, id)

	// Fetch if any paket exists locally with this KodeTender
	var paket models.Paket
	models.GetDB().Where("kode_tender = ?", id).First(&paket)

	mp["tender"] = tender
	mp["tenderSelesai"] = tenderSelesai
	if paket.ID > 0 {
		mp["paket"] = paket
		mp["dok_persiapan"] = models.GetDokPaketJenis(paket.ID, "Dokumen Persiapan")
		mp["dok_hasil"] = models.GetDokPaketJenis(paket.ID, "Dokumen Hasil Pengadaan")
		mp["dok_kontrak"] = models.GetDokPaketJenis(paket.ID, "Dokumen Kontrak")
		mp["dok_tambahan"] = models.GetDokPaketJenis(paket.ID, "Dokumen Tambahan")
	}

	return c.Render("paket/detail-tender-arsiparis", mp)
}

func GetJsonNontenderArsiparis(c *fiber.Ctx) error {
	mp := currentMap(c)
	isArsiparis := mp["isArsiparis"].(bool)
	if !isArsiparis {
		return Forbiden(c)
	}
	return models.GetDataTableNontenderArsiparis(c)
}

func GetDetailNontenderArsiparis(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	isArsiparis := mp["isArsiparis"].(bool)
	if !isArsiparis {
		return Forbiden(c)
	}

	var nontender models.Nontender
	models.GetDB().First(&nontender, id)
	if nontender.KdNontender == 0 {
		return c.Status(404).SendString("Non-Tender not found")
	}

	var nontenderSelesai models.NontenderSelesai
	models.GetDB().First(&nontenderSelesai, id)

	// Fetch if any paket exists locally with this KodeNontender
	var paket models.Paket
	models.GetDB().Where("kode_tender = ?", id).First(&paket)

	mp["tender"] = nontender
	mp["tenderSelesai"] = nontenderSelesai
	mp["isNontender"] = true
	if paket.ID > 0 {
		mp["paket"] = paket
		mp["dok_persiapan"] = models.GetDokPaketJenis(paket.ID, "Dokumen Persiapan")
		mp["dok_hasil"] = models.GetDokPaketJenis(paket.ID, "Dokumen Hasil Pengadaan")
		mp["dok_kontrak"] = models.GetDokPaketJenis(paket.ID, "Dokumen Kontrak")
		mp["dok_tambahan"] = models.GetDokPaketJenis(paket.ID, "Dokumen Tambahan")
	}

	return c.Render("paket/detail-tender-arsiparis", mp)
}

func UpdateHpsPaket(c *fiber.Ctx) error {
	formHps := c.FormValue("Hps")
	formHps = strings.Replace(formHps, ".", "", -1)
	hps, _ := strconv.ParseFloat(formHps, 64)
	// get id params
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanHpsPaket(id, hps)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Update Hps Sukses", "/paket/"+utils.UintToString(id))
}

func SimpanPersyaratanPaket(c *fiber.Ctx) error {
	log.Info("SimpanPersyaratanPaket....")
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanPersyaratanPaket(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Persyaratan Gagal", "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan Persyaratan Sukses", "/paket/"+utils.UintToString(id))
}

func ApprovePaket(c *fiber.Ctx) error {
	// get id params
	id := utils.StringToUint(c.Params("id"))
	approve, _ := strconv.ParseBool(c.FormValue("approve"))
	reject, _ := strconv.ParseBool(c.FormValue("reject"))
	alasan := c.FormValue("alasan")
	prioritas, _ := strconv.ParseBool(c.FormValue("prioritas"))
	err := services.ApprovePaket(id, approve, reject, alasan, prioritas)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Paket Gagal","/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan Paket Berhasil", "/paket/"+utils.UintToString(id))
}

func KirimPaketUkpbj(c *fiber.Ctx) error {
	// get id params
	id := utils.StringToUint(c.Params("id"))
	err := services.KirimPaketUkpbj(id)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Berhasil Kirim ke UKPBJ","/paket/"+utils.UintToString(id))
}

func KirimPaketPokja(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	pntId := utils.StringToUint(c.FormValue("pnt_id"))
	err := services.AssignPaketPokja(id, pntId)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal Assign Pokja", "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Berhasil Assign Pokja","/paket/"+utils.UintToString(id))
}

func KirimPaketPp(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	ppId := utils.StringToUint(c.FormValue("pp_id"))
	err := services.AssignPaketPp(id, ppId)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal Assign PP", "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Berhasil Assign PP", "/paket/"+utils.UintToString(id))
}

func PilihPokjaPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["paketId"] = c.Params("id")
	return c.Render("paket/pilih-pokja", mp)
}

func PilihPPPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["paketId"] = c.Params("id")
	return c.Render("paket/pilih-pp", mp)
}

func SuratPenunjukan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	mp["paket"] = paket
	if paket.PpId != 0 {
		mp["pegawai"] = paket.Pp()
		return c.Render("paket/surat-penunjukan-pp", mp)
	}
	mp["panitia"] = paket.Pokja()
	return c.Render("paket/surat-penunjukan-pokja", mp)
}

func DokPersiapanPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	dokPersiapans := paket.DokPersiapan()
	allowCetak := true
	approvedCategoryMap := make(map[uint]bool)
	for i := range dokPersiapans {
		dokPersiapans[i].CheckPersetujuanPegawai()
		if !dokPersiapans[i].IsSudahPersetujuanSemua() {
			allowCetak = false
		} else {
			// Jika dokumen ini sudah final, maka kategorinya (ChkId) dianggap sudah disetujui
			approvedCategoryMap[dokPersiapans[i].ChkId] = true
		}
	}
	mp["approvedCategoryMap"] = approvedCategoryMap
	isLocked := paket.IsLockedReview && !paket.IsAddendum
	// Fetch foto rapat dari dok_paket
	fotoRapat := models.GetDokPaketJenis(id, models.FOTO_RAPAT)

	mp["isLocked"] = isLocked
	mp["isAddendum"] = paket.IsAddendum
	mp["allowCetak"] = allowCetak
	mp["paket"] = paket
	mp["dokPersiapan"] = dokPersiapans
	mp["allowUpload"] = (mp["isPPK"].(bool) || mp["isPokja"].(bool) || mp["isPP"].(bool)) && !isLocked
	mp["fotoRapat"] = fotoRapat
	mp["addendums"] = models.GetReviewAddendumList(id)
	mp["addendumCount"] = paket.AddendumCount
	mp["addendumMaxReached"] = paket.AddendumCount >= 3
	
	// Fetch all snapshots for history
	var historySnapshots []models.ReviewAddendumSnapshot
	models.GetDB().Where("pkt_id = ?", id).Order("id DESC").Find(&historySnapshots)
	mp["historySnapshots"] = historySnapshots

	return c.Render("paket/dok-persiapan", mp)
}

func SimpanDokumenPersiapanPaket(c *fiber.Ctx) error {
	log.Info("SimpanDokumenPersiapanPaket....")
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	
	// Locking Logic: Cannot save if already locked OR if any users have already approved
	if paket.IsLockedReview && !paket.IsAddendum {
		return flashError(c, "Dokumen sudah dikunci(Final). Silakan aktifkan Addendum untuk mengubah.", "/dok-final/"+utils.UintToString(id))
	}
	err := services.SimpanDokPersiapanPaket(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Dokumen Persiapan gagal", "/dok-final/" + utils.UintToString(id))
	}
	return flashSuccess(c, "Simpan Dokumen Persiapan Sukses","/dok-final/"+utils.UintToString(id))
}

func HapusDokPersiapanSnapshot(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var snap models.ReviewAddendumSnapshot
	models.GetDB().First(&snap, id)
	if snap.ID > 0 {
		models.GetDB().Delete(&snap)
	}
	return c.Redirect("/dok-final/" + utils.UintToString(snap.PktId))
}

func UnlockAddendumReview(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	reason := c.Query("reason")
	if reason == "" {
		reason = c.FormValue("reason")
	}

	if !mp["isPPK"].(bool) && !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}

	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}

	// Cek batas maksimum addendum (3x)
	const maxAddendum = 3
	if paket.AddendumCount >= maxAddendum {
		return flashError(c, fmt.Sprintf("Addendum telah mencapai batas maksimum (%d kali). Hubungi Pokja untuk mereset kuota addendum.", maxAddendum), "/dok-final/"+utils.UintToString(id))
	}

	// 1. Create Addendum Record
	var currentVersion int
	models.GetDB().Model(&models.ReviewAddendum{}).Where("pkt_id = ?", id).Count(new(int64))
	models.GetDB().Raw("SELECT COALESCE(MAX(version), 0) FROM review_addendum WHERE pkt_id = ?", id).Scan(&currentVersion)
	
	addendum := models.ReviewAddendum{
		PktId:     id,
		Version:   currentVersion + 1,
		Reason:    reason,
		CreatedBy: userid,
	}
	models.GetDB().Save(&addendum)

	// 2. Snapshot current DokPersiapan & Approvals
	dokPersiapans := paket.DokPersiapan()
	for _, dp := range dokPersiapans {
		// Capture approvals as JSON
		approvals := dp.Persetujuan()
		type ApprovalState struct {
			Name   string `json:"name"`
			Status bool   `json:"status"`
		}
		var state []ApprovalState
		for _, a := range approvals {
			state = append(state, ApprovalState{
				Name:   a.Pegawai().PegNama,
				Status: a.Status,
			})
		}
		stateJson, _ := json.Marshal(state)

		snapshot := models.ReviewAddendumSnapshot{
			AddendumId: addendum.ID,
			ChkId:      dp.ChkId,
			DokId:      dp.DokId,
			Approvals:  string(stateJson),
		}
		models.GetDB().Save(&snapshot)

		// 3. Reset Approvals for this document
		models.DeleteAllPersetujuanDokPersiapan(dp.ID)
	}

	// 4. Update Paket status - increment addendum counter
	paket.IsLockedReview = false
	paket.IsAddendum = true
	paket.AddendumCount = paket.AddendumCount + 1
	models.SavePaket(&paket)

	// 5. Send Inbox Notifications to relevant parties
	subject := "Pemberitahuan Addendum Reviu Dokumen: " + paket.Nama
	content := fmt.Sprintf(
		"Addendum ke-%d telah diaktifkan untuk paket %s.\n\nAlasan Perubahan:\n%s\n\nSilakan buka tab Dokumen Final untuk melakukan peninjauan dan persetujuan ulang.",
		paket.AddendumCount+1, paket.Nama, reason,
	)
	
	services.SendNotification(paket.PpkId, subject, content)
	services.SendNotification(paket.PpId, subject, content)
	for _, a := range paket.Pokja().AnggotaList() {
		services.SendNotification(a.ID, subject, content)
	}

	return flashSuccess(c, fmt.Sprintf("Addendum ke-%d Berhasil Diaktifkan. Silakan upload revisi dokumen.", paket.AddendumCount), "/dok-final/"+utils.UintToString(id))
}

func FinishAddendumReview(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)

	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}

	// Kembalikan status Addendum ke false, tapi JANGAN lock dulu
	// Lock akan otomatis terjadi setelah semua pihak menyetujui dokumen
	paket.IsAddendum = false
	paket.IsLockedReview = false
	models.SavePaket(&paket)

	// Send Notifications
	subject := "Penyelesaian Addendum Reviu Dokumen: " + paket.Nama
	content := "Proses addendum selesai. Silakan lakukan persetujuan ulang pada halaman Dokumen Final."
	
	services.SendNotification(paket.PpkId, subject, content)
	services.SendNotification(paket.PpId, subject, content)
	for _, a := range paket.Pokja().AnggotaList() {
		services.SendNotification(a.ID, subject, content)
	}

	return flashSuccess(c, "Addendum Berhasil Diselesaikan. Silakan lakukan persetujuan ulang.", "/dok-final/"+utils.UintToString(id))
}

// ResetAddendumCount - hanya bisa dilakukan oleh Pokja atau Admin
func ResetAddendumCount(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))

	// Hanya Pokja atau Admin yang bisa reset kuota addendum
	if !mp["isPokja"].(bool) && !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}

	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.SendStatus(404)
	}

	paket.AddendumCount = 0
	models.SavePaket(&paket)

	// Notifikasi ke PPK bahwa kuota addendum sudah direset
	subject := "Kuota Addendum Direset: " + paket.Nama
	content := fmt.Sprintf(
		"Kuota addendum untuk paket \"%s\" telah direset oleh Pokja. Anda kini dapat mengaktifkan addendum kembali (sisa: 3 kali).",
		paket.Nama,
	)
	services.SendNotification(paket.PpkId, subject, content)

	return flashSuccess(c, "Kuota addendum berhasil direset. PPK dapat mengaktifkan addendum kembali.", "/dok-final/"+utils.UintToString(id))
}

func DokPersiapanPaketPersetujuan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	dok_persiapan := services.GetDokPersiapan(id)
	if dok_persiapan.ID == 0 {
		return c.SendStatus(404)
	}
	// checklist := services.GetChecklistsBYjenis(paket.KgrId)
	mp["dokPersiapan"] = dok_persiapan
	mp["allowUpload"] = mp["isPPK"].(bool) || mp["isPokja"].(bool) || mp["isPP"].(bool)
	return c.Render("paket/dok-persiapan-persetujuan", mp)
}

func SimpanDokumenPersiapanPaketPersetujuan(c *fiber.Ctx) error {
	setuju, _ := strconv.ParseBool(c.FormValue("status", "false"));
	dokId, _ := strconv.Atoi(c.FormValue("id"))
	mp := currentMap(c)
	dokPersiapan := services.GetDokPersiapan(uint(dokId));
	if dokPersiapan.ID == 0 {
		log.Error("Dok final tidak ditemukan")
		if c.XHR() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Dok final tidak ditemukan"})
		}
		return flashError(c, "Simpan Dokumen Persiapan gagal", "/dok-final/" + c.Params("id"))
	}
	userid := utils.InterfaceToUint(mp["id"])
	err := dokPersiapan.SavePersetujuanPegawai(userid, setuju)
	if err != nil {
		log.Error(err)
		if c.XHR() {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return flashError(c, "Simpan Dokumen Persiapan gagal", "/dok-final/" + c.Params("id"))
	}

	// Auto-lock: cek apakah semua dokumen sudah disetujui semua pihak
	paket := services.GetPaket(dokPersiapan.PktId)
	if !paket.IsLockedReview {
		allApproved := true
		dokPersiapans := paket.DokPersiapan()
		if len(dokPersiapans) == 0 {
			allApproved = false
		}
		for _, dp := range dokPersiapans {
			if !dp.IsSudahPersetujuanSemua() {
				allApproved = false
				break
			}
		}
		if allApproved {
			paket.IsLockedReview = true
			models.SavePaket(&paket)
			log.Info("Paket ", paket.ID, " otomatis dikunci karena semua persetujuan selesai")
		}
	}

	if c.XHR() {
		return c.JSON(fiber.Map{"status": "success", "message": "Persetujuan berhasil disimpan"})
	}
	return flashSuccess(c, "Simpan Dokumen Persiapan Sukses","/dok-final/"+c.Params("id"))
}

func PengadaanPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	mp["paket"] = paket
	return c.Render("paket/pengadaan", mp)
}

func SimpanKodeTender(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	kode := utils.StringToUint(c.FormValue("kode"))
	err := services.SimpanKodeTender(id, kode)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Kode Tender Gagal", "/pengadaan/" + utils.UintToString(id))
	}
	return flashSuccess(c, "Simpan Kode Tender Berhasil", "/pengadaan/"+utils.UintToString(id))
}


func HasilPengadanPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	mp["paket"] = paket

	// Try to get from Tender Selesai
	tender := paket.GetTender()
	mp["tender"] = tender
	realisasiTender := tender.GetRealisasi()
	if len(realisasiTender) > 0 {
		mp["realisasi"] = realisasiTender[0]
	} else {
		// Try to get from Non-Tender Selesai
		nontender := paket.GetNontender()
		mp["nontender"] = nontender
		realisasiNontender := nontender.GetRealisasi()
		if len(realisasiNontender) > 0 {
			mp["realisasi"] = realisasiNontender[0]
		}
	}

	// Also fetch purchase data just in case
	purchase := paket.GetPurchase()
	if purchase.OrderId != "" {
		mp["purchase"] = purchase
	}

	return c.Render("paket/hasil-pengadaan", mp)
}

func SimpanHasilPengadanPaket(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	status := utils.StringToUint(c.FormValue("status"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	paket.Status = int(status)
   err := services.SavePaket(paket)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Hasil tender Gagal", "/hasil/"+ utils.UintToString(id))
	}
	return flashSuccess(c, "Simpan Hasil tender Berhasil","/hasil/"+utils.UintToString(id))
}

func SimpanDokHasilPengadaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanDokHasilPengadaan(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Hasil pengadaan Gagal", "/hasil/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan  Hasil pengadaan Sukses", "/hasil/"+utils.UintToString(id))
}

func SimpanDokPendukungPengadaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	err := services.SimpanDokPendukungPengadaan(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Dokumen Pendukung pengadaan Gagal", "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan  Dokumen Pendukung pengadaan Sukses", "/paket/"+utils.UintToString(id))
}

func HapusDokPaket(c *fiber.Ctx) error {
	log.Info("hapus dok paket")
	id := utils.StringToUint(c.Params("id"))
	pktId, err := services.HapusDokPaket(id)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Dokumen Gagal", "/paket/" + utils.UintToString(pktId))
	}
	
	// If the deletion was triggered from the edit page, stay there
	referer := c.Get("Referer")
	if strings.Contains(referer, "/edit-manual/") {
		return flashSuccess(c, "Hapus Dokumen Sukses", referer)
	}
	
	return flashSuccess(c, "Hapus Dokumen Sukses", "/paket/"+utils.UintToString(pktId))
}

func HapusDokPersiapan(c *fiber.Ctx) error {
	log.Info("hapus dok persiapan")
	id := utils.StringToUint(c.Params("id"))
	pktId, err := services.HapusDokPersiapan(id)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Dokumen Gagal", "/dok-final/"+utils.UintToString(pktId))
	}

	return flashSuccess(c, "Hapus Dokumen Sukses", "/dok-final/"+utils.UintToString(pktId))
}

func GantiPPK(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	id := utils.StringToUint(c.Params("id"))
	ppkId := utils.StringToUint(c.FormValue("ppk_id"))
	err := services.AssignPaketPpk(id, ppkId, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal Ganti PPK", "/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Berhasil Ganti PPK","/paket/"+utils.UintToString(id))
}

func SavePPk(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}

	if paket.IsOnlyPpk() {
		paket.Status = 2
		paket.TglDisetujui = models.Datetime(time.Now())
	} else if paket.IsPaketPp() {
		if paket.PpId == 0 {
			return flashError(c, "Gagal: Silakan Pilih Pejabat Pengadaan terlebih dahulu sebelum menyimpan", "/paket/" + utils.UintToString(id))
		}
		paket.Status = 2
		paket.TglDisetujui = models.Datetime(time.Now())
	}

	err := services.SavePaket(paket)
    if err != nil {
		log.Error(err)
		return flashError(c, "Simpan Paket Gagal","/paket/" + utils.UintToString(id))
	}
	// Return the updated user
	return flashSuccess(c, "Simpan Paket Berhasil", "/paket/"+utils.UintToString(id))
}

func DownloadPendukung(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id") // id paket
	paket := services.GetPaket(uint(id))
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	log.Info("create zip file");
	var files []string
	for _,v := range paket.DokPendukungList() {
		files = append(files, v.Document().Filepath)
	}
	zipFile,  err := utils.CreateZip(files, "dok-pendukung.zip");
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}

func DownloadHasilPengadaan(c *fiber.Ctx) error {
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
	zipFile,  err := utils.CreateZip(files, "dok-hasil-pengadaan.zip");
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}

func SimpanDokTambahan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	userid := utils.InterfaceToUint(mp["id"])
	err := services.SimpanDokTambahan(c, id, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), c.Get("Referer"))
	}
	return flashSuccess(c, "Simpan Dokumen Tambahan Berhasil", c.Get("Referer"))
}

func DownloadDokTambahan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id") // id paket
	paket := services.GetPaket(uint(id))
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	log.Info("create zip file")
	var files []string
	for _, v := range paket.DokTambahanList() {
		files = append(files, v.Document().Filepath)
	}
	zipFile, err := utils.CreateZip(files, "dok-tambahan.zip")
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}

func GetMetodeFilter(c *fiber.Ctx) error {
	return c.JSON(models.GetActiveMetodePaket())
}
