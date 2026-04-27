package services

import (
	"arsip/config"
	"arsip/models"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetPaket(id uint) models.Paket {
	return models.GetPaket(id)
}

func CreateManualPaket(c *fiber.Ctx, userId uint, ppkId uint, nama string, tahun int, pagu float64, hps float64, metode int, metodeStr string, jenis string, satkerId uint, keterangan string) (uint, error) {
	paket := models.Paket{
		Nama:        nama,
		Tahun:       tahun,
		Pagu:        pagu,
		Hps:         hps,
		Metode:      metode,
		MetodeArsip: metodeStr,
		JenisArsip:  jenis,
		SatkerId:    satkerId,
		Status:      0, // Draft
		Keterangan:  keterangan,
		CreatedBy:   userId,
		PpkId:       ppkId, // Auto-assigned when created by PPK
	}

	err := models.SavePaket(&paket)
	if err != nil {
		return 0, err
	}

	// Generate checklist
	paket.GeneratePersyaratan()

	// Handle Evidence File Upload if exists
	HandleManualEvidence(c, &paket, userId)

	return paket.ID, nil
}

func UpdateManualPaket(c *fiber.Ctx, id uint, userId uint, nama string, tahun int, pagu float64, hps float64, metode int, metodeStr string, jenis string, satkerId uint, keterangan string) error {
	paket := GetPaket(id)
	if paket.ID == 0 {
		return fmt.Errorf("Paket tidak ditemukan")
	}

	paket.Nama = nama
	paket.Tahun = tahun
	paket.Pagu = pagu
	paket.Hps = hps
	paket.Metode = metode
	paket.MetodeArsip = metodeStr
	paket.JenisArsip = jenis
	paket.SatkerId = satkerId
	paket.Keterangan = keterangan
	paket.UpdatedBy = userId

	err := models.UpdatePaket(&paket)
	if err != nil {
		return err
	}

	// Handle New Evidence File Upload if exists
	HandleManualEvidence(c, &paket, userId)

	return nil
}

func HandleManualEvidence(c *fiber.Ctx, paket *models.Paket, userId uint) {
	form, err := c.MultipartForm()
	if err != nil {
		return
	}

	files := form.File["bukti"]
	for _, file := range files {
		dokId, err := models.SaveDocumentHeader(c, userId, models.TAMBAHAN_PRIVATE, file)
		if err == nil {
			dokPaket := models.DokPaket{
				PktId: paket.ID,
				PegId: userId,
				DokId: dokId,
				Jenis: "Bukti Manual",
			}
			models.SaveDokPaket(&dokPaket)
		}
	}
}

func CreatePaket(rupId uint, userId uint) (uint, error) {
	paketSirup := GetPaketSirup(rupId)
	if paketSirup.ID == 0 {
		return uint(0), fmt.Errorf("Paket dengan kode %d tidak ditemukan", rupId)
	}
	return models.CreatePaket(paketSirup, userId)
}

func SimpanPaket(id uint, paket models.Paket) error {
	obj := models.GetPaket(id)
	if obj.ID == 0 {
		return fmt.Errorf("Paket dengan kode %d tidak ditemukan", id)
	}
	obj.Nama = paket.Nama
	obj.Tahun = paket.Tahun
	obj.Pagu = paket.Pagu
	obj.Hps = paket.Hps
	obj.Metode = paket.Metode
	obj.KgrId = paket.KgrId
	obj.SatkerId = paket.SatkerId
	obj.Status = paket.Status
	if paket.PpId != 0 {
		obj.PpId = paket.PpId
		obj.TglAssignPp = models.Datetime(time.Now())
	}
	if paket.PntId != 0 {
		obj.PntId = paket.PntId
		obj.TglAssignPokja = models.Datetime(time.Now())
	}
	if paket.UkpbjId != 0 {
		obj.UkpbjId = paket.UkpbjId
		obj.TglAssignUkpbj = models.Datetime(time.Now())
	}
	return models.SavePaket(&obj)
}

func SavePaket(paket models.Paket) error {
	return models.SavePaket(&paket)
}

func HapusPaket(id uint) error {
	obj := models.GetPaket(id)
	if obj.ID == 0 {
		return fmt.Errorf("Paket dengan kode %d tidak ditemukan", id)
	}
	return models.DeletePaket(&obj)
}

func SimpanHpsPaket(id uint, hps float64) error {
	paket := GetPaket(id)
	if paket.ID == 0 {
		return fmt.Errorf("Paket dengan kode %d tidak ditemukan", id)
	}
	if hps > paket.Pagu {
		return errors.New("HPS tidak boleh melebihi Pagu")
	}
	paket.Hps = hps
	return models.SavePaket(&paket)
}

func SimpanPersyaratanPaket(c *fiber.Ctx, id uint, userid uint) error {
	paket := GetPaket(id)
	return paket.SimpanPersyaratanPaket(c, userid)
}

func ApprovePaket(id uint, approve bool, reject bool, alasan string, prioritas bool) error {
	paket := GetPaket(id)
	if approve {
		paket.Prioritas = prioritas
		paket.Status = 2
		paket.TglDisetujui = models.Datetime(time.Now())
	} else if reject {
		paket.Status = 3
		paket.AlasanDitolak = alasan
		paket.TglDitolak = models.Datetime(time.Now())
	}
	return models.SavePaket(&paket)
}

// ajukan paket ke ukpbj
func KirimPaketUkpbj(id uint) error {
	paket := GetPaket(id)
	
	// Validasi HPS
	if paket.Hps <= 0 {
		return errors.New("HPS wajib diisi dan tidak boleh Rp0 sebelum kirim ke UKPBJ")
	}

	// Validasi Persyaratan Dokumen
	if !paket.IsPersyaratanLengkap() {
		return errors.New("Seluruh dokumen persyaratan (Checklist) wajib diunggah terlebih dahulu")
	}

	ukpbj := models.GetUkpbjAktif()
	paket.TglAssignUkpbj = models.Datetime(time.Now())
	paket.UkpbjId = ukpbj.ID
	paket.Status = 1
	return models.SavePaket(&paket)
}

func AssignPaketPokja(id uint, pntId uint) error {
	paket := GetPaket(id)
	// log.Info("paket id ", id, "pnt_id", pntId);
	if paket.Status < 2 {
		return fmt.Errorf("Paket Gagal assign ke Pokja dikarenakan status paket belom disetujui UKPBJ / masih draft")
	}
	panitia := GetPanitia(pntId)
	paket.TglAssignPokja = models.Datetime(time.Now())
	paket.PntId = panitia.ID
	return models.SavePaket(&paket)
}

func AssignPaketPp(id uint, pegId uint) error {
	paket := GetPaket(id)
	// if paket.Status < 2 {
	// 	return fmt.Errorf("Paket Gagal assign ke PP dikarenakan status paket belom disetujui UKPBJ / masih draft")
	// }
	pp := GetPegawai(pegId)
	paket.TglAssignPp = models.Datetime(time.Now())
	paket.PpId = pp.ID
	return models.SavePaket(&paket)
}

func AssignPaketPpk(id uint, pegId uint, userid uint) error {
	paket := GetPaket(id)
	if paket.Status < 2 {
		return fmt.Errorf("Paket Gagal assign ke PP dikarenakan status paket belom disetujui UKPBJ / masih draft")
	}
	pegawai := GetPegawai(pegId)
	// save to paket_ppk
	paketPpk := models.PaketPPk{
		PaketId: id,
		PpkId: pegawai.ID,
		TglUpdate: time.Now(),
		PegId: userid,
	}
	err := models.SavePaketPPk(&paketPpk)
	if err != nil {
		log.Error(err)
		return err
	}
	paket.TglGantiPpk = models.Datetime(time.Now())
	paket.PpkId = pegawai.ID
	return models.SavePaket(&paket)
}


func GetKajiUlangPaket(id uint) []models.KajiUlang {
	return models.GetKajiUlangPaket(id)
}

func GetKajiUlang(id uint) models.KajiUlang {
	return models.GetKajiUlang(id)
}

func SimpanKajiUlang(c *fiber.Ctx, userid uint, obj *models.KajiUlang) error {
	obj.DokId, _ = models.SaveDocument(c, userid, models.KAJIULANG, "dok_id")
	return models.SaveKajiUlang(obj)
}

func SimpanFotoRapatKajiUlang(c *fiber.Ctx, pktId uint, userid uint) error {
	dokId, err := models.SaveDocument(c, userid, models.FOTO_RAPAT, "foto_rapat")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket{
		PktId: pktId,
		PegId: userid,
		DokId: dokId,
		Jenis: models.FOTO_RAPAT,
	}
	return models.SaveDokPaket(&dokPaket)
}

func GetDokPersiapan(id uint) models.DokPersiapan {
	return models.GetDokPersiapan(id)
}

func SimpanDokPersiapanPaket(c *fiber.Ctx, id uint, userid uint) error {
	return models.SaveAllDokPersiapan(c, id, userid)
}

func SimpanPersetujuanDokPersiapan(dokId uint, pegId uint, status bool) error {
	dokPersiapan := GetDokPersiapan(dokId);
	if dokPersiapan.ID == 0 {
		return errors.New("Dok final tidak ditemukan")
	}
	err := dokPersiapan.SavePersetujuanPegawai(pegId, status)
	if err != nil {
		return err
	}
	return nil
}

func SimpanKodeTender(id uint, kode uint) error {
	paket := GetPaket(id)
	if paket.Status < 4 && !paket.IsOnlyPpk() {
		return fmt.Errorf("Paket Gagal assign ke Pokja dikarenakan status paket belum proses pengadaan")
	}
	paket.KodeTender = kode
	err := models.SavePaket(&paket)
	if err != nil {
		return err
	}

	// Trigger Sync from LPSE
	go func() {
		if paket.IsPaketPokja() {
			FetchTenderByCode(kode)
		} else if paket.IsPaketPp() {
			FetchNontenderByCode(kode)
		} else {
			// Fallback try both or just Tender
			FetchTenderByCode(kode)
			FetchNontenderByCode(kode)
		}
	}()

	return nil
}


func GetReviu(id uint) models.Reviu {
	return models.GetReviu(id)
}

func CreateReviu(reviu models.Reviu) error {
	return models.SaveReviu(&reviu)
}

func UpdateReviu(id uint, formreviu models.Reviu) error {
	reviu := GetReviu(uint(id))
	if reviu.ID == 0 {
		return errors.New("Reviu tidak ditemukan")
	}
	reviu.Content = formreviu.Content
	reviu.Bidang = formreviu.Bidang
	reviu.Opsi1 = formreviu.Opsi1
	reviu.Opsi2 = formreviu.Opsi2
	return models.SaveReviu(&reviu)
}

func DeleteReviu(id uint) error {
	reviu := GetReviu(uint(id))
	if reviu.ID == 0 {
		return errors.New("Reviu tidak ditemukan")
	}
	return models.DeleteReviu(&reviu)
}


func SimpanDokHasilPengadaan(c *fiber.Ctx, id uint, userid uint) error {
	//paket := GetPaket(id)
	dokId, err := models.SaveDocument(c, userid,  models.HASIL_PENGADAAN, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.HASIL_PENGADAAN,
	}
	return models.SaveDokPaket(&dokPaket)
}

func SimpanDokPendukungPengadaan(c *fiber.Ctx, id uint, userid uint) error {
	dokId, err := models.SaveDocument(c, userid,  models.PENDUKUNG, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.PENDUKUNG,
	}
	log.Info("simpan dok paket....")
	return models.SaveDokPaket(&dokPaket)
}


func SimpanDokKontrak(c *fiber.Ctx, id uint, userid uint) error {
	//paket := GetPaket(id)
	dokId, err := models.SaveDocument(c, userid,  models.KONTRAK, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.KONTRAK,
	}
	return models.SaveDokPaket(&dokPaket)
}


func SimpanDokHasilPekerjaan(c *fiber.Ctx, id uint, userid uint) error {
	dokId, err := models.SaveDocument(c, userid,  models.HASIL_PEKERJAAN, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.HASIL_PEKERJAAN,
	}
	return models.SaveDokPaket(&dokPaket)
}

func HapusDokPaket(id uint) (uint, error) {
	paket := models.GetDokPaket(id)
	if paket.ID == 0 {
		return 0, errors.New("File Tidak ditemukan")
	}
	DeleteDocument(paket.Document())
	err := models.DeleteDokPaket(&paket)
	return paket.PktId, err
}

func HapusDokPersiapan(id uint) (uint, error) {
	dok := models.GetDokPersiapan(id)
	if dok.ID == 0 {
		return 0, errors.New("Dokumen tidak ditemukan")
	}

	// 1. Delete associated approvals
	models.DeleteAllPersetujuanDokPersiapan(dok.ID)

	// 2. Delete actual document and record
	DeleteDocument(dok.Dokumen())

	// 3. Delete DokPersiapan record
	err := models.GetDB().Unscoped().Delete(&dok).Error

	return dok.PktId, err
}

func GetBeritaAcara(id uint) models.BeritaAcara {
	return models.GetBeritaAcara(id)
}

func AuthorisasiPaket(paket models.Paket, sessionMap fiber.Map) bool {
	if paket.ID == 0 {
		return false
	}
	id := sessionMap["id"].(uint)
	
	// Always allow the creator (especially important for Private Documents)
	if paket.CreatedBy == id {
		return true
	}

	// Always allow Admin, Ukpbj, and Arsiparis (Reviewers)
	if sessionMap["isAdmin"].(bool) || sessionMap["isUkpbj"].(bool) || sessionMap["isArsiparis"].(bool) {
		return true
	}

	if sessionMap["isPPK"].(bool) {
		return paket.Ppk().ID == id
	} else if sessionMap["isPokja"].(bool) {
		// Pokja bisa akses paket jika status sudah tahap pokja atau dia anggota panitia paket ini
		return paket.Status >= 2 || paket.Pokja().IsAnggota(id)
	} else if sessionMap["isPP"].(bool) {
		return paket.Pp().ID == id
	}
	return false
}

func GetTahunRupList() []int {
	var result []int
	start := config.TahunStart()
	now := time.Now().Year()
	for i := start; i <= now; i++ {
		result = append(result, i)
	}
	return result
}
func SimpanDokTambahan(c *fiber.Ctx, id uint, userid uint) error {
	dokId, err := models.SaveDocument(c, userid,  models.TAMBAHAN, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.TAMBAHAN,
	}
	return models.SaveDokPaket(&dokPaket)
}

func SimpanDokTambahanPrivate(c *fiber.Ctx, id uint, userid uint) error {
	dokId, err := models.SaveDocument(c, userid,  models.TAMBAHAN_PRIVATE, "file")
	if err != nil {
		return err
	}
	dokPaket := models.DokPaket {
		PktId: id,
		PegId: userid,
		DokId: dokId,
		Jenis: models.TAMBAHAN_PRIVATE,
	}
	return models.SaveDokPaket(&dokPaket)
}
