package handlers

import (
	"arsip/models"
	"arsip/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func Feedback(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["jenisLayanan"] = models.JENIS_LAYANAN
	mp["skorLayanan"] = models.SKOR_LAYANAN
	return c.Render("feedback/feedback", mp)
}

func FeedbackView(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	feedback := services.GetFeedback(uint(id))
	if feedback.ID == 0 {
		return c.SendStatus(404)
	}
	mp["feedback"] = feedback
	return c.Render("feedback/feedback-view", mp)
}

func SubmitFeedback(c *fiber.Ctx) error {
	log.Info("submit feedback")
	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		return flashError(c, "Input Feedback Gagal","/feedback")
	}
	kualitas := form.Value["kualitas[]"]
	fasilitas:= form.Value["fasilitas[]"]
	kelengkapan := form.Value["kelengkapan[]"]
	komentar := ""
	if len(form.Value["komentar"]) > 0 {
		komentar = form.Value["komentar"][0]
	}
	err = services.SaveFeedback(kualitas, fasilitas, kelengkapan, komentar)
	if err != nil {
		log.Error(err)
		return flashError(c, "Input Feedback Tamu Gagal","/feedback")
	}
	return flashSuccess(c, "Input Feedback tamu Sukses", "/feedback/success")
}

func FeedbackKonfirmasi(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["id_tamu"] = c.Params("id")
	return c.Render("feedback/feedback-konfirmasi", mp)
}

func FeedbackList(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["summary"] = services.GetSummaryFeedback()
	mp["allFeedbacks"] = services.GetAllFeedbacks()
	mp["totalResponses"] = services.GetTotalFeedbackResponses()
	return c.Render("feedback/feedback-list", mp)
}

func GetJsonFeedback(c *fiber.Ctx) error {
	return services.GetDataTableFeedback(c)
}
