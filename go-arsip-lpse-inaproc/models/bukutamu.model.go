package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type BukuTamu struct {
	gorm.Model
	Nama			string 		`json:"nama" form:"nama"`
	NamaPerusahaan	string 		`json:"nama_perusahaan"  form:"nama_perusahaan"`
	Kategori		string 		`json:"kategori" form:"kategori"`
	Keperluan		string 		`json:"keperluan" form:"keperluan"`
	KodeTender		uint 		`json:"kode_tender" form:"kode_tender"`
	Feedback		string		`json:"feedback" form:"feedback"`
	Layanan			string 		`json:"layanan" form:"layanan"`
	Status			int			`json:"status" form:"status"` // 0 blom terlayani, 1 proses, 2 sudah terlayani
	Email			string 		`json:"email" form:"email"`
	Phone			string 		`json:"phone" form:"phone"`
	DokId			uint 		`json:"dok_id" form:"dok_id"`
	PegId			uint 		`json:"peg_id" form:"peg_id"`
	TglProses		time.Time	`json:"tgl_proses" form:"tgl_proses"`
	TglSelesai		time.Time	`json:"tgl_selesai" form:"tgl_selesai"`
	Catatan 		string 		`json:"catatan" form:"catatan"`
	TindakLanjut	string 		`json:"tindak_lanjut" form:"tindak_lanjut"`
}

func (BukuTamu) TableName() string {
	return "buku_tamu"
}

func (obj BukuTamu) StatusLabel() string {
	if obj.Status == 1 {
		return "Diproses"
	}
	if obj.Status == 2 {
		return "Sudah Selesai"
	}
	return "Belum Diproses"
}

func (obj BukuTamu) Pegawai() Pegawai {
	var res Pegawai
	db.First(&res, obj.PegId)
	return res
}

var JENIS_LAYANAN = []string{"Helpdesk", "Pendaftaran", "Pembuktian Kualifikasi"}
var SKOR_LAYANAN = []string{"Kualitas Layanan", "Fasilitas Layanan", "Kelengkapan Informasi"}

type Feedback struct {
	gorm.Model
	Jenis 			int         `json:"jenis" form:"jenis"`      // jenis layanan
	Nama			string 		`json:"nama" form:"nama"`
	NamaPerusahaan	string 		`json:"nama_perusahaan"  form:"nama_perusahaan"`
	Kualitas		int 		`json:"kualitas" form:"kualitas"`   // skor kualitas
	Fasilitas		int			`json:"fasilitas" form:"fasilitas"` // skor fasilitas
	Kelengkapan		int 		`json:"kelengkapan" form:"kelengkapan"` // skor kelengkapan
}

func (Feedback) TableName() string {
	return "feedback"
}

func GetBukuTamu(id uint) BukuTamu {
	var buku BukuTamu
	db.First(&buku, id)
	return buku
}

func SaveBukuTamu(buku *BukuTamu) (uint, error) {
	err := db.Save(buku).Error
	if err != nil {
		return 0, err
	}
	return buku.ID, nil
}

func GetFeedback(id uint) Feedback {
	var feedback Feedback
	db.First(&feedback, id)
	return feedback
}

func SaveFeedback(kualitas, fasilitas, kelengkapan []string) error {
	var feedbacks []Feedback
	for i := range JENIS_LAYANAN {
		kualitasInt , _ := strconv.Atoi(kualitas[i])
		fasilitasInt , _ := strconv.Atoi(fasilitas[i])
		kelengkapanInt , _ := strconv.Atoi(kelengkapan[i])
		feedbacks = append(feedbacks, Feedback{
			Jenis: i+1,
			Nama: "Guest",
			Kualitas: kualitasInt,
			Fasilitas: fasilitasInt,
			Kelengkapan: kelengkapanInt,
		})
	}
	return  db.Save(&feedbacks).Error
}

func GetCountFeedbackByJenisKualitas(jenis int, start int) int64 {
	var result int64
	db.Model(&Feedback{}).Where("jenis = ? AND kualitas = ? and deleted_at IS NULL", jenis, start).Count(&result)
	return result
}

func GetCountFeedbackByJenisFasilitas(jenis int, start int) int64 {
	var result int64
	db.Model(&Feedback{}).Where("jenis = ? AND fasilitas = ? and deleted_at IS NULL", jenis, start).Count(&result)
	return result
}

func GetCountFeedbackByJenisKelengkapan(jenis int, start int) int64 {
	var result int64
	db.Model(&Feedback{}).Where("jenis = ? AND kelengkapan = ? and deleted_at IS NULL", jenis, start).Count(&result)
	return result
}

type SummaryFeedBack struct {
	Jenis 		int
	Kualitas	[5]int64
	Fasilitas	[5]int64
	Kelengkapan	[5]int64
}

func (obj SummaryFeedBack) Label() string {
	return JENIS_LAYANAN[obj.Jenis - 1]
}
