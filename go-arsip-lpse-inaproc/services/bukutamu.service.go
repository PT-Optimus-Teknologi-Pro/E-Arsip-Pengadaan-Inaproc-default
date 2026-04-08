package services

import (
	"arsip/models"

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

func SaveFeedback(kualitas, fasilitas, kelengkapan []string, komentar string) error {
	return  models.SaveFeedback(kualitas, fasilitas, kelengkapan, komentar)
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

func GetAllFeedbacks() []models.Feedback {
	var feedbacks []models.Feedback
	// Only fetch feedbacks that have non-empty comments
	models.GetDB().Where("komentar IS NOT NULL AND komentar <> ''").Order("created_at desc").Limit(20).Find(&feedbacks)
	return feedbacks
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

