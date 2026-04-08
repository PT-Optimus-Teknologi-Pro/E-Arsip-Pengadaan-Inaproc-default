package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type Kontrak struct {
	AlamatSatker					string 			`json:"alamat_satker"`
    AlasanAdendum					string 			`json:"alasan_addendum"`
    AlasanNilaiKontrak10Persen 		string 			`json:"alasan_nilai_kontrak_10_persen"`
    AlamatPenetapanStatusKontrak	string 			`json:"alasan_penetapan_status_kontrak"`
    AlasanUbahNilaiKontrak			string 			`json:"alasan_ubah_nilai_kontrak"`
    AnggotaKso						string 			`json:"anggota_kso"`
    ApakahAdendum					string 			`json:"apakah_addendum"`
    BentukUsahaPenyedia				string 			`json:"bentuk_usaha_penyedia"`
    InformasiLainnya				string 			`json:"informasi_lainnya"`
    JabatanPpk						string 			`json:"jabatan_ppk"`
    JabatanWakilPenyedia			string 			`json:"jabatan_wakil_penyedia"`
    JenisKlpd						string 			`json:"jenis_klpd"`
    JenisKontrak					string 			`json:"jenis_kontrak"`
    KdKlpd							string 			`json:"kd_klpd"`
    KdLpse							int 			`json:"kd_lpse"`
    KdPenyedia						uint 			`json:"kd_penyedia"`
    KdSatker						string 			`json:"kd_satker"`
    KdSatkerStr						string 			`json:"kd_satker_str"`
    KdTender						uint 			`json:"kd_tender"`
    KotaKontrak						string 			`json:"kota_kontrak"`
    LingkupPekerjaan				string 			`json:"lingkup_pekerjaan"`
    NamaKlpd						string 			`json:"nama_klpd"`
    NamaPaket						string 			`json:"nama_paket"`
    NamaPemilikRekBank				string 			`json:"nama_pemilik_rek_bank"`
    NamaPenyedia					string 			`json:"nama_penyedia"`
    NamaPpk							string 			`json:"nama_ppk"`
    NamaRekBank						string 			`json:"nama_rek_bank"`
    NamaSatker						string 			`json:"nama_satker"`
    NilaiKontrak					float64			`json:"nilai_kontrak"`
    NilaiPdnKontrak					float64			`json:"nilai_pdn_kontrak"`
    NilaiUmkKontrak					float64			`json:"nilai_umk_kontrak"`
    NipPpk							string 			`json:"nip_ppk"`
    NoKontrak						string 			`gorm:"primaryKey;autoIncrement:false" json:"no_kontrak"`
    NoRekBank						string 			`json:"no_rek_bank"`
    NoSkPpk							string 			`json:"no_sk_ppk"`
    NoSppbj							string 			`json:"no_sppbj"`
    Npwp16Penyedia					string			`json:"npwp_16_penyedia"`
    NpwpPenyedia					string 			`json:"npwp_penyedia"`
    StatusKontrak					string 			`json:"status_kontrak"`
    TahunAnggaran					int 			`json:"tahun_anggaran"`
    TglKontrak						time.Time		`json:"tgl_kontrak"`
    TglKontrakAkhir					time.Time		`json:"tgl_kontrak_akhir"`
    TglKontrakAwal					time.Time		`json:"tgl_kontrak_awal"`
    TglPenetapanStatusKontrak		time.Time		`json:"tgl_penetapan_status_kontrak"`
    TipePenyedia					string 			`json:"tipe_penyedia"`
    VersiAdendum					int 			`json:"versi_addendum"`
    WakilSahPenyedia				string 			`json:"wakil_sah_penyedia"`
}

func (c Kontrak) TableName() string {
	return "kontrak"
}

func SaveKontrak(datas *map[string]Kontrak, tahun int)  {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Kontrak{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type KontrakNontender struct {
	AlamatSatker					string 			`json:"alamat_satker"`
    AlasanAdendum					string 			`json:"alasan_addendum"`
    AlasanNilaiKontrak10Persen 		string 			`json:"alasan_nilai_kontrak_10_persen"`
    AlamatPenetapanStatusKontrak	string 			`json:"alasan_penetapan_status_kontrak"`
    AlasanUbahNilaiKontrak			string 			`json:"alasan_ubah_nilai_kontrak"`
    AnggotaKso						string 			`json:"anggota_kso"`
    ApakahAdendum					string 			`json:"apakah_addendum"`
    BentukUsahaPenyedia				string 			`json:"bentuk_usaha_penyedia"`
    InformasiLainnya				string 			`json:"informasi_lainnya"`
    JabatanPpk						string 			`json:"jabatan_ppk"`
    JabatanWakilPenyedia			string 			`json:"jabatan_wakil_penyedia"`
    JenisKlpd						string 			`json:"jenis_klpd"`
    JenisKontrak					string 			`json:"jenis_kontrak"`
    KdKlpd							string 			`json:"kd_klpd"`
    KdLpse							int 			`json:"kd_lpse"`
    KdSatker						string 			`json:"kd_satker"`
    KdSatkerStr						string 			`json:"kd_satker_str"`
    KdNontender						uint 			`json:"kd_nontender"`
    KotaKontrak						string 			`json:"kota_kontrak"`
    LingkupPekerjaan				string 			`json:"lingkup_pekerjaan"`
    NamaKlpd						string 			`json:"nama_klpd"`
    NamaPaket						string 			`json:"nama_paket"`
    NamaPemilikRekBank				string 			`json:"nama_pemilik_rek_bank"`
    NamaPenyedia					string 			`json:"nama_penyedia"`
    NamaPpk							string 			`json:"nama_ppk"`
    NamaRekBank						string 			`json:"nama_rek_bank"`
    NamaSatker						string 			`json:"nama_satker"`
    NilaiKontrak					float64			`json:"nilai_kontrak"`
    NilaiPdnKontrak					float64			`json:"nilai_pdn_kontrak"`
    NilaiUmkKontrak					float64			`json:"nilai_umk_kontrak"`
    NipPpk							string 			`json:"nip_ppk"`
    NoKontrak						string 			`gorm:"primaryKey;autoIncrement:false" json:"no_kontrak"`
    NoRekBank						string 			`json:"no_rek_bank"`
    NoSkPpk							string 			`json:"no_sk_ppk"`
    NoSppbj							string 			`json:"no_sppbj"`
    Npwp16Penyedia					string			`json:"npwp_16_penyedia"`
    NpwpPenyedia					string 			`json:"npwp_penyedia"`
    StatusKontrak					string 			`json:"status_kontrak"`
    TahunAnggaran					int 			`json:"tahun_anggaran"`
    TglKontrak						time.Time		`json:"tgl_kontrak"`
    TglKontrakAkhir					time.Time		`json:"tgl_kontrak_akhir"`
    TglKontrakAwal					time.Time		`json:"tgl_kontrak_awal"`
    TglPenetapanStatusKontrak		time.Time		`json:"tgl_penetapan_status_kontrak"`
    TipePenyedia					string 			`json:"tipe_penyedia"`
    VersiAdendum					int 			`json:"versi_addendum"`
    WakilSahPenyedia				string 			`json:"wakil_sah_penyedia"`
    KontrakId						uint 			`json:"kontrak_id"`
    MtdPengadaan					string 			`json:"mtd_pengadaan"`
    SpkId							uint 			`json:"spk_id"`
    SppbjId							uint 			`json:"sppbj_id"`
}



func (c KontrakNontender) TableName() string {
	return "kontrak_nontender"
}

func SaveKontrakNontender(datas *map[string]KontrakNontender, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&KontrakNontender{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
