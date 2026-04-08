package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type Pencatatan struct {
 	AlasanPembatalan			string 			`json:"alasan_pembatalan"`
    BuktiPembayaran				string 			`json:"bukti_pembayaran"`
    InformasiLainnya			string 			`json:"informasi_lainnya"`
    JenisKlpd					string 			`json:"jenis_klpd"`
    KategoriPengadaan			string 			`json:"kategori_pengadaan"`
    KdKlpd						string 			`json:"kd_klpd"`
    KdLpse						int				`json:"kd_lpse"`
    KdNontenderPct				uint			`gorm:"primaryKey;autoIncrement:false" json:"kd_nontender_pct"`
    KdPktDce					uint 			`json:"kd_pkt_dce"`
    KdRup						string 			`json:"kd_rup"`
    KdSatker					string			`json:"kd_satker"`
    KdSatkerStr					string 			`json:"kd_satker_str"`
    MtdPemilihan				string 			`json:"mtd_pemilihan"`
    NamaKlpd					string 			`json:"nama_klpd"`
    NamaPaket					string 			`json:"nama_paket"`
    NamaPpk						string 			`json:"nama_ppk"`
    NamaSatker					string 			`json:"nama_satker"`
    NilaiPdnPct					float64 		`json:"nilai_pdn_pct"`
    NilaiUmkPct					float64			`json:"nilai_umk_pct"`
    NipPpk						string 			`json:"nip_ppk"`
    Pagu						float64			`json:"pagu"`
    StatusNontenderPct			string 			`json:"status_nontender_pct"`
    StatusNontenderPctKet		string 			`json:"status_nontender_pct_ket"`
    SumberDana					string 			`json:"sumber_dana"`
    TahunAnggaran				int 			`json:"tahun_anggaran"`
    TglBuatPaket				time.Time		`json:"tgl_buat_paket"`
    TglMulaiPaket				time.Time		`json:"tgl_mulai_paket"`
    TglSelesaiPaket				time.Time		`json:"tgl_selesai_paket"`
    TotalRealisasi				float64			`json:"total_realisasi"`
    UraianPekerjaan				string 			`json:"uraian_pekerjaan"`
}

func (c Pencatatan) TableName() string {
	return "pencatatan"
}

func SavePencatatan(datas *map[uint]Pencatatan, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Pencatatan{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type PencatatanRealisasi struct {
	DokRealisasi			string 		`json:"dok_realisasi"`
    JenisKlpd				string 		`json:"jenis_klpd"`
    JenisRealisasi			string 		`json:"jenis_realisasi"`
    KdKlpd					string 		`json:"kd_klpd"`
    KdLpse					int 		`json:"kd_lpse"`
    KdNontenderPct			uint		`gorm:"primaryKey;autoIncrement:false" json:"kd_nontender_pct"`
    KdPaketDce				uint 		`json:"kd_paket_dce"`
    KdRupPaket				string 		`json:"kd_rup_paket"`
    KdSatker				string 		`json:"kd_satker"`
    KdSatkerStr				string 		`json:"kd_satker_str"`
    KetRealisasi			string 		`json:"ket_realisasi"`
    NamaKlpd				string 		`json:"nama_klpd"`
    NamaLpse				string 		`json:"nama_lpse"`
    NamaPaket				string 		`json:"nama_paket"`
    NamaPenyedia			string 		`json:"nama_penyedia"`
    NamaPpk					string 		`json:"nama_ppk"`
    NamaSatker				string 		`json:"nama_satker"`
    NilaiRealisasi			float64		`json:"nilai_realisasi"`
    NipPpk					string 		`json:"nip_ppk"`
    NoRealisasi				string 		`json:"no_realisasi"`
    NpwpPenyedia			string 		`json:"npwp_penyedia"`
    Pagu					float64		`json:"pagu"`
    TahunAnggaran			int 		`json:"tahun_anggaran"`
    TglRealisasi			time.Time	`json:"tgl_realisasi"`
}

func (c PencatatanRealisasi) TableName() string {
	return "pencatatan_realisasi"
}

func SavePencatatanRealisasi(datas *map[uint]PencatatanRealisasi, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&PencatatanRealisasi{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
