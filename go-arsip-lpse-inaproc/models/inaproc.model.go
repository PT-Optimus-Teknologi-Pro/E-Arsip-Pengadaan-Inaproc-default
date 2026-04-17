package models

import (
	"time"
)

// type Satker struct {
// 	KdSatker		uint		`gorm:"primaryKey;autoIncrement:false" json:"kd_satker"`
// 	NamaSatker		string 		`gorm:"nama_satker" json:"nama_satker"`
// 	Alamat 			string 		`gorm:"alamat" json:"alamat"`
//     Fax				string 		`gorm:"fax" json:"fax"`
//     JenisKlpd		string 		`gorm:"jenis_klpd" json:"jenis_klpd"`
//     JenisSatker		string 	 	`gorm:"jenis_satker" json:"jenis_satker"`
//     KdKlpd			string 		`gorm:"kd_klpd" json:"kd_klpd"`
//     KdSatkerStr		string 		`gorm:"kd_satker_str" json:"kd_satker_str"`
//     KetSatker		string 		`gorm:"ket_satker" json:"ket_satker"`
//     KodeEseleon		string 		`gorm:"kode_eselon" json:"kode_eselon"`
//     Kodepos			string 		`gorm:"kodepos" json:"kodepos"`
//     NamaKlpd		string 		`gorm:"nama_klpd" json:"nama_klpd"`
//     StatusSatker	string 		`gorm:"status_satker" json:"status_satker"`
//     TahunAktif		string 		`gorm:"tahun_aktif" json:"tahun_aktif"`
//     Telepon			string 		`gorm:"telepon" json:"telepon"`
// }

// func (c Satker) TableName() string {
// 	return "satker"
// }

// func GetAllSatker(tahun int) []Satker {
// 	var satkers []Satker
// 	filter := "%"+strconv.Itoa(tahun)+"%"
// 	db.Where("tahun_aktif like ?", filter).Find(&satkers)
// 	return satkers
// }

// func GetSatker(id uint) Satker {
// 	var satker Satker
// 	db.First(&satker, id)
// 	return satker
// }

type Program struct {
	IsDeleted		bool 		`json:"is_deleted"`
    JenisKlpd		string 		`json:"jenis_klpd"`
    KdKlpd			string 		`json:"kd_klpd"`
    KdProgram		uint 		`gorm:"primaryKey;autoIncrement:false" json:"kd_program"`
    KdProgramLokal	uint		`json:"kd_program_lokal"`
    KdProgramStr	string 		`json:"kd_program_str"`
    KdSatker		uint 		`json:"kd_satker"`
    NamaKlpd		string 		`json:"nama_klpd"`
    NamaProgram		string 		`json:"nama_program"`
    PaguProgram		float64		`json:"pagu_program"`
    TahunAnggaran	int 		`json:"tahun_anggaran"`
}

func (c Program) TableName() string {
	return "program"
}

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

func (obj Tender) GetRealisasi() []TenderSelesai {
	var paketSelesai []TenderSelesai
	db.Find(&paketSelesai, "kd_tender=?", obj.KdTender)
	return paketSelesai
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

func (obj Nontender) GetRealisasi() []NontenderSelesai {
	var paketSelesai []NontenderSelesai
	db.Find(&paketSelesai, "kd_nontender=?", obj.KdNontender)
	return paketSelesai
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
/**
 * katalog v6
 */
type Katalog struct {
	CountProduct		int 				`json:"count_product"`
    FiscalYear			int 				`json:"fiscal_year"`
    FundingYear			string 				`json:"funding_source"`
    KodeKldi			string 				`json:"kode_klpd"`
    KodePenyedia		string 				`json:"kode_penyedia"`
    KodeSatker			string 				`json:"kode_satker"`
    Mak					string 				`json:"mak"`
    NamaSatker			string 				`json:"nama_satker"`
    OrderDate			time.Time			`json:"order_date"`
    OrderId				string 				`gorm:"primaryKey;autoIncrement:false" json:"order_id"`
    RekanId				uint 				`json:"rekan_id"`
    RupCode				string 				`json:"rup_code"`
    RupDesc				string 				`json:"rup_desc"`
    RupName				string 				`json:"rup_name"`
    ShipmentStatus		string 				`json:"shipment_status"`
    ShippingFee			float64				`json:"shipping_fee"`
    Status				string 				`json:"status"`
    Total				float64				`json:"total"`
    TotalQty			float32				`json:"total_qty"`
}

func (c Katalog) TableName() string {
	return "katalog"
}

func (obj Katalog) Penyedia() Penyedia {
	var res Penyedia
	db.First(&res, "kode_penyedia = ?", obj.KodePenyedia)
	return res
}

type Penyedia struct {
	AlamatPenyedia 				string 		`json:"alamat_penyedia"`
    BentukUsaha					string 		`json:"bentuk_usaha"`
    Email						string 		`json:"email"`
    JenisPerusahaan				string 		`json:"jenis_perusahaan"`
    KbliId						string		`json:"kbli_id"`
    KbliName					string 		`json:"kbli_name"`
    KodePenyedia				string 		`gorm:"primaryKey;autoIncrement:false" json:"kode_penyedia"`
    NamaPenyedia				string 		`json:"nama_penyedia"`
    Nib							string 		`json:"nib"`
    NpwpPenyedia				string		`json:"npwp_penyedia"`
    RekanId						uint 		`json:"rekan_id"`
    StatusAktif					string 		`json:"status_aktif"`
    StatusUmk					int 		`json:"status_umkk"`
    Telepon 					string 		`json:"telepon"`
}

func (c Penyedia) TableName() string {
	return "penyedia"
}

/**
 * katalog v5
 */
type KatalogArchive struct {
	AlamatSatker				string 			`json:"alamat_satker"`
    CatatanProduk				string 			`json:"catatan_produk"`
    Deskripsi					string 			`json:"deskripsi"`
    EmailUserPokja				string 			`json:"email_user_pokja"`
    HargaSatuan					float64			`json:"harga_satuan"`
    JabatanPpk					string 			`json:"jabatan_ppk"`
    JmlJenisProduk				int 			`json:"jml_jenis_produk"`
    KdKabupatenWilayahHarga		uint 			`json:"kd_kabupaten_wilayah_harga"`
    KdKlpd						string 			`json:"kd_klpd"`
    KdKomoditas					uint 			`json:"kd_komoditas"`
    KdPaket						uint 			`json:"kd_paket"`
    KdPaketProduk				uint 			`json:"kd_paket_produk"`
    KdPenyedia					uint 			`json:"kd_penyedia"`
    KdPenyediaDistributor		uint 			`json:"kd_penyedia_distributor"`
    KdProduk					uint			`json:"kd_produk"`
    KdProvinsiWilayahHarga		uint 			`json:"kd_provinsi_wilayah_harga"`
    KdRup						uint 			`json:"kd_rup"`
    KdUserPokja					uint 			`json:"kd_user_pokja"`
    KdUserPpk					uint 			`json:"kd_user_ppk"`
    KodeAnggaran				string 			`json:"kode_anggaran"`
    Kuantitas					float32			`json:"kuantitas"`
    NamaPaket					string 			`json:"nama_paket"`
    NamaSatker					string 			`json:"nama_satker"`
    NamaSumberDana				string 			`json:"nama_sumber_dana"`
    NoPaket						string 			`json:"no_paket"`
    NoTelpUserPokja				string 			`json:"no_telp_user_pokja"`
    NpwpSatker					string 			`json:"npwp_satker"`
    OngkosKirim					float64			`json:"ongkos_kirim"`
    PaketStatusStr				string			`json:"paket_status_str"`
    PpkNip						string 			`json:"Ppk_nip"`
    SatkerId					uint	 		`json:"satker_id"`
    StatusPaket 				string 			`json:"status_paket"`
    TahunAnggaran				int 			`json:"tahun_anggaran"`
    TanggalBuatPaket			Date			`json:"tanggal_buat_paket"`
    TanggalEditPaket			Date			`json:"tanggal_edit_paket"`
    TotalHarga					float64			`json:"total_harga"`
}

func (c KatalogArchive) TableName() string {
	return "katalog_archive"
}

type PenyediaArchive struct {
 	AlamatPenyedia			string 		`json:"alamat_penyedia"`
    EmailPenyedia			string 		`json:"email_penyedia"`
    Kbli2020Penyedia		string 		`json:"kbli2020_penyedia"`
    KdPenyedia				uint 		`gorm:"primaryKey;autoIncrement:false" json:"kd_penyedia"`
    KodePenyediaSikap		uint 		`json:"kode_penyedia_sikap"`
    NamaPenyedia			string 		`json:"nama_penyedia"`
    NoTelpPenyedia			string 		`json:"no_telp_penyedia"`
    Mpwp16					string 		`json:"npwp_16"`
    NpwpPenyedia			string 		`json:"npwp_penyedia"`
    PenyediaUkm				string 		`json:"penyedia_ukm"`
}

func (c PenyediaArchive) TableName() string {
	return "penyedia_archive"
}

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
