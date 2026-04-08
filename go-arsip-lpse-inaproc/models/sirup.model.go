package models

import (
	"arsip/utils"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/datatypes"
)


var sumberDanas = []string{"APBD", "APBN", "PHLN","PNBP","APBNP","APBDP", "BLU", "BLUD", "BUMN", "BUMD", "Lainnya"}
var jenisPengadaan = []string {"Belum Ditentukan", "Barang", "Pekerjaan Konstruksi", "Jasa Konsultansi", "Jasa Lainnya", "Terintegrasi"}
var metodePengadaan = []string{"Belum Ditentukan", "Lelang Umum", "Lelang Sederhana", "Lelang Terbatas", "Seleksi Umum",
	"Seleksi Sederhana","Pemilihan Langsung", "Penunjukan Langsung", "Pengadaan Langsung", "e-Purchasing","Sayembara",
	"Kontes","Lelang Cepat", "Tender", "Tender Cepat", "Seleksi", "Dikecualikan"}

type Provinsi struct {
	ID 		uint 	`gorm:"autoIncrement:false" json:"prp_id"`
	Nama 	string 	`json:"prp_nama"`
}

func (s Provinsi) TableName() string {
	return "provinsi"
}

func GetCountProvinsi() int64 {
	var count int64
	db.Model(&Provinsi{}).Count(&count)
	return count
}

func SaveAllPropinsi(propinsis *[]Provinsi) error {
	return db.Save(propinsis).Error
}

type Kabupaten struct {
	ID 		uint 	`gorm:"autoIncrement:false" json:"kbp_id"`
	PrpId	int 	`json:"prp_id"`
	Nama 	string  `json:"kbp_nama"`
}

func (s Kabupaten) TableName() string {
	return "kabupaten"
}

func SaveAllKabupaten(kabupatens *[]Kabupaten) error {
	return db.Save(kabupatens).Error
}


type SatkerSirup struct {
	ID				uint 		`gorm:"autoIncrement:false" json:"id"`
    IdSatker		string		`gorm:"id_satker" json:"id_satker"`
    IdKldi			string 		`gorm:"id_kldi" json:"id_kldi"`
    IsDeleted		bool		`gorm:"is_deleted" json:"is_deleted"`
    Nama			string 		`gorm:"nama" json:"nama"`
    Auditupdate		SirupTime	`gorm:"auditupdate" json:"auditupdate"` //Apr 10, 2025, 5:24:01 PM
    TahunAktif		string		`gorm:"tahun_aktif" json:"tahun_aktif"`
    Blu				bool		`gorm:"blu" json:"blu"`
    JenisSatkerId 	int			`gorm:"jenis_satker_id" json:"jenis_satker_id"`
}

func (s SatkerSirup) TableName() string {
	return "satker"
}

func GetSatkerSirup(id uint) SatkerSirup {
	var result SatkerSirup
	db.First(&result, id)
	return result
}

func GetAllSatkerSirup(tahun int) []SatkerSirup {
	var satkers []SatkerSirup
	filter := "%"+strconv.Itoa(tahun)+"%"
	db.Where("tahun_aktif like ?", filter).Find(&satkers)
	return satkers
}

func SaveSatkerSirup(satker *SatkerSirup) error {
	return db.Save(satker).Error
}

func SaveAllSatkerSirup(satkers *[]SatkerSirup) error {
	tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    for _, satker := range *satkers {
        if err := tx.Save(satker).Error; err != nil { // Save updates all fields if PK exists
            tx.Rollback()
            return err
        }
    }
    return tx.Commit().Error
}

func DeleteSatkerSirup(satker *SatkerSirup) error {
	return db.Delete(satker).Error
}

func GetCountSatker(tahun int) int64 {
	var count int64
	filter := "%"+strconv.Itoa(tahun)+"%"
	db.Model(&SatkerSirup{}).Where("tahun_aktif like ?", filter).Count(&count)
	return count
}

type APISatkerSirup struct {
 	ID   		uint
   	Nama 		string
    IdKldi		string
}

func (s APISatkerSirup) TableName() string {
	return "satker"
}

func GetSatkerAPI(tahun int) []APISatkerSirup {
	var satkers []APISatkerSirup
	filter := "%"+strconv.Itoa(tahun)+"%"
	db.Debug().Select([]string{"id", "nama", "id_kldi"}).Where("tahun_aktif like ?", filter).Find(&satkers)
	return satkers
}

type PaketSirup struct {
	ID						uint 			`gorm:"autoIncrement:false" json:"id"`
    Nama 					string 			`json:"nama"`
    PaketLokasiJson			datatypes.JSON	`json:"paket_lokasi_json"`
    Volume					string 			`json:"volume"`
    Keterangan 				string 			`json:"keterangan"`
    Spesifikasi 			string 			`json:"spesifikasi"`
    IsTkdn					bool 			`json:"is_tkdn"`
    IsPradipa				bool			`json:"is_pradipa"`
    PaketAnggaranJson		datatypes.JSON	`json:"paket_anggaran_json"`
    Pagu					float64			`json:"pagu"`
    PaketJenisJson			datatypes.JSON	`json:"paket_jenis_json"`
    MetodePengadaan			int 			`json:"metode_pengadaan"`
    TanggalKebutuhan		Date			`json:"tanggal_kebutuhan"`
    TanggalAwalPengadaan	Date			`json:"tanggal_awal_pengadaan"`
    TanggalAkhirPengadaan	Date			`json:"tanggal_akhir_pengadaan"`
    TanggalAwalPekerjaan	Date			`json:"tanggal_awal_pekerjaan"`
    TanggalAkhirPekerjaan	Date			`json:"tanggal_akhir_pekerjaan"`
    TanggalPengumuman		Datetime		`json:"tanggal_pengumuman"`
    IdSwakelola				uint			`json:"id_swakelola"`
    IdPpk					uint			`json:"id_ppk"`
    Umkm					bool 			`json:"umkm"`
    KodeKldi				string 			`json:"kode_kldi"`
    IdSatker				uint			`json:"id_satker"`
    EncryptedUsernamePpk	string 			`json:"encrypted_username_ppk"`
    PaketAktif				bool 			`json:"paket_aktif"`
    PaketTerhapus			bool			`gorm:"index:idx_paket_sirup_terhapus" json:"paket_terhapus"`
    StatusPaket				int				`json:"status_paket"`
    PaketTerumumkan			bool 			`json:"paket_terumumkan"`
    Tahun					int				`json:"tahun"`
    JenisPaket 				int 			`json:"jenis_paket"`
}

func (PaketSirup) TableName() string {
	return "paket_sirup"
}

func (obj PaketSirup) Satker() SatkerSirup {
	var res SatkerSirup
	db.First(&res, obj.IdSatker)
	return res
}

func (obj PaketSirup) Jenis() []PaketJenisSirup {
	var res []PaketJenisSirup
	if err := json.Unmarshal(obj.PaketJenisJson, &res); err != nil {
		log.Error(err)
	}
	return res
}

func (obj PaketSirup) Lokasi() []PaketLokasiSirup {
	var res []PaketLokasiSirup
	if err := json.Unmarshal(obj.PaketLokasiJson, &res); err != nil {
		log.Error(err)
	}
	return res
}

func (obj PaketSirup) Anggaran() []PaketAnggaranSirup {
	var res []PaketAnggaranSirup
	if err := json.Unmarshal(obj.PaketAnggaranJson, &res); err != nil {
		log.Error(err)
	}
	return res
}

func (obj PaketSirup) AnggaranLabel() string {
	anggarans := obj.Anggaran()
	var res string
	for _, o := range anggarans {
		label := o.SumberDanaLabel() + " " + utils.IntToString(o.TahunAnggaranDana)
		if res == label {
			continue
		}
		res += label
	}
	return res
}

func (obj PaketSirup) Metode() string {
	return metodePengadaan[obj.MetodePengadaan]
}

func (obj PaketSirup) JenisPengadaan() string {
	return jenisPengadaan[obj.JenisPaket]
}

func GetPaketSirup(id uint) PaketSirup {
	var paket PaketSirup
	db.First(&paket, id)
	return paket
}

func SavePaketSirupTransaction(pakets *[]PaketSirup) error {
	tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    for _, paket := range *pakets {
        if err := tx.Save(paket).Error; err != nil { // Save updates all fields if PK exists
            tx.Rollback()
            return err
        }
    }
    err := tx.Commit().Error
    if err == nil {
		err = db.Delete(&PaketSirup{}, "paket_terhapus=?", true).Error
	}
	return err
}

func GetCountPaketSirup() int64 {
	var count int64
	db.Model(&PaketSirup{}).Count(&count)
	return count
}

type SwakelolaSirup struct {
	ID 						uint 			`gorm:"autoIncrement:false" json:"id"`
    Nama 					string 			`json:"nama"`
    PaketLokasiJson			datatypes.JSON 	`json:"paket_lokasi_json"`
    LlsVolume				string 			`jsob:"lls_volume"`
    Keterangan				string 			`json:"keterangan"`
    IsPradiap				bool 			`json:"is_pradipa"`
    PaketAnggaranJson		datatypes.JSON	`json:"paket_anggaran_json"`
    JumlahPagu				float64			`json:"jumlah_pagu"`
    TanggalAwalPekerjaan	Date			`json:"tanggal_awal_pekerjaan"`
    TanggalAkhirPekerjaan	Date			`json:"tanggal_akhir_pekerjaan"`
    TanggalPengumuman		Datetime		`json:"tanggal_pengumuman"`
    IdPpk					uint			`json:"id_ppk"`
    KodeKldi				string 			`json:"kode_kldi"`
    IdSatker				uint			`json:"id_satker"`
    TipeSwakelola			int				`json:"tipe_swakelola"`
    SatkerLain				int 			`json:"satker_lain"`
    NamaSatkerLain			string			`json:"nama_satker_lain"`
    KldLain					string			`json:"kld_lain"`
    NamaKldLain				string			`json:"nama_kld_lain"`
    Aktif					bool			`json:"aktif"`
    Umumkan					bool 			`json:"umumkan"`
    IsDelete				bool			`json:"is_delete"`
    Status					int 			`json:"status"`
    Tahun					int				`json:"tahun"`
}

func (SwakelolaSirup) TableName() string {
	return "swakelola_sirup"
}

func GetSwakelolaSirup(id uint) SwakelolaSirup {
	var paket SwakelolaSirup
	db.First(&paket, id)
	return paket
}

// func SaveAllSwakelolaSirup(pakets *[]SwakelolaSirup) error {
// 	err := db.Save(pakets).Error
// 	if err == nil {
// 		err = db.Delete(&SwakelolaSirup{}, "is_delete=?", true).Error
// 	}
// 	return err
// }

func SaveSwakelolaSirupTransaction(pakets *[]SwakelolaSirup) error {
	tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    for _, paket := range *pakets {
        if err := tx.Save(paket).Error; err != nil { // Save updates all fields if PK exists
            tx.Rollback()
            return err
        }
    }
    err := tx.Commit().Error
    if err == nil {
		err = db.Delete(&SwakelolaSirup{}, "is_delete=?", true).Error
	}
	return err
}


func GetCountSwakelolaSirup() int64 {
	var count int64
	db.Model(&SwakelolaSirup{}).Count(&count)
	return count
}

type SirupTime time.Time

const sirupTimeLayout = "Jan 2, 2006, 15:04:05 PM" //Apr 10, 2025, 5:24:01 PM

func (ct *SirupTime) UnmarshalJSON(b []byte) (err error) {
  s := strings.Trim(string(b), "\"")
  nt, err := time.Parse(sirupTimeLayout, s)
  *ct = SirupTime(nt)
  return
}

func (date *SirupTime) Scan(value interface{}) (err error) {
	if value == nil {
		*date = SirupTime(time.Time{})
		return nil
	}
	*date = SirupTime(value.(time.Time))
	return nil
}

func (date SirupTime) Value() (driver.Value, error) {
	t := time.Time(date)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format("2006-01-02 15:04:05"), nil
}

type PaketLokasiSirup struct {
	ID 			uint 		`json:"id"`
	PktId 		uint 		`json:"pkt_id"`
	Auditupdate SirupTime	`json:"auditupdate"`
	IdProvinsi	int 		`json:"id_provinsi"`
	DetilLokasi	string		`json:"detil_lokasi"`
	IdKabupaten	uint 		`json:"id_kabupaten"`
}

func (obj PaketLokasiSirup) Propinsi() string {
	var propinsi Provinsi
	db.First(&propinsi, obj.IdProvinsi)
	return propinsi.Nama
}

func (obj PaketLokasiSirup) Kabupaten() string {
	var kabupaten Kabupaten
	db.First(&kabupaten, obj.IdKabupaten)
	return kabupaten.Nama
}

type PaketAnggaranSirup struct {
	ID 					uint 		`json:"id"`
    Mak 				string 		`json:"mak"`
    Pagu 				float64		`json:"pagu"`
    IdTp				uint 		`json:"id_tp"`
    Jenis				int 		`json:"jenis"`
    PktId				uint 		`json:"pkt_id"`
    AsalDana			string 		`json:"asal_dana"`
    DanaApbd			string 		`json:"dana_apbd"`
    AuditUpdate			SirupTime	`json:"auditupdate"`
    IdKegiatan			uint 		`json:"id_kegiatan"`
    IdKomponen			int 		`json:"id_komponen"`
    KodeSatker			uint 		`json:"kode_satker"`
    SumberDana			int 		`json:"sumber_dana"`
    IdDanaApbn			int 		`json:"id_dana_apbn"`
    KodeEseleon			int 		`json:"kode_esselon"`
    KodeInstansi		string 		`json:"kode_instansi"`
    PktIdClient			uint 		`json:"pkt_id_client"`
    AsalDanaSatker		uint		`json:"asal_dana_satker"`
    IdRinciObjectAkun	int 		`json:"id_rinci_objek_akun"`
    TahunAnggaranDana	int 		`json:"tahun_anggaran_dana"`
}

func (obj PaketAnggaranSirup) SumberDanaLabel() string {
	if obj.SumberDana == 0 {
		return sumberDanas[obj.SumberDana]
	}
	return sumberDanas[obj.SumberDana - 1]
}

type PaketJenisSirup struct {
	ID 			uint 	`json:"id"`
    Jenisid		int 	`json:"jenisid"`
    JumlahPagu  float64	`json:"jumlah_pagu"`
    PktId		uint 	`json:"pkt_id"`
}



func GetAllJenisPengadaan() []string {
	return jenisPengadaan
}

func GetAllMetodePengadaan() []string {
	return metodePengadaan
}

type StrukturAnggaran struct {
	ID 						uint 			`gorm:"primaryKey"`
    TahunAnggaran 			int				`json:"tahun_anggaran"`
    IdSatker 				uint 			`json:"id_satker"`
    IdKldi				 	string 			`json:"id_kldi"`
    BelanjaOperasi			float64			`json:"belanja_operasi"`
    BelanjaModal			float64			`json:"belanja_modal"`
    BelanjaBtt				float64			`json:"belanja_btt"`
    IdClient				uint			`json:"id_client"`
    BelanjaNonPengadaan 	float64			`json:"belanja_non_pengadaan"`
    BelanjaPengadaan		float64			`json:"belanja_pengadaan"`
    TotalBelanja			float64			`json:"total_belanja"`
    // Auditupdate				time.Time		`json:"Auditupdate"`
}

func (s StrukturAnggaran) TableName() string {
	return "struktur_anggaran"
}

func SaveAllStrukturAnggaran(data *[]StrukturAnggaran) error {
	tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    for i := range *data {
        if err := tx.Save(&(*data)[i]).Error; err != nil { // Use pointer to avoid copy
            tx.Rollback()
            return err
        }
    }
   return tx.Commit().Error
}
