package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func Feedback(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["jenisLayanan"] = models.JENIS_LAYANAN
	mp["skorLayanan"] = models.SKOR_LAYANAN
	mp["kategoriList"] = services.GetAllFeedbackKategori()
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

	nama := ""
	if len(form.Value["nama"]) > 0 { nama = form.Value["nama"][0] }
	email := ""
	if len(form.Value["email"]) > 0 { email = form.Value["email"][0] }
	instansi := ""
	if len(form.Value["nama_perusahaan"]) > 0 { instansi = form.Value["nama_perusahaan"][0] }
	kategori := ""
	if len(form.Value["kategori"]) > 0 { kategori = form.Value["kategori"][0] }

	kualitas := form.Value["kualitas[]"]
	fasilitas:= form.Value["fasilitas[]"]
	kelengkapan := form.Value["kelengkapan[]"]
	komentar := ""
	if len(form.Value["komentar"]) > 0 {
		komentar = form.Value["komentar"][0]
	}
	err = services.SaveFeedback(nama, email, instansi, kategori, kualitas, fasilitas, kelengkapan, komentar)
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
	filter := c.Query("filter", "semua")
	mp["filter"] = filter
	mp["summary"] = services.GetSummaryFeedback()
	mp["allFeedbacks"] = services.GetAllFeedbacks(filter)
	total := services.GetTotalFeedbackResponses()
	mp["totalResponses"] = total
	globalAvg := services.GetGlobalAverageScore()
	mp["globalAvg"] = globalAvg
	mp["globalAvgPercent"] = int(globalAvg * 20)
	return c.Render("feedback/feedback-list", mp)
}

func DeleteFeedback(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	err := services.DeleteFeedback(id)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Feedback Gagal", "/feedback/list")
	}
	return flashSuccess(c, "Hapus Feedback Berhasil", "/feedback/list")
}

func GetJsonFeedback(c *fiber.Ctx) error {
	return services.GetDataTableFeedback(c)
}

func FeedbackKategoriList(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["list"] = services.GetAllFeedbackKategori()
	return c.Render("feedback/kategori-list", mp)
}

func SubmitFeedbackKategori(c *fiber.Ctx) error {
	nama := c.FormValue("nama")
	if nama == "" {
		return flashError(c, "Nama Kategori tidak boleh kosong", "/feedback/kategori")
	}
	err := services.SaveFeedbackKategori(nama)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal simpan kategori", "/feedback/kategori")
	}
	return flashSuccess(c, "Kategori berhasil disimpan", "/feedback/kategori")
}

func DeleteFeedbackKategori(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	err := services.DeleteFeedbackKategori(id)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal hapus kategori", "/feedback/kategori")
	}
	return flashSuccess(c, "Kategori berhasil dihapus", "/feedback/kategori")
}
