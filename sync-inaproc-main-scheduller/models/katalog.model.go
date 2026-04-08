package models

import (
	"log/slog"
	"time"
	"golang.org/x/exp/maps"
)

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

func SaveKatalog(datas *map[string]Katalog, tahun int) {
	if len(*datas) > 0 {
		values := maps.Values(*datas)
		DB.Unscoped().Where("fiscal_year=?", tahun).Delete(&Katalog{})
		err := DB.CreateInBatches(values, MAX_DATA).Error
		if err != nil {
			slog.Error(err.Error())
		}
	}
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

func SavePenyedia(value *[]Penyedia) error {
	return DB.Save(value).Error
}

func GetListPenyedia() []string {
	var result []string
	DB.Raw("SELECT DISTINCT kode_penyedia FROM katalog").Scan(&result)
	return result
}



