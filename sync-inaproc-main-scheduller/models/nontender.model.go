package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type JadwalNontender struct {
	KdAkt			uint 		`json:"kd_akt"`
    KdKlpd			string 		`json:"kd_klpd"`
    KdLpse			int 		`json:"kd_lpse"`
    KdNontender		uint 		`gorm:"index:idx_jadwal_kdnontender" json:"kd_nontender"`
    KdSatker		string 		`json:"kd_satker"`
    KdSatkerStr		string 		`json:"kd_satker_str"`
    KdTahapan		uint 		`json:"kd_tahapan"`
    NamaAkt			string 		`json:"nama_akt"`
    NamaTahapan		string 		`json:"nama_tahapan"`
    TahunAnggaran	int			`json:"tahun_anggaran"`
    TglAkhir		time.Time	`json:"tgl_akhir"`
    TglAwal			time.Time 	`json:"tgl_awal"`
}


func (c JadwalNontender) TableName() string {
	return "jadwal_nontender"
}

func SaveJadwalNontender(values *[]JadwalNontender, tahun int) {
	if len(*values) > 0 {
		DB.Where("tahun_anggaran = ?", tahun).Delete(&JadwalNontender{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type Nontender struct {
	Hps 					float64 		`json:"hps"`
    JenisKlpd				string 			`json:"jenis_klpd"`
    JenisPengadaan			string 			`json:"jenis_pengadaan"`
    KdKlpd					string 			`json:"kd_klpd"`
    KdLpse					int 			`json:"kd_lpse"`
    KdNontender				uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_nontender"`
    KdPktDce				uint			`json:"kd_pkt_dce"`
    KdRup					string 			`json:"kd_rup"`
    KdSatker				string			`json:"kd_satker"`
    KdSatkerStr				string 			`json:"kd_satker_str"`
    KetDitutup				string 			`json:"ket_ditutup"`
    KetDiulang				string 			`json:"ket_diulang"`
    KontrakPembayaran		string 			`json:"kontrak_pembayaran"`
    KualifikasiPaket		string 			`json:"kualifikasi_paket"`
    LlsId					uint 			`json:"lls_id"`
    Mak						string 			`json:"mak"`
    MtdPemilihan			string 			`json:"mtd_pemilihan"`
    NamaKlpd				string 			`json:"nama_klpd"`
    NamaLpse				string 			`json:"nama_lpse"`
    NamaPaket				string 			`json:"nama_paket"`
    NamaSatker				string 			`json:"nama_satker"`
    NipNamaPokja			string 			`json:"nip_nama_pokja"`
    NipNamaPp				string 			`json:"nip_nama_pp"`
    NipNamaPpk				string 			`json:"nip_nama_ppk"`
    Pagu					float64			`json:"pagu"`
    RepeatOrder				string 			`json:"repeat_order"`
    StatusNontender			string 			`json:"status_nontender"`
    SumberDana				string 			`json:"sumber_dana"`
    TahunAnggaran			int 			`json:"tahun_anggaran"`
    TglBuatPaket			time.Time		`json:"tgl_buat_paket"`
    TglKolektifKolegial		time.Time		`json:"tgl_kolektif_kolegial"`
    TglPengumumanNontender	CustomeDate		`json:"tgl_pengumuman_nontender"`
    UrlLpse					string 			`json:"url_lpse"`
    VersiNontender			int				`json:"versi_nontender"`
}

func (c Nontender) TableName() string {
	return "nontender"
}

func SaveNontender(datas *map[uint]Nontender, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Nontender{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type NontenderSelesai struct {
	Hps						float64			`json:"hps"`
    JenisKlpd				string 			`json:"jenis_klpd"`
    JenisPengadaan			string 			`json:"jenis_pengadaan"`
    KdKlpd					string 			`json:"kd_klpd"`
    KdLpse					int 			`json:"kd_lpse"`
    KdNontender				uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_nontender"`
    KdPenyedia				uint 			`json:"kd_penyedia"`
    KdPktDce				uint			`json:"kd_pkt_dce"`
    KdRup					string 			`json:"kd_rup"`
    KdSatker				string			`json:"kd_satker"`
    KdSatkerStr				string 			`json:"kd_satker_str"`
    KontrakPembayaran		string 			`json:"kontrak_pembayaran"`
    KualifikasiPaket 		string 			`json:"kualifikasi_paket"`
    LlsId					uint 			`json:"lls_id"`
    LpseId					int 			`json:"lpse_id"`
    Mak						string 			`json:"mak"`
    MtdPemilihan			string 			`json:"mtd_pemilihan"`
    NamaKlpd				string 			`json:"nama_klpd"`
    NamaLpse				string 			`json:"nama_lpse"`
    NamaPaket				string 			`json:"nama_paket"`
    NamaPenyedia			string 			`json:"nama_penyedia"`
    NamaSatker				string 			`json:"nama_satker"`
    NilaiKontrak			float64			`json:"nilai_kontrak"`
    NilaiNegosiasi			float64			`json:"nilai_negosiasi"`
    NilaiPdnKontrak			float64			`json:"nilai_pdn_kontrak"`
    NilaiPenawaran			float64			`json:"nilai_penawaran"`
    NilaiTerkoreksi			float64			`json:"nilai_terkoreksi"`
    NilaiUmkKontrak			float64			`json:"nilai_umk_kontrak"`
    Npwp16Penyedia			string 			`json:"npwp16_penyedia"`
    NpwpPenyedia			string 			`json:"npwp_penyedia"`
    Pagu					float64			`json:"pagu"`
    StatusNontender			string 			`json:"status_nontender"`
    SumberDana				string 			`json:"sumber_dana"`
    TahunAnggaran			int 			`json:"tahun_anggaran"`
    TglPengumumanNontender 	time.Time 		`json:"tgl_pengumuman_nontender"`
    TglSelesaiNontender		time.Time		`json:"tgl_selesai_nontender"`
    UrlLpse					string 			`json:"url_lpse"`
}

func (c NontenderSelesai) TableName() string {
	return "nontender_selesai"
}

func SaveNontenderSelesai(datas *map[uint]NontenderSelesai, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&NontenderSelesai{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
