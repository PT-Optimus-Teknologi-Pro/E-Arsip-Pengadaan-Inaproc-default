package services

import (
	"arsip/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetBukuTamu(id uint) models.BukuTamu {
	return models.GetBukuTamu(id)
}

func SaveDocBukuTamu(c *fiber.Ctx, buku *models.BukuTamu, name string) (uint, error) {
	buku.DokId, _ = models.SaveDocument(c, 0, models.BUKUTAMU, name)
	buku.Status = 0
	return models.SaveBukuTamu(buku)
}

func SaveBukuTamu(buku *models.BukuTamu) (uint, error) {
	return models.SaveBukuTamu(buku)
}


func GetFeedback(id uint) models.Feedback {
	return models.GetFeedback(id)
}

func SaveFeedback(nama, email, instansi, kategori string, kualitas, fasilitas, kelengkapan []string, komentar string) error {
	return models.SaveFeedback(nama, email, instansi, kategori, kualitas, fasilitas, kelengkapan, komentar)
}

func GetAllFeedbackKategori() []models.FeedbackKategori {
	var list []models.FeedbackKategori
	models.GetDB().Order("nama asc").Find(&list)
	return list
}

func SaveFeedbackKategori(nama string) error {
	return models.GetDB().Create(&models.FeedbackKategori{Nama: nama}).Error
}

func DeleteFeedbackKategori(id uint) error {
	return models.GetDB().Delete(&models.FeedbackKategori{}, id).Error
}

func GetSummaryFeedback() []models.SummaryFeedBack {
	var summary []models.SummaryFeedBack
	total := float64(GetTotalFeedbackResponses())
	if total == 0 { total = 1 } // Avoid division by zero

	for i := 0; i < 3; i++ {
		jenis := i + 1
		kDist := summaryFeedbackByKualitas(jenis)
		fDist := summaryFeedbackByFasilitas(jenis)
		keDist := summaryFeedbackByKelengkapan(jenis)

		obj := models.SummaryFeedBack{
			Jenis:          jenis,
			Kualitas:       kDist,
			Fasilitas:      fDist,
			Kelengkapan:    keDist,
			KualitasAvg:    float64(kDist[0]+kDist[1]*2+kDist[2]*3+kDist[3]*4+kDist[4]*5) / total,
			FasilitasAvg:   float64(fDist[0]+fDist[1]*2+fDist[2]*3+fDist[3]*4+fDist[4]*5) / total,
			KelengkapanAvg: float64(keDist[0]+keDist[1]*2+keDist[2]*3+keDist[3]*4+keDist[4]*5) / total,
		}
		summary = append(summary, obj)
	}
	return summary
}

func summaryFeedbackByKualitas(jenis int) [5]int64 {
	var summary [5]int64
	summary[0] = models.GetCountFeedbackByJenisKualitas(jenis, 1)
	summary[1] = models.GetCountFeedbackByJenisKualitas(jenis, 2)
	summary[2] = models.GetCountFeedbackByJenisKualitas(jenis, 3)
	summary[3] = models.GetCountFeedbackByJenisKualitas(jenis, 4)
	summary[4] = models.GetCountFeedbackByJenisKualitas(jenis, 5)
	return summary
}

func summaryFeedbackByFasilitas(jenis int) [5]int64 {
	var summary [5]int64
	summary[0] = models.GetCountFeedbackByJenisFasilitas(jenis, 1)
	summary[1] = models.GetCountFeedbackByJenisFasilitas(jenis, 2)
	summary[2] = models.GetCountFeedbackByJenisFasilitas(jenis, 3)
	summary[3] = models.GetCountFeedbackByJenisFasilitas(jenis, 4)
	summary[4] = models.GetCountFeedbackByJenisFasilitas(jenis, 5)
	return summary
}

func summaryFeedbackByKelengkapan(jenis int) [5]int64 {
	var summary [5]int64
	summary[0] = models.GetCountFeedbackByJenisKelengkapan(jenis, 1)
	summary[1] = models.GetCountFeedbackByJenisKelengkapan(jenis, 2)
	summary[2] = models.GetCountFeedbackByJenisKelengkapan(jenis, 3)
	summary[3] = models.GetCountFeedbackByJenisKelengkapan(jenis, 4)
	summary[4] = models.GetCountFeedbackByJenisKelengkapan(jenis, 5)
	return summary
}

func GetAllFeedbacks(filter string) []models.Feedback {
	var feedbacks []models.Feedback
	// Base query: Only fetch feedbacks that have non-empty comments
	db := models.GetDB().Where("komentar IS NOT NULL AND komentar <> '' AND jenis = 1")

	// Apply time-based filters
	switch filter {
	case "baru":
		// New: records from the last 7 days
		db = db.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "lama":
		// Old: records older than 7 days
		db = db.Where("created_at < ?", time.Now().AddDate(0, 0, -7))
	}

	db.Order("created_at desc").Find(&feedbacks)
	return feedbacks
}

func DeleteFeedback(id uint) error {
	return models.GetDB().Delete(&models.Feedback{}, id).Error
}

func GetTotalFeedbackResponses() int64 {
	return models.GetTotalFeedbackCount()
}

func GetGlobalAverageScore() float64 {
	summary := GetSummaryFeedback()
	if len(summary) == 0 { return 0 }
	
	var totalAvg float64
	for _, s := range summary {
		totalAvg += (s.KualitasAvg + s.FasilitasAvg + s.KelengkapanAvg) / 3
	}
	return totalAvg / float64(len(summary))
}

