package models

import "log/slog"

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

func SaveKatalogArchive(value *[]KatalogArchive, tahun int)  {
	if len(*value) > 0 {
		DB.Where("tahun_anggaran = ?", tahun).Delete(&KatalogArchive{})
		err := DB.CreateInBatches(value, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
		*value = nil
	}
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


func SavePenyediaArchive(value *[]PenyediaArchive) error {
	return DB.Save(value).Error
}

func GetListPenyediaArchive() []string {
	var result []string
	DB.Raw("SELECT DISTINCT kd_penyedia FROM katalog_archive").Scan(&result)
	return result
}
