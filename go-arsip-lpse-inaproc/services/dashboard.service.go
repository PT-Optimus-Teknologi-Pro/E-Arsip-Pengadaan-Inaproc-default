package services

import (
	"arsip/models"
	"arsip/utils"

	"github.com/gofiber/fiber/v2/log"
)

// DashboardHeroStats represents the aggregated data for the dashboard hero section
type DashboardHeroStats struct {
	TotalPagu      string
	Realisasi      string
	PaketSelesai   int64
	SkpdAktif      int
	RawTotalPagu   float64
	RawTotalBelanja float64
}

// GetDashboardHeroStats calculates all major statistics for the dashboard hero section
func GetDashboardHeroStats(tahun int) DashboardHeroStats {
	// 1. Get Progress from RUP/Structure
	progress := models.GetRupProgress(tahun)
	var totalPagu, totalBelanja float64
	var skpdAktif int
	for _, p := range progress {
		totalPagu += p.Pagu + p.PaguPds + p.PaguSwakelola
		totalBelanja += p.Belanja
		if (p.Pagu + p.PaguPds + p.PaguSwakelola) > 0 {
			skpdAktif++
		}
	}

	// 2. Get Paket Selesai Count
	// We combine logic from SPSE tables (Tender, Nontender, etc.)
	// and potentially internal package tracking if desired.
	var totalSelesai int64
	
	// SPSE Data Counts
	var count int64
	models.GetDB().Model(&models.Tender{}).Where("tahun_anggaran = ?", tahun).Count(&count)
	totalSelesai += count
	
	models.GetDB().Model(&models.Nontender{}).Where("tahun_anggaran = ?", tahun).Count(&count)
	totalSelesai += count
	
	models.GetDB().Model(&models.Pencatatan{}).Where("tahun_anggaran = ?", tahun).Count(&count)
	totalSelesai += count
	
	models.GetDB().Model(&models.Swakelola{}).Where("tahun_anggaran = ?", tahun).Count(&count)
	totalSelesai += count

	// Optional: Factor in internal archives that are marked as 'Selesai' (Status 6)
	// countInternal := GetCountPaketSelesaiInternal(tahun)
	// totalSelesai += countInternal

	return DashboardHeroStats{
		TotalPagu:      utils.FormatRupiah(totalPagu),
		Realisasi:      utils.Prosentase(totalBelanja, totalPagu),
		PaketSelesai:   totalSelesai,
		SkpdAktif:      skpdAktif,
		RawTotalPagu:   totalPagu,
		RawTotalBelanja: totalBelanja,
	}
}

// GetCountPaketSelesaiInternal returns the count of packages completed in the internal archive system
func GetCountPaketSelesaiInternal(tahun int) int64 {
	var count int64
	// Logic from ITKP: Join with paket_sirup to check the year
	err := models.GetDB().Table("paket").
		Joins("JOIN paket_sirup ON paket.rup_id = paket_sirup.id").
		Where("paket.status = 6 AND paket_sirup.tahun = ? AND paket.deleted_at IS NULL", tahun).
		Count(&count).Error
	
	if err != nil {
		log.Error("Error calculating internal paket selesai: ", err)
		return 0
	}
	return count
}

// GetCountSkpdAktif returns the number of Satkers that have RUP data in the given year
func GetCountSkpdAktif(tahun int) int {
	var count int64
	models.GetDB().Model(&models.SatkerSirup{}).
		Where("tahun_aktif LIKE ?", "%"+utils.IntToString(tahun)+"%").
		Count(&count)
	return int(count)
}

// GetDashboardRupProgress fetches the detailed progress per SKPD
func GetDashboardRupProgress(tahun int) []models.RupProgress {
	return models.GetRupProgress(tahun)
}
