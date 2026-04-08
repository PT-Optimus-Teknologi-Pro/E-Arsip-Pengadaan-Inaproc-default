package services

import (
	"arsip/config"
	"arsip/models"
)


func GetSatker(id uint) models.SatkerSirup {
	return models.GetSatkerSirup(id)
}


func GetTahunList() []int {
	return config.GetTahunList()
}

func GetRupProgress(tahun int) []models.RupProgress {
	return models.GetRupProgress(tahun)
}

func GetPaketPrioritas(tahun int) []models.PaketPrioritas {
	return models.GetPaketPrioritas(tahun)
}

func GetBebanPersonil(tahun int) []models.BebanPersonil {
	return models.GetBebanPersonil(tahun)
}

func GetRekapPaketSatker(tahun int) []models.RekapPaketSatker {
	return models.GetRekapPaketSatker(tahun)
}

func GetRekapFeedback(tahun int) models.RekapFeedbackDashboard {
	tidakPuas := make([]int, 12)
	puas := make([]int, 12)
	sangatPuas := make([]int, 12)
	rekap := models.GetRekapFeedbackBulanan(tahun)
	for _, v := range rekap {
		if v.Kepuasan == 1 {
			tidakPuas[v.Bulan - 1] = v.Count
		}
		if v.Kepuasan == 2 {
			puas[v.Bulan - 1] = v.Count
		}
		if v.Kepuasan == 3 {
			sangatPuas[v.Bulan - 1] = v.Count
		}
	}
	return models.RekapFeedbackDashboard{
		Tahun: tahun,
		TidakPuas: tidakPuas,
		Puas: puas,
		SangatPuas: sangatPuas,
	}
}

func CalculateRekapSatker(tahun int) []models.RekapSatkerDashboard {
	satkers := models.GetRekapSatkerDashboard(tahun)
	result := make([]models.RekapSatkerDashboard, len(satkers))
	rencanaRekap := models.GetRekapRencanaSatkerBulanan(tahun)
	tenderRekap := models.GetRealisasiTenderPerSkpd(tahun)
	nontenderRekap := models.GetRealisasiNontenderPerSkpd(tahun)
	pencatatanRekap := models.GetRealisasiPencatatanPerSkpd(tahun)
	purchaseRekap := models.GetRealisasiPurchasePerSkpd(tahun)
	for i, obj := range satkers {
		result[i].KdSatker = obj.KdSatker
		result[i].NamaSatker = obj.NamaSatker
		result[i].Tahun = tahun
		result[i].Rencana = make([]float64, 12)
		result[i].PaketRencana = make([]int, 12)
		result[i].Realisasi = make([]float64, 12)
		result[i].PaketRealisasi = make([]int, 12)
		result[i].Purchase = make([]float64, 12)
		result[i].PaketPurchase = make([]int, 12)
		for _, v := range rencanaRekap {
			if result[i].KdSatker == v.KdSatker {
				result[i].Rencana[v.Bulan-1] = v.Pagu
				result[i].PaketRencana[v.Bulan-1] = v.Paket
			}
		}
		for _, v := range tenderRekap {
			if result[i].KdSatker == v.KdSatker {
				result[i].Realisasi[v.Bulan-1] += v.Realisasi
				result[i].PaketRealisasi[v.Bulan-1] += v.Paket
			}
		}
		for _, v := range nontenderRekap {
			if result[i].KdSatker == v.KdSatker {
				result[i].Realisasi[v.Bulan-1] += v.Realisasi
				result[i].PaketRealisasi[v.Bulan-1] += v.Paket
			}
		}
		for _, v := range pencatatanRekap {
			if result[i].KdSatker == v.KdSatker {
				result[i].Realisasi[v.Bulan-1] += v.Realisasi
				result[i].PaketRealisasi[v.Bulan-1] += v.Paket
			}
		}
		for _, v := range purchaseRekap {
			if result[i].KdSatker == v.KdSatker {
				result[i].Purchase[v.Bulan-1] += v.Realisasi
				result[i].PaketPurchase[v.Bulan-1] += v.Paket
			}
		}
	}
	return result
}

func GetRekapPaketPpk(tahun int) []models.PPKRekap {
	rekap := make(map[string]models.PPKRekap)
	tender := models.GetPaketTenderPpk(tahun)
	for _, v := range tender {
		rekap[v.NipPpk] = v
	}
	nontender := models.GetPaketNontenderPpk(tahun)
	for _, v := range nontender {
		obj, ok := rekap[v.NipPpk]
		if ok {
			obj.Count += v.Count
		} else {
			obj = v
		}
		rekap[v.NipPpk] = obj
	}
	pencatatan := models.GetPaketPencatatanPpk(tahun)
	for _, v := range pencatatan {
		obj, ok := rekap[v.NipPpk]
		if ok {
			obj.Count += v.Count
		} else {
			obj = v
		}
		rekap[v.NipPpk] = obj
	}
	swakelola := models.GetPaketSwakelolaPpk(tahun)
	for _, v := range swakelola {
		obj, ok := rekap[v.NipPpk]
		if ok {
			obj.Count += v.Count
		} else {
			obj = v
		}
		rekap[v.NipPpk] = obj
	}
	purchase := models.GetPaketPurchasePpk(tahun)
	for _, v := range purchase {
		obj, ok := rekap[v.NipPpk]
		if ok {
			obj.CountPurchase += v.Count
		} else {
			obj = v
		}
		rekap[v.NipPpk] = obj
	}
	values := make([]models.PPKRekap, 0, len(rekap))
	for key , value := range rekap {
		if key != ""{
			values = append(values, value)
		}
	}
	return values
}

func GetRekapRencanaPaketSatkerBulan(tahun int, satkerId uint, bulan int) []models.Rup {
	return models.GetRencanaPaketSatkerBulanan(tahun, satkerId, bulan)
}

func GetRekapRealisasiPaketSatkerBulan(tahun int, satkerId uint, bulan int) []models.Rup {
	return models.GetRencanaPaketSatkerBulanan(tahun, satkerId, bulan)
}

func GetRekapRealisasiPaketTenderSatkerBulan(tahun int, satkerId uint, bulan int) []models.DetailPaketBulanan {
	return models.GetRealisasiPaketTenderBulanan(tahun, satkerId, bulan)
}

func GetRekapRealisasiPaketNontenderSatkerBulan(tahun int, satkerId uint, bulan int) []models.DetailPaketBulanan {
	return models.GetRealisasiPaketNontenderBulanan(tahun, satkerId, bulan)
}

func GetRekapRealisasiPencatatanSatkerBulan(tahun int, satkerId uint, bulan int) []models.DetailPaketBulanan {
	return models.GetRealisasiPencatatanBulanan(tahun, satkerId, bulan)
}

func GetRekapRealisasiPurchaseSatkerBulan(tahun int, satkerId uint, bulan int) []models.DetailPaketBulanan {
	return models.GetRealisasiPurchaseBulanan(tahun, satkerId, bulan)
}

func GetRekapTenderPPk(nip string, tahun int) []models.Tender {
	return models.GetDetilPaketTenderPpk(nip, tahun)
}

func GetRekapNontenderPPk(nip string, tahun int) []models.Nontender {
	return models.GetDetilPaketNontenderPpk(nip, tahun)
}

func GetRekapPencacatanPPk(nip string, tahun int) []models.Pencatatan {
	return models.GetDetilPaketPencatatanPpk(nip, tahun)
}

func GetRekapPurchasePPK(nip string, tahun int) []models.Katalog {
	return models.GetDetilPaketPurchasePpk(nip, tahun)
}
