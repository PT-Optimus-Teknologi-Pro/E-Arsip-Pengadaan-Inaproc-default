package models

import (
	"arsip/utils"
	"database/sql"

	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)
		
const (
	BELUM_SETUJU     = 0
	SETUJU			 = 1
	TIDAK_SETUJU	 = 2
)

type Checklist struct {
	gorm.Model
	Jenis       	int    		`gorm:"not null" json:"jenis"`
	Metode			int 		`gorm:"not null,default:0" json:"metode"`
	PeriodeAwal  	time.Time 	`json:"periode_awal"` // Tanggal penggunaan surat misalkan 1 jan 2025 sd 31 des 2025
	PeriodeAkhir 	time.Time 	`json:"periode_akhir"`
}

func (s Checklist) TableName() string {
	return "checklist"
}

func (c Checklist) JenisLabel() string {
	return jenisPengadaan[c.Jenis]
}

func (c Checklist) MetodeLabel() string {
	return metodePengadaan[c.Metode]
}

func (c Checklist) ChecklistDok() []ChecklistDok {
	return GetChecklistDok(c.ID)
}

func Jenis(id int) string {
	return jenisPengadaan[id]
}

type ChecklistDok struct {
	gorm.Model
	ChkId		uint          `json:"chk_id" form:"chk_id"`
	DokId		uint 		  `json:"dok_id" form:"dok_id"`
	Status		int		 	  `json:"status" form:"status"`
}

func (c ChecklistDok) Template() DokTemplate {
	return GetDokTemplate(c.DokId)
}

func (c ChecklistDok) StatusLabel() string {
	switch c.Status {
	case 0:
		return "opsional"
	case 1:
		return "wajib"
	default:
		return "belom ditentukan"
	}
}

func (s ChecklistDok) TableName() string {
	return "checklist_dok"
}

func GetChecklistDok(checkId uint) []ChecklistDok {
	var res []ChecklistDok
	db.Find(&res, "chk_id=?", checkId)
	return res
}

type ChecklistPaket struct {
	gorm.Model
	PktId      	uint          `json:"pkt_id" form:"pkt_id"`
	ChkId      	uint          `json:"chk_id" form:"chk_id"`
	Jenis 	   	string 		  `json:"jenis" form:"jenis"`
	DokTemplate	uint 		  `json:"dok_template" form:"dok_template"`
	CreatedBy  	sql.NullInt64 `json:"created_by"`
	DokId      	uint          `json:"dok_id" form:"dok_id"`
	Status     	uint          `json:"status" form:"status"` // 0 : diajukan, 1: diapprove, 2: direvisi, 3:ditolak
	TglAjukan  	sql.NullTime  `json:"tgl_ajukan"`
	TglApprove 	sql.NullTime  `json:"tgl_approve"`
	TglRevisi  	sql.NullTime  `json:"tgl_revisi"`
	TglTolak   	sql.NullTime  `json:"tgl_tolak"`
}

func (c ChecklistPaket) Dokumen() Document {
	var dokumen Document
	if c.DokId > 0 {
		db.First(&dokumen, c.DokId)
	}
	return dokumen
}

func (c ChecklistPaket) Template() DokTemplate {
	var dokumen DokTemplate
	if c.DokTemplate > 0 {
		db.First(&dokumen, c.DokTemplate)
	}
	return dokumen
}

func (ChecklistPaket) TableName() string {
	return "checklist_paket"
}

type ChecklistPaketHistory struct {
	gorm.Model
	CreatedBy sql.NullInt64 `json:"created_by"`
	CheckId   uint          `json:"check_id" form:"check_id"` // refer to ID ChecklistPaket
	DokId     uint          `json:"dok_id" form:"dok_id"`
}

func (ChecklistPaketHistory) TableName() string {
	return "checklist_paket_history"
}

type DokPersiapan struct {
	gorm.Model
	PktId      uint          `json:"pkt_id" form:"pkt_id"`
	ChkId      uint          `json:"chk_id" form:"chk_id"`
	CreatedBy  sql.NullInt64 `json:"created_by"`
	DokId      uint          `json:"dok_id" form:"dok_id"`
}

func (DokPersiapan) TableName() string {
	return "dok_persiapan"
}

func (c DokPersiapan) Dokumen() Document {
	var dokumen Document
	if c.DokId > 0 {
		db.First(&dokumen, c.DokId)
	}
	return dokumen
}

func (c DokPersiapan) Persetujuan() []PersetujuanDokPersiapan {
	var res []PersetujuanDokPersiapan
	db.Find(&res, "dkp_id=?", c.ID)
	return res
}

func (c DokPersiapan) PersetujuanPegawai(pegId uint) PersetujuanDokPersiapan {
	var res PersetujuanDokPersiapan
	db.First(&res, "dkp_id=? AND peg_id=?", c.ID, pegId)
	return res
}

func (c DokPersiapan) SavePersetujuanPegawai(pegId uint, status bool) error {
	persetujuan := c.PersetujuanPegawai(pegId)
	if persetujuan.ID > 0 {
		persetujuan.Status = status
	} else {
		persetujuan = PersetujuanDokPersiapan{
			DkpId: c.ID,
			PegId: pegId,
			Status: status,
		}
	}
	return db.Save(&persetujuan).Error
}

func (c *DokPersiapan) CheckPersetujuanPegawai() {
	if c.ID == 0 {
		return
	}
	
	persetujuanList := c.Persetujuan()
	paket := GetPaket(c.PktId)
	panitia := paket.Pokja()
	ppk := paket.Ppk()
	pp := paket.Pp()
	
	// 1. Bersihkan yang sudah tidak relevan
	for _, v := range persetujuanList {
		peg := v.Pegawai()
		isInvalid := false
		if peg.IsPokja() && !panitia.IsAnggota(v.PegId) {
			isInvalid = true
		} else if peg.IsPPK() && v.PegId != ppk.ID {
			isInvalid = true
		} else if peg.IsPP() && v.PegId != pp.ID {
			isInvalid = true
		}
		
		if isInvalid {
			db.Unscoped().Delete(&v)
		}
	}

	// 2. Tambahkan yang belum ada
	if ppk.ID > 0 {
		c.EnsurePersetujuanExist(ppk.ID)
	}
	if pp.ID > 0 {
		c.EnsurePersetujuanExist(pp.ID)
	}
	anggota := panitia.AnggotaList()
	for _, a := range anggota {
		c.EnsurePersetujuanExist(a.ID)
	}
}

func (c *DokPersiapan) EnsurePersetujuanExist(pegId uint) {
	if c.ID == 0 || pegId == 0 {
		return
	}
	var count int64
	db.Model(&PersetujuanDokPersiapan{}).Where("dkp_id = ? AND peg_id = ?", c.ID, pegId).Count(&count)
	if count == 0 {
		persetujuan := PersetujuanDokPersiapan{
			DkpId:  c.ID,
			PegId:  pegId,
			Status: false,
		}
		db.Create(&persetujuan)
	}
}

func (c DokPersiapan) IsBelumAdaPersetujuan() bool {
	belumAdaSetuju := true
	persetujuanList := c.Persetujuan()
	for _, v := range persetujuanList {
		if v.Status {
			belumAdaSetuju = false
			break
		}
	}
	return belumAdaSetuju
}

func (c DokPersiapan) IsSudahPersetujuanSemua() bool {
	persetujuanList := c.Persetujuan()
	// Jika belum ada persetujuan sama sekali, anggap belum selesai
	if len(persetujuanList) == 0 {
		return false
	}
	for _, v := range persetujuanList {
		if !v.Status {
			return false
		}
	}
	return true
}

func GetDokPersiapan(id uint) DokPersiapan {
	var res DokPersiapan
	db.First(&res, id)
	return res
}

func SaveAllDokPersiapan(c *fiber.Ctx, id uint, userid uint) error {
	paket := GetPaket(id)
	checks := []DokPersiapan{}
	for _, obj := range paket.Checklist() {
		log.Info("save document ", obj.ID)
		dokId, err := SaveDocument(c, userid,  DOKFINAL, "checklist_"+utils.UintToString(obj.ID))
		if err != nil {
			continue
		}
		
		// Selalu buat record baru (Multi-file support)
		checklistpaket := DokPersiapan {
			DokId: dokId,
			PktId: paket.ID,
			ChkId: obj.ID,
		}
		db.Save(&checklistpaket)
		checks = append(checks, checklistpaket)
	}
	
	for i := range checks {
		// Reset/Ensure approvals exist for the (updated) record
		checks[i].CheckPersetujuanPegawai()
	}
	return nil
}

func GetChecklists() []Checklist {
	var res []Checklist
	db.Find(&res)
	return res
}

func GetChecklistsBYjenis(jenis int) []ChecklistDok {
	var res []ChecklistDok
	db.Find(&res, "chk_id IN (SELECT ID FROM checklist WHERE jenis=? AND deleted_at IS NULL)", jenis)
	return res
}

func GetChecklistsBYJenisMetode(jenis int, metode int) []ChecklistDok {
	var res []ChecklistDok
	db.Find(&res, "chk_id IN (SELECT ID FROM checklist WHERE jenis=? and metode=? AND deleted_at IS NULL)", jenis, metode)
	return res
}

func GetChecklist(id uint) Checklist {
	var res Checklist
	db.First(&res, id)
	return res
}

func SaveChecklist(checklist *Checklist) error {
	return db.Save(checklist).Error
}

func HapusChecklist(id uint) error {
	checklist := GetChecklist(id)
	if checklist.ID == 0 {
		return fmt.Errorf("Checklist %d tidak ditemukan", id)
	}
	return db.Delete(&checklist, id).Error
}

func SimpanChecklist(checklist []ChecklistDok) error {
	save := []ChecklistDok{}
	var count int64
	if len(checklist) > 0 {
		db.Delete(&ChecklistDok{}, "chk_id = ?", checklist[0].ChkId)
		for _,o := range checklist {
			count = 0
			db.Model(&ChecklistDok{}).Where("chk_id = ? and dok_id=? and deleted_at IS NULL", o.ChkId, o.DokId).Count(&count)
			if count == 0 {
				save = append(save, o);
			}
		}
	}
	if len(save) == 0 {
		return nil
	}
	return db.Create(&save).Error
}

func JenisPengadaan(c int) string {
	return jenisPengadaan[c]
}

type PersetujuanDokPersiapan struct {
	gorm.Model
	DkpId		uint		`form:"dkp_id" json:"dkp_id"`
	PegId 		uint 		`form:"peg_id" json:"peg_id"`
	Status 		bool 		`form:"status" json:"status"`
}

func (PersetujuanDokPersiapan) TableName() string {
	return "persetujuan_dok_persiapan"
}

func (c PersetujuanDokPersiapan) Pegawai() Pegawai {
	var res Pegawai
	db.First(&res, c.PegId)
	return res
}

func (c PersetujuanDokPersiapan) DokPersiapan() DokPersiapan {
	var dokumen DokPersiapan
	if c.DkpId > 0 {
		db.First(&dokumen, c.DkpId)
	}
	return dokumen
}

func DeleteAllPersetujuanDokPersiapan(dkpId uint) error {
	return db.Unscoped().Where("dkp_id = ?", dkpId).Delete(&PersetujuanDokPersiapan{}).Error
}

type ReviewAddendum struct {
	gorm.Model
	PktId     uint   `json:"pkt_id"`
	Version   int    `json:"version"`
	Reason    string `json:"reason"`
	CreatedBy uint   `json:"created_by"`
}

func (ReviewAddendum) TableName() string {
	return "review_addendum"
}

func (c ReviewAddendum) Snapshot() []ReviewAddendumSnapshot {
	var res []ReviewAddendumSnapshot
	db.Find(&res, "addendum_id=?", c.ID)
	return res
}

func (c ReviewAddendum) Pegawai() Pegawai {
	var res Pegawai
	db.First(&res, c.CreatedBy)
	return res
}

type ReviewAddendumSnapshot struct {
	gorm.Model
	AddendumId uint   `json:"addendum_id"`
	PktId      uint   `json:"pkt_id"`
	ChkId      uint   `json:"chk_id"`
	DokId      uint   `json:"dok_id"`
	Approvals  string `json:"approvals"` // Store as JSON string [{pegawai: "Name", status: true}, ...]
}

func (ReviewAddendumSnapshot) TableName() string {
	return "review_addendum_snapshot"
}

func (c ReviewAddendumSnapshot) Dokumen() Document {
	var res Document
	db.First(&res, c.DokId)
	return res
}

func (c ReviewAddendumSnapshot) Checklist() ChecklistPaket {
	var res ChecklistPaket
	db.First(&res, c.ChkId)
	return res
}

func GetReviewAddendumList(pktId uint) []ReviewAddendum {
	var res []ReviewAddendum
	db.Where("pkt_id = ?", pktId).Order("version DESC").Find(&res)
	return res
}

