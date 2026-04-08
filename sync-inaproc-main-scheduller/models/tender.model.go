package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)

type Tender struct {
	Hps						float64 		`json:"hps"`
    JenisKlpd				string 			`json:"jenis_klpd"`
    JenisPengadaan			string 			`json:"jenis_pengadaan"`
    KdKlpd					string 			`json:"kd_klpd"`
    KdLpse					int				`json:"kd_lpse"`
    KdPktDce				uint 			`json:"kd_pkt_dce"`
    KdRup					string			`json:"kd_rup"`
    KdSatker				string 			`json:"kd_satker"`
    KdSatkerStr				string 			`json:"kd_satker_str"`
    KdTender				uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_tender"`
    KetDitutup				string 			`json:"ket_ditutup"`
    KetDiulang				string 			`json:"ket_diulang"`
    KgrId					int 			`json:"kgr_id"`
    KgrNama					string 			`json:"kgr_nama"`
    KontrakPembayaran		string 			`json:"kontrak_pembayaran"`
    KualifikasiPaket		string 			`json:"kualifikasi_paket"`
    ListTahunAnggaran		string 			`json:"list_tahun_anggaran"`
    LokasiPekerjaan			string 			`json:"lokasi_pekerjaan"`
    MtdEvaluasi 			string 			`json:"mtd_evaluasi"`
    MtdEValuasi				string 			`json:"mtd_kualifikasi"`
    MtdPemilihan			string 			`json:"mtd_pemilihan"`
    NamaKlpd				string 			`json:"nama_klpd"`
    NamaLpse				string 			`json:"nama_lpse"`
    NamaPaket				string 			`json:"nama_paket"`
    NamaPokja				string 			`json:"nama_pokja"`
    NamaPpk					string 			`json:"nama_ppk"`
    NamaSatker				string 			`json:"nama_satker"`
    NipPokja				string 			`json:"nip_pokja"`
    NipPpk					string 			`json:"nip_ppk"`
    Pagu 					float64			`json:"pagu"`
    StatusTender			string 			`json:"status_tender"`
    SumberDana				string 			`json:"sumber_dana"`
    TahunAnggaran			int 			`json:"tahun_anggaran"`
    TanggalStatus			time.Time		`json:"tanggal_status"`
    TglBuatPaket			time.Time		`json:"tgl_buat_paket"`
    TglKolektifKolegial		time.Time		`json:"tgl_kolektif_kolegial"`
    TglPengumumanTender		time.Time 		`json:"tgl_pengumuman_tender"`
    UrlLpse					string 			`json:"url_lpse"`
    VersiTender				int 			`json:"versi_tender"`
}

func (c Tender) TableName() string {
	return "tender"
}

func SaveTender(datas *map[uint]Tender, tahun int)  {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Tender{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type Jadwal struct {
	KdAkt				int			`json:"kd_akt"`
	KdKlpd				string 		`json:"kd_klpd"`
	KdLpse				int 		`json:"kd_lpse"`
    KdSatker			string 		`json:"kd_satker"`
    KdSatkerStr			string 		`json:"kd_satker_str"`
    KdTahapan			uint 		`json:"kd_tahapan"`
    KdTender			uint 		`gorm:"index:idx_jadwal_kdtender" json:"kd_tender"`
    NamaAkt				string 		`json:"nama_akt"`
    NamaTahapan			string 		`json:"nama_tahapan"`
    TahunAnggaran		int 		`json:"tahun_anggaran"`
    TglAkhir			time.Time	`json:"tgl_akhir"`
    TglAwal 			time.Time	`json:"tgl_awal"`
}

func (c Jadwal) TableName() string {
	return "jadwal"
}

func SaveJadwal(values *[]Jadwal, tahun int) {
	if len(*values) > 0 {
		DB.Where("tahun_anggaran = ?", tahun).Delete(&Jadwal{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}


type Peserta struct {
	  Alasan			string 			`json:"alasan"`
      KdKlpd			string 			`json:"kd_klpd"`
      KdLpse			int 			`json:"kd_lpse"`
      KdPenyedia		uint			`json:"kd_penyedia"`
      KdPeserta			uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_peserta"`
      KdPktDce			uint 			`json:"kd_pkt_dce"`
      KdSatker			string 			`json:"kd_satker"`
      KdSatkerStr		string 			`json:"kd_satker_str"`
      KdTender			uint 			`json:"kd_tender"`
      NamaPenyedia		string 			`json:"nama_penyedia"`
      NilaiPenawaran	float64			`json:"nilai_penawaran"`
      NilaiTerkoreksi	float64			`json:"nilai_terkoreksi"`
      NpwpPenyedia		string 			`json:"npwp_penyedia"`
      NpwpPenyedia16	string 			`json:"npwp_penyedia_16"`
      Pemenang			int				`json:"pemenang"`
      PemenangTerverifikasi		int 	`json:"pemenang_terverifikasi"`
      TahunAnggaran				int 	`json:"tahun_anggaran"`
}

func (c Peserta) TableName() string {
	return "peserta"
}

func SavePeserta(datas *map[uint]Peserta, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Peserta{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

type TenderSelesai struct {
	  Hps 				float64			`json:"hps"`
      JenisKlpd			string 			`json:"jenis_klpd"`
      KdKlpd			string 			`json:"kd_klpd"`
      KdLpse			int 			`json:"kd_lpse"`
      KdPaket			uint 			`json:"kd_paket"`
      KdPenyedia		uint 			`json:"kd_penyedia"`
      KdRupPaket		string 			`json:"kd_rup_paket"`
      KdSatker			string 			`json:"kd_satker"`
      KdTender			uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_tender"`
      NamaKlpd			string			`json:"nama_klpd"`
      NamaPenyedia		string 			`json:"nama_penyedia"`
      NamaSatker		string 			`json:"nama_satker"`
      NilaiKontrak		float64 		`json:"nilai_kontrak"`
      NilaiNegosiasi	float64			`json:"nilai_negosiasi"`
      NilaiPdnKontrak	float64			`json:"nilai_pdn_kontrak"`
      NilaiPenawaran	float64			`json:"nilai_penawaran"`
      NilaiTerkoreksi	float64			`json:"nilai_terkoreksi"`
      NilaiUmkKontrak	float64			`json:"nilai_umk_kontrak"`
      Npwp16Penyedia	string 			`json:"npwp_16_penyedia"`
      NpwpPenyedia		string 			`json:"npwp_penyedia"`
      Pagu				float64			`json:"pagu"`
      PsrId				uint			`json:"psr_id"`
      TahunAnggaran		int 			`json:"tahun_anggaran"`
      TglPenetapanPemenang 		time.Time	`json:"tgl_penetapan_pemenang"`
      TglPengumumanTender		time.Time 	`json:"tgl_pengumuman_tender"`
}

func (c TenderSelesai) TableName() string {
	return "tender_selesai"
}

func SaveTenderSelesai(datas *map[uint]TenderSelesai, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&TenderSelesai{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
