package models

import (
	"log/slog"
	"time"

	"golang.org/x/exp/maps"
)


type RupAnggaran struct {
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

func (c RupAnggaran) TableName() string {
	return "rup_anggaran"
}

func SaveRupAnggaran(datas *[]RupAnggaran, tahun int) {
	if len(*datas) > 0 {
		// slog.Info("saving data RupAnggaran", "size", len(*datas), "tahun", tahun)
		DB.Where("tahun_anggaran = ?", tahun).Delete(&RupAnggaran{})
		err := DB.CreateInBatches(datas, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
		*datas = nil
	}
}


type Rup struct {
	  AlasanDikecualikan					string 			`json:"alasan_dikecualikan"`
      AlasanNonUkm							string 			`json:"alasan_non_ukm"`
      JenisKlpd								string 			`json:"jenis_klpd"`
      JenisPengadaan						string 			`json:"jenis_pengadaan"`
      KdJenisPengadaan						string 			`json:"kd_jenis_pengadaan"`
      KdKlpd								string 			`json:"kd_klpd"`
      KdMetodePengadaan						int				`json:"kd_metode_pengadaan"`
      KdRup									uint 			`gorm:"primaryKey;autoIncrement:false" json:"kd_rup"`
      KdRupLokal							uint 			`json:"kd_rup_lokal"`
      KdRupSwakelola						uint 			`json:"kd_rup_swakelola"`
      KdSatker								uint 			`json:"kd_satker"`
      KdSatkerStr							string 			`json:"kd_satker_str"`
      KodeRupTahunPertama					uint 			`json:"kode_rup_tahun_pertama"`
      MetodePengadaan						string 			`json:"metode_pengadaan"`
      NamaKlpd								string 			`json:"nama_klpd"`
      NamaPaket								string 			`json:"nama_paket"`
      NamaPpk								string 			`json:"nama_ppk"`
      NamaSatker							string 			`json:"nama_satker"`
      NipPpk								string 			`json:"nip_ppk"`
      NomorKontrak							string 			`json:"nomor_kontrak"`
      Pagu									float64			`json:"pagu"`
      SpesifikasiPekerjaan					string 			`json:"spesifikasi_pekerjaan"`
      SppAspekEkonomi						bool			`json:"spp_aspek_ekonomi"`
      SppAspekLingkungan					bool			`json:"spp_aspek_lingkungan"`
      SppAspekSosial						bool			`json:"spp_aspek_sosial"`
      StatusAktifRup						bool 			`json:"status_aktif_rup"`
      StatusDeleteRup						bool 			`json:"status_delete_rup"`
      StatusDikecualikan					bool 			`json:"status_dikecualikan"`
      StatusKonsolidasi						string 			`json:"status_konsolidasi"`
      StatusPdn								string 			`json:"status_pdn"`
      StatusPradipa							string 			`json:"status_pradipa"`
      StatusUkm								string 			`json:"status_ukm"`
      StatusUmumkanRup						string 			`json:"status_umumkan_rup"`
      TahunAnggaran							int 			`json:"tahun_anggaran"`
      TahunPertama							int 			`json:"tahun_pertama"`
      TglAkhirKontrak 						time.Time		`json:"tgl_akhir_kontrak"`
      TglAkhirPemanfaatan					Date 			`json:"tgl_akhir_pemanfaatan"`
      TglAkhirPemulihan						Date			`json:"tgl_akhir_pemilihan"`
      TglAwalKontrak						Date			`json:"tgl_awal_kontrak"`
      TglAwalPemanfaatan					Date			`json:"tgl_awal_pemanfaatan"`
      TglAwalPemilihan						Date			`json:"tgl_awal_pemilihan"`
      TglBuatPaket							time.Time    	`json:"tgl_buat_paket"`
      TglPengumumanPaket 					time.Time		`json:"tgl_pengumuman_paket"`
      TipePaket								string 			`json:"tipe_paket"`
      UraianPekerjaan						string 			`json:"urarian_pekerjaan"`
      UsernamePpk							string 			`json:"username_ppk"`
      VolumePekerjaan						string 			`json:"volume_pekerjaan"`
}

func (c Rup) TableName() string {
	return "rup"
}


func SaveRup(datas *map[uint]Rup, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("tahun_anggaran=?", tahun).Delete(&Rup{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
