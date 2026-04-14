package models

import (
	"arsip/utils"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const (
	VERIFIKASI     	= 0
	APPROVED      	= 1
	APPROVED_UKPBJ 	= 2
	REJECT         	= 3
	ADMIN		   	= "ADMIN"
	UKPBJ		   	= "UKPBJ"
	PPK				= "PPK"
	ADM_AGENCY		= "ADM_AGENCY"
	PEGAWAI			= "PEGAWAI"
	POKJA			= "POKJA"
	PP				= "PP"
	ARSIPARIS		= "ARSIPARIS"
)

// Pegawai model
type Pegawai struct {
	gorm.Model
	PegNip            string `form:"peg_nip" json:"peg_nip"`
	CreatedBy         sql.NullInt64
	UpdatedBy         sql.NullInt64
	PegNama           string        `form:"peg_nama" json:"peg_nama"`
	PegAlamat         string        `form:"peg_alamat" json:"peg_alamat"`
	PegTelepon        string        `form:"peg_telepon" json:"peg_telepon"`
	PegEmail          string        `form:"peg_email" json:"peg_email"`
	PegMobile         string        `form:"peg_mobile" json:"peg_mobile"`
	PegGolongan       string        `form:"peg_golongan" json:"peg_golongan"`
	PegPangkat        string        `form:"peg_pangkat" json:"peg_pangkat"`
	PegJabatan        string        `form:"peg_jabatan" json:"peg_jabatan"`
	PegIsactive       int           `form:"peg_isactive" json:"peg_isactive"`
	PegNamauser       string        `form:"peg_namauser" json:"peg_namauser"`
	PegNoSk           string        `form:"peg_no_sk" json:"peg_no_sk"`
	PegMasaBerlaku    time.Time     `form:"peg_masa_berlaku" json:"peg_masa_berlaku"`
	AgcId             sql.NullInt64 `form:"agc_id"`
	PegNoPbj          string        `form:"peg_no_pbj" json:"peg_no_pbj"`
	Passw             string        `form:"passw"`
	ResetPassword     string        `form:"reset_password"`
	UkpbjId           sql.NullInt64 `form:"ukpbj_id"`
	PegNik            string        `form:"peg_nik" json:"peg_nik"`
	Usrgroup          string        `form:"usrgroup" json:"usrgroup"`
	PegTipeSertifikat int           `form:"peg_tipe_sertifikat" json:"peg_tipe_sertifikat"`
	LastChangePassw   time.Time     `form:"last_change_passw"`
	PegStatus         int           `form:"peg_status" json:"peg_status"`
	PegCatatan        string        `form:"peg_catatan" json:"peg_catatan"`
	TglApprove		  time.Time		`form:"tgl_approve" json:"tgl_approve"`
	TglReject		  time.Time		`form:"tgl_reject" json:"tgl_reject"`
}

func (Pegawai) TableName() string {
	return "pegawai"
}

func AnyAdmin() bool {
	var count int64
	db.Model(&Pegawai{}).Where("usrgroup=? and deleted_at IS NULL", ADMIN).Count(&count)
	return count > 0
}

func (u Pegawai) GetMasaBerlaku() string {
	return u.PegMasaBerlaku.Format("02-01-2006")
}

func (u Pegawai) GetStatusVerifikasi() string {
	switch u.PegStatus {
	case VERIFIKASI:
		return "Verifikasi"
	case APPROVED:
		return "Disetujui dan Akun SPSE Sudah dibuat"
	case APPROVED_UKPBJ:
		return "Disetujui"
	case REJECT:
		return "Ditolak"
	}
	return ""
}

func (u Pegawai) IsVerifikasi() bool {
	return u.PegStatus == VERIFIKASI
}

func (u Pegawai) IsReject() bool {
	return u.PegStatus == REJECT
}

func (u Pegawai) IsApprove() bool {
	return u.PegStatus == APPROVED || u.PegStatus == APPROVED_UKPBJ
}

func (u Pegawai) IsAktif() bool {
	return u.PegIsactive != 0;
}

func (u Pegawai) GetTglBuat() string {
	return u.CreatedAt.Format("02-01-2006 15:04")
}

func (u Pegawai) RoleLabel() string {
	switch u.Usrgroup {
	case ADMIN:
		return "Admin"
	case ADM_AGENCY:
		return "Admin Agency"
	case UKPBJ:
		return "Admin UKPBJ"
	case PPK:
		return "PPK"
	case POKJA:
		return "Pokja"
	case PP:
		return "Pejabat Pengadaan"
	case PEGAWAI:
		return "Pejabat Pengadaan / Pokja"
	case ARSIPARIS:
		return "Arsiparis"
	}
	return ""
}

func (u Pegawai) IsPPK() bool {
	return u.Usrgroup == PPK
}

func (u Pegawai) IsPokja() bool {
	return u.Usrgroup == POKJA
}

func (u Pegawai) IsPP() bool {
	return u.Usrgroup == PP
}

func (u Pegawai) IsArsiparis() bool {
	return u.Usrgroup == ARSIPARIS
}

func GetPegawai(id uint) Pegawai {
	var user Pegawai
	db.First(&user, id)
	return user
}

func GetPegawaiByUserPass(userid, hashedPassword string) Pegawai {
	var user Pegawai
	db.First(&user, "peg_namauser = ? AND passw = ?", userid, hashedPassword)
	return user
}

func GetCountPegawai() int64 {
	var result int64
	db.Model(&Pegawai{}).Count(&result)
	return result
}

func GetPPs(satkerid uint) []Pegawai {
	var res []Pegawai
	ppSatker := GetPejabatPengadaanSatker(satkerid)
	db.Find(&res, "usrgroup IN ('PP', 'PEGAWAI') AND peg_status IN (1, 2) AND id IN (SELECT peg_id FROM pejabat_pengadaan_pegawai WHERE pp_id=?)", ppSatker.PpId)
	return res
}

func GetPpks() []Pegawai {
	var res []Pegawai
	db.Find(&res, "usrgroup = ? AND peg_status IN (1, 2)", PPK)
	return res
}

func GetListAnggotaPokja() []Pegawai {
	var anggotas []Pegawai
	db.Find(&anggotas, "usrgroup IN ('POKJA', 'PEGAWAI') AND peg_isactive=1 AND peg_status IN (1,2) ")
	return anggotas
}

func GetListAnggotaPp() []Pegawai {
	var anggotas []Pegawai
	db.Find(&anggotas, "usrgroup IN ('PP', 'PEGAWAI') AND peg_isactive=1 AND peg_status IN (1,2) ")
	return anggotas
}

func GetPegawaiPengadaan() []Pegawai {
	var users []Pegawai
	db.Find(&users, "usrgroup IN ('POKJA','PP')")
	return users
}

func GetPegawaiNonPengadaan() []Pegawai {
	var users []Pegawai
	db.Find(&users, "usrgroup = 'PEGAWAI' ")
	return users
}

func GetAdminUKPBJ() []Pegawai {
	var users []Pegawai
	db.Find(&users, "usrgroup = 'UKPBJ' ")
	return users
}

func SavePegawai(user *Pegawai) error {
	return db.Save(user).Error
}

func DeletePegawai(user *Pegawai) error {
	return db.Delete(user).Error
}

type PerubahanData struct {
	gorm.Model
	Nomor		string		`form:"nomor" json:"nomor"`
	Perihal		string 		`form:"perihal" json:"perihal"`
	PegId		uint		`form:"peg_id" json:"peg_id"`
	DokId		uint 		`form:"dok_id" json:"dok_id"`
	Status		int 		`form:"status" json:"status"`
	Pegawai 	Pegawai		`gorm:"foreignKey:PegId"`
}

func (c PerubahanData) Dokumen() Document {
	var dokumen Document
	if c.DokId > 0 {
		db.First(&dokumen, c.DokId)
	}
	return dokumen
}

// func (c PerubahanData) Pegawai() Pegawai {
// 	var pegawai Pegawai
// 	if c.PegId > 0 {
// 		db.First(&pegawai, c.PegId)
// 	}
// 	return pegawai
// }

func (c PerubahanData) StatusLabel() string {
	switch c.Status {
	case 0:
		return "Pengajuan"
	case 1:
		return "Proses"
	case 2:
		return "Selesai"
	}
	return ""
}

func (c PerubahanData) IsProses() bool {
	return c.Status == 1
}

func (c PerubahanData) IsSelesai() bool {
	return c.Status == 2
}

func GetPerubahanData(id uint) PerubahanData {
	var res PerubahanData
	db.First(&res, id)
	return res
}

func SavePerubahanData(data *PerubahanData) error {
	return db.Save(data).Error
}

func DeletePerubahanData(data *PerubahanData) error {
	return db.Delete(data).Error
}

func AutoCreateAdmin()  {
	if !AnyAdmin() {
		log.Info("setup user admin with password 123456")
		pegawai := Pegawai{
			PegNip: "-",
			PegNama: "Admin",
			PegEmail: "admin@domain.com",
			PegIsactive: 1,
			PegNamauser: "ADMIN",
			Usrgroup: "ADMIN",
			PegStatus: 1,
			Passw: utils.HashPassword("123456"), // default password admin
		}
		db.Save(&pegawai)
	}
}
