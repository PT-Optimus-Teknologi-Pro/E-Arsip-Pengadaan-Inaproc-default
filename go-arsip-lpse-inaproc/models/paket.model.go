package models

import (
	"arsip/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var statusPaket = []string {"Draft", "Pengajuan", "Disetujui", "Tolak", "Kaji Ulang", "Proses", "selesai"}
var HPS_BATAS float64 = 2e8 // 200 jt
var HPS_BATAS_KONSTRUKSI float64 = 4e8 // 400jt

type Paket struct {
	gorm.Model
	Nama           string        `json:"nama" form:"nama"`
	CreatedBy      uint			 `json:"created_by" form:"created_by"`
	UpdatedBy      uint			 `json:"updated_by" form:"updated_by"`
	Pagu           float64       `json:"pagu" form:"pagu"`
	Hps            float64       `json:"hps" form:"hps"`
	UkpbjId        uint			 `json:"ukpbj_id" form:"ukpbj_id"`
	TglAssignUkpbj Datetime  	 `json:"tgl_assign_ukpbj" form:"tgl_assign_ukpbj"`
	PntId          uint			 `json:"pnt_id" form:"pnt_id"`
	TglAssignPokja Datetime  	 `json:"tgl_assign_pokja" form:"tgl_assign_pokja"`
	PpId           uint 		 `json:"pp_id" form:"pp_id"`
	TglAssignPp    Datetime		 `json:"tgl_assign_pp" form:"tgl_assign_pp"`
	Status         int			 `json:"status" form:"status"`
	KgrId          int           `json:"kgr_id" form:"kgr_id"`
	KodeTender     uint          `json:"kode_tender" form:"kode_tender"` // kode tender/nontender/swakelola
	PpkId		   uint			 `json:"ppk_id" form:"ppk_id"`
	RupId		   uint 		 `json:"rup_id" form:"rup_id"`
	SatkerId	   uint 		 `json:"satker_id" form:"satker_id"`
	Metode 		   int			 `json:"metode" form:"metode"`
	TglDisetujui   Datetime		 `json:"tgl_disetujui" form:"tgl_disetujui"`
	TglDitolak	   Datetime		 `json:"tgl_ditolak" form:"tgl_ditolak"`
	AlasanDitolak  string		 `json:"alasan_ditolak" form:"alasan_ditolak"`
	Prioritas 	   bool 		 `json:"prioritas" form:"prioritas"`
	TglGantiPpk   Datetime		 `json:"tgl_ganti_ppk" form:"tgl_ganti_ppk"`
}

func (Paket) TableName() string {
	return "paket"
}

func (obj Paket) IsDraft() bool {
	return obj.Status == 0
}

func (obj Paket) IsPengajuan() bool {
	return obj.Status == 1
}

func (obj Paket) IsDisetujui() bool {
	return obj.Status == 2
}

func (obj Paket) IsDitolak() bool {
	return obj.Status == 3
}

func (obj Paket) IsKajiUlang() bool {
	return obj.Status == 4
}

func (obj Paket) IsProses() bool {
	return obj.Status == 5
}

func (obj Paket) IsSelesai() bool {
	return obj.Status == 6
}

func (obj Paket) StatusLabel() string {
	return statusPaket[obj.Status]
}

func (obj Paket) Jenis() string {
	return jenisPengadaan[obj.KgrId]
}

func (obj Paket) MetodePengadaan() string {
	return metodePengadaan[obj.Metode]
}

func (obj Paket) Satker() SatkerSirup {
	var res SatkerSirup
	db.First(&res, obj.SatkerId)
	return res
}

func (obj Paket) Rup() PaketSirup {
	var res PaketSirup
	db.First(&res, obj.RupId)
	return res
}

func (obj Paket) Checklist() []ChecklistPaket {
	var res []ChecklistPaket
	db.Find(&res, "pkt_id=?", obj.ID)
	return res
}

func (obj Paket) CountChecklist() int64 {
	var count int64
	db.Model(&ChecklistPaket{}).Where("pkt_id=? and deleted_at IS NULL", obj.ID).Count(&count)
	return count
}

func (obj Paket) CountChecklistWithDok() int64 {
	var count int64
	db.Model(&ChecklistPaket{}).Where("pkt_id=? and dok_id > 0 and deleted_at IS NULL", obj.ID).Count(&count)
	return count
}

func (obj Paket) Ukpbj() Ukpbj {
	var res Ukpbj
	db.First(&res, obj.UkpbjId)
	return res
}

func (obj Paket) Pokja() Panitia {
	var res Panitia
	db.First(&res, obj.PntId)
	return res
}

func (obj Paket) Ppk() Pegawai {
	var res Pegawai
	db.First(&res, obj.PpkId)
	return res
}

func (obj Paket) Pp() Pegawai {
	var res Pegawai
	db.First(&res, obj.PpId)
	return res
}

func (obj Paket) DokPersiapan() []DokPersiapan {
	var res []DokPersiapan
	db.Find(&res, "pkt_id=?", obj.ID)
	return res
}

func (obj Paket) DokPendukungList() []DokPaket {
	var res []DokPaket
	db.Find(&res, "pkt_id=? AND jenis=?", obj.ID, PENDUKUNG)
	return res
}

func (obj Paket) DokHasilList() []DokPaket {
	var res []DokPaket
	db.Find(&res, "pkt_id=? AND jenis=?",obj.ID, HASIL_PENGADAAN)
	return res
}

func (obj Paket) DokKontrakList() []DokPaket {
	var res []DokPaket
	db.Find(&res, "pkt_id=? AND jenis=?",obj.ID, KONTRAK)
	return res
}

func (obj Paket) DokPekerjaanList() []DokPaket {
	var res []DokPaket
	db.Find(&res, "pkt_id=? AND jenis=?",obj.ID, HASIL_PEKERJAAN)
	return res
}

func (obj Paket) DokTambahanList() []DokPaket {
	var res []DokPaket
	db.Find(&res, "pkt_id=? AND jenis=?",obj.ID, TAMBAHAN)
	return res
}

func (obj Paket) GeneratePersyaratan() error {
	// generate checklist CreatePaket
	checklist := GetChecklistsBYJenisMetode(obj.KgrId, obj.Metode)
	if len(checklist) == 0 {
		return errors.New("Pembuatan paket Gagal.Admin Belum menentukan Persyaratan")
	}
	checks := []ChecklistPaket{}
	for _, o := range checklist {
		template := o.Template()
		checklistpaket := ChecklistPaket {
			PktId: obj.ID,
			Jenis: template.Jenis,
			DokTemplate: template.ID,
			ChkId: o.ID,
		}
		checks = append(checks, checklistpaket)
	}
	err := 	db.Save(&checks).Error
	if err != nil {
		log.Error(err)
		return errors.New("Pembuatan paket Gagal.")
	}
	return nil
}

func (obj Paket) SimpanPersyaratanPaket(c *fiber.Ctx, userid uint) error {
	checklist := obj.Checklist()
	checks := [] ChecklistPaket{}
	for _, o := range checklist {
		dokId, err := SaveDocument(c, userid,  CHECKLIST, "checklist_"+utils.UintToString(o.ID))
		if err != nil {
			continue
		}
		o.DokId = dokId
		checks = append(checks, o)
	}
	return db.Save(&checks).Error
}

func (obj Paket) IsAllowAjukan() bool {
	count := obj.CountChecklist()
	countDok := obj.CountChecklistWithDok()
	return obj.Status == 0 && count == countDok && count > 0 && obj.Hps > 0
}

func (obj Paket) IsOnlyPpk() bool {
	if obj.Metode == 7 || obj.Metode == 8 || obj.Metode == 9 { // PL, pengadaan langsung , e-purchasing
		if obj.KgrId == 2 {
			return obj.Hps > HPS_BATAS_KONSTRUKSI
		}
		return obj.Hps > HPS_BATAS
	}
	return false
}

func (obj Paket) IsPaketPokja() bool {
	return obj.Metode == 12 || obj.Metode == 13 || obj.Metode == 14 || obj.Metode == 15
}

func (obj Paket) IsPaketPp() bool {
	return obj.Metode == 8 || obj.Metode == 9
}

func (obj Paket) GetTender() Tender {
	var result Tender
	db.First(&result, obj.KodeTender)
	return result
}

func (obj Paket) GetNontender() Nontender {
	var result Nontender
	db.First(&result, obj.KodeTender)
	return result
}

func (obj Paket) GetPurchase() Katalog {
	var result Katalog
	db.First(&result, obj.KodeTender)
	return result
}

func (obj Paket) GetAllDocument(isPPK bool) []Document {
	var documents []Document
	for _, v := range obj.Checklist() {
		documents = append(documents, v.Dokumen())
	}
	for _, v := range obj.DokPersiapan() {
		documents = append(documents, v.Dokumen())
	}
	if isPPK {
		for _, v := range obj.DokPendukungList() {
			documents = append(documents, v.Document())
		}
	}
	for _, v := range obj.DokHasilList() {
		documents = append(documents, v.Document())
	}
	if isPPK {
		for _, v := range obj.DokKontrakList() {
			documents = append(documents, v.Document())
		}
		for _, v := range obj.DokPekerjaanList() {
			documents = append(documents, v.Document())
		}
	}
	return documents
}

func GetPaket(id uint) Paket {
	var paket Paket
	db.First(&paket, id)
	return paket
}

func CreatePaket(sirup PaketSirup, userId uint) (uint, error) {
	paket := Paket{
		Nama  : sirup.Nama,
	    CreatedBy : userId,
	    Pagu : sirup.Pagu,
	    Status : 0,
	    KgrId : sirup.JenisPaket,
		PpkId: userId,
		RupId: sirup.ID,
		SatkerId: sirup.IdSatker,
		Metode: sirup.MetodePengadaan,
	}
	err := db.Save(&paket).Error
	if err != nil {
		log.Error(err)
		return uint(0), errors.New("Pembuatan paket Gagal.")
	}
	err = paket.GeneratePersyaratan()
	if err != nil {
		log.Error(err)
		return uint(0), errors.New("Pembuatan paket Gagal.")
	}
	return paket.ID, nil
}

func SavePaket(paket *Paket) error {
	return db.Save(paket).Error
}

func DeletePaket(paket *Paket) error {
	return db.Delete(paket).Error
}


type PaketAnggaran struct {
	gorm.Model
	PktId uint `json:"pkt_id" form:"pkt_id"`
	AngId uint `json:"ang_id" form:"ang_id"`
	PpkId uint `json:"ppk_id" form:"ppk_id"`
	RupId uint `json:"rup_id" form:"rup_id"`
}

func (PaketAnggaran) TableName() string {
	return "paket_anggaran"
}

type PaketSatker struct {
	gorm.Model
	PktId uint `json:"pkt_id"`
	StkId uint `json:"stk_id"`
	RupId uint `json:"rup_id"`
}

func (PaketSatker) TableName() string {
	return "paket_satker"
}

type PaketLokasi struct {
	gorm.Model
	PKtId  uint   `json:"pkt_id"`
	KbpId  uint   `json:"kbp_id"`
	Lokasi string `json:"lokasi" form:"lokasi"`
}

func (PaketLokasi) TableName() string {
	return "paket_lokasi"
}


type DokPaket struct {
	gorm.Model
	PktId		uint 		`json:"pkt_id" form:"pkt_id"`
	DokId		uint 		`json:"dok_id" form:"dok_id"`
	PegId		uint 		`json:"peg_id" form:"peg_id"`
	Jenis		string 		`json:"jenis" form:"jenis"`
}

func (DokPaket) TableName() string {
	return "dok_paket"
}

func (obj DokPaket) Document() Document {
	var res Document
	db.First(&res, obj.DokId)
	return res
}

func GetDokPaket(id uint) DokPaket {
	var res DokPaket
	db.First(&res, id)
	return res
}

func SaveDokPaket(paket *DokPaket) error {
	return db.Save(paket).Error
}

func DeleteDokPaket(paket *DokPaket) error {
	return db.Delete(paket).Error
}

type PaketPPk struct {
	PaketId 	uint 		`gorm:"primaryKey"`
	PpkId		uint		`gorm:"primaryKey"`
	TglUpdate	time.Time	`json:"tgl_update"`
	PegId		uint 		`json:"peg_id"`
}

func (PaketPPk) TableName() string {
	return "paket_ppk"
}

func SavePaketPPk(paketPpk *PaketPPk) error {
	return db.Save(paketPpk).Error
}
