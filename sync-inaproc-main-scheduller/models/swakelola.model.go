package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type Swakelola struct {
	AlasanPembatalan			string 		`json:"alasan_pembatalan"`
    InformasiLainnya			string 		`json:"informasi_lainnya"`
    JenisKlpd					string 		`json:"jenis_klpd"`
    KdKlpd						string 		`json:"kd_klpd"`
    KdLpse						int 		`json:"kd_lpse"`
    KdPktDce					uint		`json:"kd_pkt_dce"`
    KdRup						string 		`json:"kd_rup"`
    KdSatker					string 		`json:"kd_satker"`
    KdSatkerStr					string 		`json:"kd_satker_str"`
    KdSwakelolaPct				uint		`gorm:"primaryKey;autoIncrement:false" json:"kd_swakelola_pct"`
    NamaKlpd					string 		`json:"nama_klpd"`
    NamaPaket					string 		`json:"nama_paket"`
    NamaPpk						string 		`json:"nama_ppk"`
    NamaSatker					string 		`json:"nama_satker"`
    NilaiPdnPct					float64		`json:"nilai_pdn_pct"`
    NilaiUmkPct					float64		`json:"nilai_umk_pct"`
    NipPpk						string 		`json:"nip_ppk"`
    Pagu						float64		`json:"pagu"`
    StatusSwakelolaPct			string 		`json:"status_swakelola_pct"`
    StatusSwakelolaPctKet		string 		`json:"status_swakelola_pct_ket"`
    SumberDana					string 		`json:"sumber_dana"`
    TahunAnggaran				int 		`json:"tahun_anggaran"`
    TglBuatPaket				time.Time	`json:"tgl_buat_paket"`
    TglMulaiPaket				time.Time	`json:"tgl_mulai_paket"`
    TglSelesaiPaket				time.Time	`json:"tgl_selesai_paket"`
    TipeSwakelola				int 		`json:"tipe_swakelola"`
    TipeSwakelolaNama			string 		`json:"tipe_swakelola_nama"`
    TotalRealisasi				float64		`json:"total_realisasi"`
    UraianPekerjaan				string 		`json:"uraian_pekerjaan"`
}

func (c Swakelola) TableName() string {
	return "swakelola"
}

func SaveSwakelola(datas *map[uint]Swakelola, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Swakelola{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type SwakelolaRealisasi struct {
	  DokRealisasi			string		`json:"dok_realisasi"`
      JenisRealisasi		string 		`json:"jenis_realisasi"`
      KdKlpd				string 		`json:"kd_klpd"`
      KdLpse				int 		`json:"kd_lpse"`
      KdSatker				string 		`json:"kd_satker"`
      KdSwakelolaPct		uint 		`gorm:"primaryKey;autoIncrement:false" json:"kd_swakelola_pct"`
      KetRealisasi			string 		`json:"ket_realisasi"`
      NamaPelaksana			string 		`json:"nama_pelaksana"`
      NamaPpk				string 		`json:"nama_ppk"`
      NilaiRealisasi		float64		`json:"nilai_realisasi"`
      NipPpk				string 		`json:"nip_ppk"`
      NoRealisasi			string 		`json:"no_realisasi"`
      NpwpPelaksana			string 		`json:"npwp_pelaksana"`
      RskId					uint		`json:"rsk_id"`
      TahunAnggaran			int 		`json:"tahun_anggaran"`
      TglRealisasi			time.Time	`json:"tgl_realisasi"`
}

func (c SwakelolaRealisasi) TableName() string {
	return "swakelola_realisasi"
}

func SaveSwakelolaRealisasi(datas *map[uint]SwakelolaRealisasi, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&SwakelolaRealisasi{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
