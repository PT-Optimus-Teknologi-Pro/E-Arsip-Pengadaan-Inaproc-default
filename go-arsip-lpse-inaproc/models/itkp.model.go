package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Itkp struct {
	// satker / SKPD
	KdSatker				uint 		`gorm:"primaryKey" json:"kd_satker"`
	// kode satker str
	KdSatkerStr				string 		`json:"kd_satker_str"`
	// Tahun Anggaran
	Tahun					int 		`gorm:"primaryKey" json:"tahun"`
	// ITKP SDM
	ItkpSdm					float32		`json:"itkp_sdm"`
	// total pagu RUP
	PaguRup					float64		`json:"pagu_rup"`
	// Belanja Pengadaan
	Belanja					float64		`json:"belanja"`
	// persentase ITKP RUP
	PersenRup				float32		`gorm:"type:numeric(10,2)" json:"persen_rup"`
	// ITKP RUP
	ItkpRup					float32		`gorm:"type:numeric(10,2)" json:"itkp_rup"`
	// jumlah paket transaksional
	PaketPemilihan			int			`json:"paket_pemilihan"`
	// jumlah ekontrak selesai
	PaketPemilihanSelesai	int 		`json:"paket_pemilihan_selesai"`
	// Persentasi Ekontrak
	PersenPemilihan			float32		`gorm:"type:numeric(10,2)" json:"persen_pemilihan"`
	// ITKP Ekontrak
	ItkpPemilihan			float32		`gorm:"type:numeric(10,2)" json:"itkp_pemilihan"`
	// jumlah paket tender
	PaketTender				int			`json:"paket_tender"`
	// jumlah paket tender selesai
	PaketTenderSelesai		int 		`json:"paket_tender_selesai"`
	// Pagu RUP tender
	PaguTender				float64		`json:"pagu_tender"`
	// Nilai ekontrak tender
	KontrakTender			float64 	`json:"kontrak_tender"`
	// persen itkp tender
	PersenTender			float32		`gorm:"type:numeric(10,2)" json:"persen_tender"`
	// Skor ITKP Tender
	ItkpTender				float32		`gorm:"type:numeric(10,2)" json:"itkp_tender"`
	// paket purchase aktif
	Purchase				int			`json:"purchase"`
	// paket purchase selesai
	PurchaseSelesai			int			`json:"purchase_selesai"`
	// Persen ITKP purchase
	PersenPurchase			float32		`gorm:"type:numeric(10,2)" json:"persen_purchase"`
	// ITKP Purchase
	ItkpPurchase			float32		`gorm:"type:numeric(10,2)" json:"itkp_purchase"`
	// Pagu nontender / non purchase
	PaguNontender			float64		`json:"pagu_nontender"`
	// Nilai ekontrak nontender / nonpurchase
	KontrakNontender		float64		`json:"kontrak_nontender"`
	// Paket nontender
	PaketNontender			int 		`json:"paket_nontender"`
	// Paket nontender Selesai
	PaketNontenderSelesai	int			`json:"paket_nontender_selesai"`
	// Persen ITKP non tender / non purchase
	PersenNontender			float32		`gorm:"type:numeric(10,2)" json:"persen_nontender"`
	// ITKP nontender / nonpurchase
	ItkpNontender			float32		`gorm:"type:numeric(10,2)" json:"itkp_nontender"`
	// Total Skor ITKP
	ItkpLokal				float32		`json:"itkp_lokal"`
}

func (Itkp) TableName() string {
	return "itkp"
}

func calculateItkp(v *Itkp) {
	v.ItkpSdm = 20 // set default itkp SDM
	// hitung ITKP RUP
	if v.Belanja > 0 {
		v.PersenRup = float32((v.PaguRup / v.Belanja) * 100)
	}
	v.ItkpRup = skorItkpRup(v.PersenRup)
	// hitung ITKP Ekontrak
	if v.PaketPemilihan > 0 {
 		v.PersenPemilihan = float32((v.PaketPemilihanSelesai / v.PaketPemilihan) * 100)
		v.ItkpPemilihan = skorItkpEkontrak(v.PersenPemilihan)
	}
	// hitung ITKP Tender
	if v.PaguTender > 0 {
		v.PersenTender = float32((v.KontrakTender / v.PaguTender) * 100)
		v.ItkpTender = skorItkpTender(v.PersenTender)
	}
	// hitung ITKP Purchase
	if v.Purchase > 0 {
		v.PersenPurchase = float32((v.PurchaseSelesai / v.Purchase) * 100)
		v.ItkpPurchase = skorItkpPurchase(v.PersenPurchase)
	}
	// hitung ITKP nontender
	if v.PaguNontender > 0 {
		v.PersenNontender = float32((v.KontrakNontender / v.PaguNontender) * 100)
		v.ItkpNontender = skorItkpNontender(v.PersenNontender)
	}
	v.ItkpLokal = v.ItkpSdm + v.ItkpRup + v.ItkpPemilihan + v.ItkpTender + v.ItkpPurchase + v.ItkpNontender
}

func GetItkp(tahun int) []Itkp {
	var result []Itkp
	db.Find(&result, "tahun=?", tahun)
	return result
}

 func UpdateItkpRup(tahun int) {
	var result []Itkp
	db.Raw(`WITH rup AS (SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, sum(pagu) pagu from rup WHERE tahun_anggaran = @tahun group by kd_satker, kd_satker_str, tahun_anggaran),
anggaran AS (SELECT id_satker kd_satker, tahun_anggaran tahun, sum(belanja_pengadaan) belanja FROM struktur_anggaran WHERE tahun_anggaran = @tahun  group by kd_satker, tahun_anggaran),
tender AS (SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, count(kd_tender) paket, sum(pagu) pagu FROM tender WHERE tahun_anggaran = @tahun group by kd_satker, kd_satker_str, tahun_anggaran),
tender_selesai AS (SELECT kd_satker, tahun_anggaran tahun, count(kd_tender) paket, sum(nilai_kontrak) kontrak FROM tender_selesai WHERE tahun_anggaran = @tahun group by kd_satker, tahun_anggaran),
nontender AS (SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, count(kd_nontender) paket, sum(pagu) pagu FROM nontender WHERE tahun_anggaran = @tahun group by kd_satker, kd_satker_str, tahun_anggaran),
nontender_selesai AS (SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, count(kd_nontender) paket, sum(nilai_kontrak) kontrak FROM nontender_selesai WHERE tahun_anggaran = @tahun  group by kd_satker, kd_satker_str, tahun_anggaran),
purchase AS (SELECT satker_id, tahun_anggaran tahun, count(*) paket FROM katalog_archive WHERE tahun_anggaran = @tahun group by tahun_anggaran, satker_id),
purchase_selesai AS (SELECT satker_id, tahun_anggaran tahun, count(*) paket FROM katalog_archive WHERE tahun_anggaran = @tahun AND paket_status_str='Paket Selesai' group by tahun_anggaran, satker_id)
SELECT a.kd_satker, a.kd_satker_str, a.tahun, a.pagu pagu_rup, belanja, c.pagu pagu_tender, d.kontrak kontrak_tender, c.paket paket_tender, d.paket paket_tender_selesai, e.pagu pagu_nontender , f.kontrak kontrak_nontender, e.paket paket_nontender, f.paket paket_nontender_selesai,
g.paket purchase, h.paket purchase_selesai, (c.paket + e.paket) paket_pemilihan, (d.paket + f.paket) paket_pemilihan_selesai
FROM rup a LEFT JOIN anggaran b ON a.kd_satker=b.kd_satker AND a.tahun = b.tahun
LEFT JOIN tender c ON a.kd_satker_str=c.kd_satker_str AND a.tahun = c.tahun
LEFT JOIN tender_selesai d ON a.kd_satker_str=d.kd_satker AND a.tahun = d.tahun
LEFT JOIN nontender e ON a.kd_satker_str=e.kd_satker_str AND a.tahun = e.tahun
LEFT JOIN nontender_selesai f ON a.kd_satker_str=f.kd_satker_str AND a.tahun = f.tahun
LEFT JOIN purchase g ON a.kd_satker=g.satker_id AND a.tahun = g.tahun
LEFT JOIN purchase_selesai h ON a.kd_satker=h.satker_id AND a.tahun = h.tahun
		`, map[string]interface{}{"tahun": tahun}).Scan(&result)
	log.Info("itkp tahun", tahun, "ada ", len(result))
	for i := range result  {
		calculateItkp(&result[i])
	}
	err := db.Save(&result).Error
	if err != nil {
		log.Error(err)
	}
	log.Info("Update ITKP Done")
 }

 func isBetween(x, min, max float32) bool {
	return x >= min && x <= max
 }
 /* skor ITKP RUP
  * 70% - 100% = 20
  * 50% - 70% = 15
  * <50% = 0
  * >100% = 0
  */

  func skorItkpRup(persen float32) float32 {
  	if isBetween(persen, 70,  100) {
   		return 20
  	} else if isBetween(persen, 50, 70) {
   		return 15
   	} else {
     	return 0
     }
  }
 /**
  * skor ITKP Ekontrak
  * 70% - 100% = 20
  * 50% - 70% = 15
  * 25% - 50% = 10
  * <25% = 0
  * >100% = 0
  */
  func skorItkpEkontrak(persen float32) float32 {
  	if isBetween(persen, 70, 100) {
   		return 20
  	} else if isBetween(persen, 50, 70) {
   		return 15
    } else if isBetween(persen, 25, 50) {
   		return 10
   	} else {
     	return 0
     }
  }
  /**
   * skor ITKP E-purchase
   * 70% - 100% = 20
   * 50% - 70% = 15
   * 25% - 50% = 10
   * <25% = 0
   * >100% = 0
   */
  func skorItkpPurchase(persen float32) float32 {
	if isBetween(persen, 70, 100) {
  		return 20
 	} else if isBetween(persen, 50, 70) {
  		return 15
   	} else if isBetween(persen, 25, 50) {
  		return 10
  	} else {
    	return 0
    }
  }

  /**
   * skor ITKP Tender
   * 70% - 100% = 10
   * 50% - 70% = 7,5
   * 25% - 50% = 5
   * <25% = 0
   * >100% = 0
   */
   func skorItkpTender(persen float32) float32 {
 		if isBetween(persen, 70, 100) {
    		return 10
   		} else if isBetween(persen, 50, 70) {
    		return 7.5
     	} else if isBetween(persen, 25, 50) {
    		return 5
    	} else {
      		return 0
      	}
   }
   /**
    * skor ITKP Non-tender
    * 70% - 100% = 10
    * 50% - 70% = 7,5
    * 25% - 50% = 5
    * <25% = 0
    * >100% = 0
    */
   func skorItkpNontender(persen float32) float32 {
 		if isBetween(persen, 70, 100) {
    		return 10
   		} else if isBetween(persen, 50, 70) {
    		return 7.5
     	} else if isBetween(persen, 25, 50) {
    		return 5
    	} else {
      		return 0
      	}
   }


	type ItkpRekap struct {
		KdSatker 		uint  		`json:"kd_satker"`
		Nama			string 		`json:"nama"`
		ItkpSdm			float32 	`json:"itkp_sdm"`
		ItkpRup			float32 	`json:"itkp_rup"`
		ItkpTender		float32 	`json:"itkp_tender"`
		ItkpPurchase	float32		`json:"itkp_purchase"`
		ItkpPemilihan	float32 	`json:"itkp_pemilihan"`
		ItkpNontender	float32		`json:"itkp_nontender"`
		ItkpNontener	float32 	`json:"itkp_lokal"`
	}

	type ItkpRup struct {
		KdSatker 		uint  		`json:"kd_satker"`
		Nama			string 		`json:"nama"`
		PaguRup			float64		`json:"pagu_rup"`
		Belanja			float64		`json:"belanja"`
		PersenRup		float32		`json:"persen_rup"`
		ItkpRup			float32		`json:"itkp_rup"`
	}

	type ItkpPemilihan struct {
		KdSatker 				uint  		`json:"kd_satker"`
		Nama					string 		`json:"nama"`
		PaketPemilihan			int 		`json:"paket_pemilihan"`
		PaketPemilihanSelesai	int 		`json:"paket_pemilihan_selesai"`
		PersenPemilihan			float32		`json:"persen_pemilihan"`
		ItkpPemilhan			float32		`json:"itkp_pemilihan"`
	}

	type ItkpTender struct {
		KdSatker 		uint  		`json:"kd_satker"`
		Nama			string 		`json:"nama"`
		PaguTender		float64		`json:"pagu_tender"`
		KontrakTender	float64		`json:"kontrak_tender"`
		PersenTender	float32		`json:"persen_tender"`
		ItkpTender		float32		`json:"itkp_tender"`
	}

	type ItkpNontender struct {
		KdSatker 			uint  		`json:"kd_satker"`
		Nama				string 		`json:"nama"`
		PaguNontender		float64		`json:"pagu_nontender"`
		KontrakNontender	float64		`json:"kontrak_nontender"`
		PersenNontender		float32		`json:"persen_nontender"`
		ItkpNontender		float32		`json:"itkp_nontender"`
	}

	type ItkpPurchase struct {
		KdSatker 		uint  		`json:"kd_satker"`
		Nama			string 		`json:"nama"`
		Purchase		int 		`json:"purchase"`
		PurchaseSelesai	int 		`json:"purchase_selesai"`
		PersenPurchase	float32		`json:"persen_purchase"`
		ItkpPurchase	float32		`json:"itkp_purchase"`
	}

   func GetDataTableItkp(c *fiber.Ctx, tahun int) error {
		var datas []ItkpRekap
		columns := []string{"itkp.kd_satker", "nama", "itkp_sdm", "itkp_rup", "itkp_tender",  "itkp_purchase", "itkp_pemilihan", "itkp_nontender", "itkp_lokal"}
		query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
		return GetDataTable(c, datas, columns, query, tahun)
   }

   func GetDataTableItkpRup(c *fiber.Ctx, tahun int) error {
   		var datas []ItkpRup
     	columns := []string{"itkp.kd_satker", "nama", "pagu_rup", "belanja", "persen_rup",  "itkp_rup"}
      	query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
       return GetDataTable(c, datas, columns, query, tahun)
   }

   func GetDataTableItkpPemilihan(c *fiber.Ctx, tahun int) error {
   		var datas []ItkpPemilihan
     	columns := []string{"itkp.kd_satker", "nama", "paket_pemilihan", "paket_pemilihan_selesai", "persen_pemilihan",  "itkp_pemilihan"}
      	query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
       return GetDataTable(c, datas, columns, query, tahun)
   }

   func GetDataTableItkpTender(c *fiber.Ctx, tahun int) error {
   		var datas []ItkpTender
     	columns := []string{ "itkp.kd_satker", "nama", "pagu_tender", "kontrak_tender", "persen_tender",  "itkp_tender"}
      	query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
       return GetDataTable(c, datas, columns, query, tahun)
   }

   func GetDataTableItkpNontender(c *fiber.Ctx, tahun int) error {
   		var datas []ItkpNontender
     	columns := []string{"itkp.kd_satker", "nama", "pagu_nontender", "kontrak_nontender", "persen_nontender",  "itkp_nontender"}
      	query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
       return GetDataTable(c, datas, columns, query, tahun)
   }

   func GetDataTableItkpPurchase(c *fiber.Ctx, tahun int) error {
   		var datas []ItkpPurchase
     	columns := []string{ "itkp.kd_satker", "nama", "purchase", "purchase_selesai", "persen_purchase",  "itkp_purchase"}
      	query := `FROM itkp LEFT JOIN satker ON satker.id = itkp.kd_satker WHERE tahun = ?`
        return GetDataTable(c, datas, columns, query, tahun)
   }

   type ItkpRupSatker struct {
    	KdSatker 			uint
        KdSatkerStr			string
        Tahun 				int
        KdRup				uint
        NamaPaket	 		string
        Pagu				float64
    }

    type ItkpTenderSatker struct {
    	KdSatker 			uint
        KdSatkerStr			string
        Tahun 				int
        KdTender			uint
        NamaPaket	 		string
        Pagu				float64
        Hps					float64
    }

    type ItkpNontenderSatker struct {
    	KdSatker 			uint
        KdSatkerStr			string
        Tahun 				int
        KdNontender			uint
        NamaPaket	 		string
        Pagu				float64
        Hps					float64
    }

    type ItkpPurchaseSatker struct {
      	SatkerId			uint
        Tahun 				int
        KdPaket				uint
        NamaPaket	 		string
        NoPaket				string
    }

    type ItkpEkontrakSatker struct {
		KdSatker 			uint
    	KdSatkerStr			string
       	Tahun 				int
        KdPaket				uint
        NamaPaket	 		string
        Pagu				float64
        Kontrak				float64
    }

    func GetDetilRupSatker(satkerId uint, satkerStr string, tahun int) []ItkpRupSatker {
      		var result []ItkpRupSatker
    		db.Raw(`SELECT kd_satker, kd_satker_str, kd_rup, tahun_anggaran tahun, nama_paket, pagu from rup WHERE tahun_anggaran = ? and (kd_satker = ? OR kd_satker_str = ?)`,
      			tahun, satkerId, satkerStr).Scan(&result)
      		return result
    }

      func GetDetilTenderSatker(satkerId string, satkerStr string, tahun int) []ItkpTenderSatker {
      		var result []ItkpTenderSatker
      		db.Raw(`SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, kd_tender, pagu, nama_paket, hps FROM tender WHERE tahun_anggaran = ? and (kd_satker = ? OR kd_satker_str = ?)`,
        		tahun, satkerId, satkerStr).Scan(&result)
        	return result
      }

      func GetDetilNontenderSatker(satkerId string, satkerStr string, tahun int) []ItkpNontenderSatker {
      		var result []ItkpNontenderSatker
      		db.Raw(`SELECT kd_satker, kd_satker_str, tahun_anggaran tahun, kd_nontender, pagu, nama_paket, hps FROM nontender WHERE tahun_anggaran = ? and (kd_satker = ? OR kd_satker_str = ?)`,
        		tahun, satkerId, satkerStr).Scan(&result)
        	return result
      }

      func GetDetilPurchaseSatker(satkerId uint, tahun int) []ItkpPurchaseSatker {
      		var result []ItkpPurchaseSatker
      		db.Raw(`SELECT satker_id, tahun_anggaran tahun, kd_paket, nama_paket, no_paket FROM katalog_archive WHERE tahun_anggaran = ? and satker_id = ?`,
        		tahun, satkerId).Scan(&result)
        	return result
      }

      func GetDetilEkontrakSatker(satkerId string, satkerStr string, tahun int) []ItkpEkontrakSatker {
      		var result []ItkpEkontrakSatker
      		db.Raw(`SELECT t.kd_satker, t.kd_satker_str, t.tahun_anggaran tahun, t.kd_tender kd_paket, nama_paket, t.pagu, nilai_kontrak kontrak FROM tender t, tender_selesai ts
        			WHERE ts.kd_tender = t.kd_tender AND t.tahun_anggaran = ? AND (t.kd_satker = ? OR t.kd_satker_str = ?)
           			UNION
              		SELECT ns.kd_satker, ns.kd_satker_str, ns.tahun_anggaran tahun, ns.kd_nontender kd_paket, ns.nama_paket, ns.pagu, nilai_kontrak kontrak FROM nontender n, nontender_selesai ns
                	WHERE ns.kd_nontender = n.kd_nontender AND ns.tahun_anggaran = ? AND (ns.kd_satker = ? OR ns.kd_satker_str = ?)`,
        		tahun, satkerId, satkerStr, tahun, satkerId, satkerStr).Scan(&result)
        	return result
      }
