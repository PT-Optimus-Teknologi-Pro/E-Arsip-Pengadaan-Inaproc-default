package services

import (
	"arsip/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	for i := range 3 {
		jenis := i + 1;
		log.Info("summary jenis", jenis)
		obj := models.SummaryFeedBack {
			Jenis: jenis,
			Fasilitas: summaryFeedbackByFasilitas(jenis),
			Kualitas: summaryFeedbackByKualitas(jenis),
			Kelengkapan: summaryFeedbackByKelengkapan(jenis),
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
	// Only fetch feedbacks that have comments and are not deleted
	models.GetDB().Where("komentar <> ''").Order("created_at desc").Limit(30).Find(&feedbacks)
	return feedbacks
}

func GetTotalFeedbackResponses() int64 {
	var count int64
	// Since each form submission creates 3 records in the Feedback table (one for each service), 
	// we count by grouping or just counting total and dividing by 3? 
	// Actually, let's just count total records for simplicity or count unique names if applicable.
	// But let's just count how many distinct "created_at" timestamps we have roughly, or just total count.
	models.GetDB().Model(&models.Feedback{}).Count(&count)
	return count / 3 // Assuming 3 types of services per form
}
