package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type RupSwakelolaAnggaran struct {
	AsalDana			string 		`json:"asal_dana"`
	AsalDanaKlpd		string		`json:"asal_dana_klpd"`
	AsalDanaSatker		string 		`json:"asal_dana_satker"`
	JenisKlpd			string 		`json:"jenis_klpd"`
	KdKegiatan			uint 		`json:"kd_kegiatan"`
	KdKlpd				string 		`json:"kd_klpd"`
	KdKomponen			string 		`json:"kd_komponen"`
	KdRup				uint 		`json:"kd_rup"`
	KdRupLokal			uint 		`json:"kd_rup_lokal"`
	KdSatker			uint 		`json:"kd_satker"`
	KdSatkerStr			string 		`json:"kd_satker_str"`
	KdSubkegiatan		uint 		`json:"kd_subkegiatan"`
	Mak					string 		`json:"mak"`
	NamaKlpd			string 		`json:"nama_klpd"`
	NamaSatker			string 		`json:"nama_satker"`
	Pagu				float64		`json:"pagu"`
	StatusAktiRup		bool		`json:"status_aktif_rup"`
	StatusDeleteRup		bool		`json:"status_delete_rup"`
	StatusUmumkanRup	string 		`json:"status_umumkan_rup"`
	SumberDana			string 		`json:"sumber_dana"`
	TahunAnggaran		int 		`json:"tahun_anggaran"`
	TahunAnggaranDana	int 		`json:"tahun_anggaran_dana"`
}

func (c RupSwakelolaAnggaran) TableName() string {
	return "rup_swakelola_anggaran"
}

func SaveRupAnggaranSwakelola(datas *[]RupSwakelolaAnggaran, tahun int)  {
	if len(*datas) > 0 {
		// slog.Info("saving data RupAnggaranSwakelola", "size", len(*datas), "tahun", tahun)
		DB.Where("tahun_anggaran = ?", tahun).Delete(&RupSwakelolaAnggaran{})
		err := DB.CreateInBatches(datas, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
		*datas = nil
	}
}


type RupSwakelola struct {
 	JenisKlpd					string 			`json:"jenis_klpd"`
    KdKlpd						string 			`json:"kd_klpd"`
    KdKlpdPenyelenggara			string 			`json:"kd_klpd_penyelenggara"`
    KdRup						uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_rup"`
    KdRupLokal					uint 			`json:"kd_rup_lokal"`
    KdSatker					uint 			`json:"kd_satker"`
    KdSatkerStr					string 			`json:"kd_satker_str"`
    NamaKlpd					string 			`json:"nama_klpd"`
    NamaKlpdPenyelenggara		string 			`json:"nama_klpd_penyelenggara"`
    NamaPaket					string 			`json:"nama_paket"`
    NamaPpk						string 			`json:"nama_ppk"`
    NamaSatker					string 			`json:"nama_satker"`
    NamaSatkerPenyelenggara		string 			`json:"nama_satker_penyelenggara"`
    NipPpk						string 			`json:"nip_ppk"`
    Pagu						float64 		`json:"pagu"`
    StatusAktifRup				bool			`json:"status_aktif_rup"`
    StatusDeleteRup				bool			`json:"status_delete_rup"`
    StatusUmumkanRup			string 			`json:"status_umumkan_rup"`
    TahunAnggaran				int 			`json:"tahun_anggaran"`
    TglAkhirPelaksanaanKontrak	Date 			`json:"tgl_akhir_pelaksanaan_kontrak"`
    TglAwalPelaksanaanKontrak	Date			`json:"tgl_awal_pelaksanaan_kontrak"`
    TglBuatPaket				time.Time		`json:"tgl_buat_paket"`
    TglPengumumanPaket			time.Time		`json:"tgl_pengumuman_paket"`
    TipeSwakelola				int 			`json:"tipe_swakelola"`
    UraianPekerjaan				string 			`json:"uraian_pekerjaan"`
    UsernamePpk					string 			`json:"username_ppk"`
    VolumePekerjaan				string 			`json:"volume_pekerjaan"`
}

func (c RupSwakelola) TableName() string {
	return "rup_swakelola"
}


func SaveRupSwakelola(datas *map[uint]RupSwakelola, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&RupSwakelola{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
