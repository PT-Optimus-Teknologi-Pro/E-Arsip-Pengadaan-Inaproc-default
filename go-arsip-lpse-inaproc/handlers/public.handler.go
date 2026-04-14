package handlers

import (
	"arsip/cache"
	"arsip/config"
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"bytes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/kal72/go-captcha"
	"github.com/sujit-baniya/flash"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var Sessions *session.Store

func Design(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("index", mp)
}

func Register(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("publik/register", mp)
}

func SubmitRegister(c *fiber.Ctx) error {
	user := new(models.Pegawai)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	if err != nil {
		log.Error(err)
		return flashError(c, "Registrasi Akun Gagal","/register")
	}
	plainPassw := user.Passw
	strength := passwordvalidator.GetEntropy(plainPassw)

    log.Info("Password: %s\n", plainPassw)
    log.Info("Entropy: %.2f bits\n", strength)
    if strength < 60 {
        log.Info("Password is too weak")
        return flashError(c, "Password is too weak","/register")
    } else {
        log.Info("Password is strong")
    }
	user.PegMasaBerlaku, _ = time.Parse("2006-01-02", c.FormValue("masa_berlaku"))
	user.PegIsactive = models.VERIFIKASI
	err = services.Registrasi(c, user, plainPassw)
	if err != nil {
		log.Error(err)
		return flashError(c, "Registrasi Akun Gagal","/register")
	}
	return flashSuccess(c, "Registrasi Akun Sukses","/register/success")
}

func RegisterKonfirmasi(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("publik/register-konfirmasi", mp)
}

func Login(c *fiber.Ctx) error {
	mp := currentMap(c)
	cap := captcha.New("qwertyasdfzxcv1234")
	base64Image, text, token, err := cap.Generate()
	if err != nil {
		log.Error(err)
	}
	// log.Info("text ", text)
	// log.Info("token ", token)
	key := "capthca"+token
	cache.Set(key, text)
	mp["token"] = token
	mp["base64Image"] = base64Image
	return c.Render("publik/login", mp)
}


func PrintToPdf(c *fiber.Ctx) error {
	result := utils.ExportToPdf("https://google.com")
	reader := bytes.NewReader(result)
	// Set the Content-Type header for PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	return c.SendStream(reader)
}


func Download(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	document := services.GetDocument(uint(id))
	return c.SendFile(document.Filepath)
}

func DownloadAll(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id") // id paket
	paket := services.GetPaket(uint(id))
	if !services.AuthorisasiPaket(paket, mp) {
		return Forbiden(c)
	}
	log.Info("create zip file");
	var files []string
	isPPK := mp["isPPK"].(bool)
	for _,v := range paket.GetAllDocument(isPPK) {
		files = append(files, v.Filepath)
	}
	zipFile,  err := utils.CreateZip(files, "download-all.zip");
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
	}
	return c.SendFile(zipFile)
}

func SubmitLogin(c *fiber.Ctx) error {
	token := c.FormValue("token")
	captchaText := c.FormValue("captchaText")
	key := "capthca"+token
	text, found := cache.Get(key)
	if found {
		cache.Delete(key)
	}
	log.Info("text ", text)
	log.Info("captchaText ", captchaText)
	if text != captchaText {
		return flashError(c,  "invalid Captcha","/login")
	}
	userid := c.FormValue("userid")
	password := c.FormValue("password")
	user, err := services.Otentikasi(userid, password)
	if err != nil {
		return flashError(c,  err.Error(),"/login")
	}
	sess := getSession(c)
	sess.Set("id", user.ID)
	sess.Set("name", user.PegNama)
	sess.Set("group", user.Usrgroup)
	defer sess.Save()
	return c.Redirect("/home")
}

func Logout(c *fiber.Ctx) error {
	sess := getSession(c)
	sess.Delete("id")
	sess.Delete("name")
	sess.Delete("group")
	// Destry session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	return c.Redirect("/")
}

func Home(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	if mp["id"] != nil {
	id := mp["id"].(uint)
		mp["countAgency"] = services.GetCountAgency()
		mp["countUkpbj"] = services.GetCountUkpbj()
		mp["countPegawai"] = services.GetCountPegawai()
		mp["user"] = services.GetPegawai(id)
		if mp["isPPK"].(bool) {
			return c.Render("beranda/home-ppk", mp)
		}
		if mp["isPP"].(bool) {
			return c.Render("beranda/home-pp", mp)
		}
		if mp["isPokja"].(bool) || mp["isPegawai"].(bool) {
			return c.Render("beranda/home-pokja", mp)
		}
	}
	mp["tahun"] = tahun
	tahunList, found := cache.Get("tahunList")
	if !found {
		tahunList = config.GetTahunList()
		cache.Set("tahunList", tahunList)
	}
	rekapFeedback, found := cache.Get("rekapFeedback")
	if !found {
		rekapFeedback = services.GetRekapFeedback(tahun)
		cache.Set("rekapFeedback", rekapFeedback)
	}
	mp["rekapFeedback"] = rekapFeedback
	mp["tahunlist"] = tahunList

	var sliders []models.HeroSlider
	models.GetDB().Order("id desc").Find(&sliders)
	for i := range sliders {
		sliders[i].Image = utils.ToWebPath(sliders[i].Image)
	}
	mp["heroSliders"] = sliders

	// Use the new Dashboard Service for statistics
	stats := services.GetDashboardHeroStats(tahun)

	mp["heroTotalPagu"] = stats.TotalPagu
	mp["heroRealisasi"] = stats.Realisasi
	mp["heroPaketSelesai"] = stats.PaketSelesai
	mp["heroSkpdAktif"] = stats.SkpdAktif

	appSettings := mp["app_settings"].(models.AppSettings)
	mp["heroBadge"] = appSettings.HeroBadge
	mp["heroTitle"] = appSettings.HeroTitle
	mp["heroSubtitle"] = appSettings.HeroSubtitle

	return c.Render("beranda/home", mp)
}

func Profile(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := mp["id"].(uint)
	user := services.GetPegawai(id)
	mp["user"] = user
	return c.Render("profile", mp)
}


func LoggedMiddleware(c *fiber.Ctx) error {
	sess := getSession(c)
	id := sess.Get("id")
	if id == nil {
		return Forbiden(c)
	}
	return c.Next()
}

func getSession(c *fiber.Ctx) *session.Session {
	sess, err := Sessions.Get(c)
	if err != nil {
		panic(err)
	}
	return sess
}

func currentMap(c *fiber.Ctx) fiber.Map {
	mp := flash.Get(c)
	if mp == nil {
		mp = make(fiber.Map)
	}
	sess := getSession(c)
	if sess != nil {
		mp["id"] = sess.Get("id")
		mp["name"] = sess.Get("name")
		mp["group"] = sess.Get("group")
		mp["isPPK"] = sess.Get("group") == models.PPK
		mp["isPokja"] = sess.Get("group") == models.POKJA
		mp["isPP"] = sess.Get("group") == models.PP
		mp["isUkpbj"] = sess.Get("group") == models.UKPBJ
		mp["isAdmin"] = sess.Get("group") == models.ADMIN
		mp["isAdminAgency"] = sess.Get("group") == models.ADM_AGENCY
		mp["isPegawai"] = sess.Get("group") == models.PEGAWAI
		mp["isArsiparis"] = sess.Get("group") == models.ARSIPARIS
	}
	mp["path"] = string(c.Request().URI().Path())

	appSettings := services.GetSettings()

	logoPath := appSettings.LogoPath
	if logoPath == "" { logoPath = "/modern/images/logo.png" }
	logoPath = utils.ToWebPath(logoPath)

	faviconPath := appSettings.FaviconPath
	if faviconPath == "" { faviconPath = "/favicon.ico" }
	faviconPath = utils.ToWebPath(faviconPath)

	appTitle := appSettings.AppTitle
	if appTitle == "" { appTitle = "e-Arsip Pengadaan" }

	mp["app_title"] = appTitle
	mp["appTitle"] = appTitle
	mp["appSubtitle"] = appSettings.LoadingSubtitle
	mp["logo_path"] = logoPath
	mp["logoPath"] = logoPath
	mp["favicon_path"] = faviconPath
	mp["faviconPath"] = faviconPath
	mp["app_settings"] = appSettings

	return mp
}


func getUserSession(c *fiber.Ctx) models.UserSession {
	sess := getSession(c)
	var result models.UserSession
	if sess != nil {
		result = models.UserSession{
			Id: sess.Get("id").(uint),
			Name: sess.Get("name").(string),
			Role: sess.Get("group").(string),
		}
	}
	return result
}

func flashError(c *fiber.Ctx, message,redirect string) error {
	return flash.WithError(c, fiber.Map{"message": message, "error": true, "success": false}).Redirect(redirect)
}

func flashSuccess(c *fiber.Ctx, message,redirect string) error {
	return flash.WithSuccess(c, fiber.Map{"message": message, "error": false, "success": true}).Redirect(redirect)
}

func flashSuccessWithData(c *fiber.Ctx, message, data, redirect string) error {
	return flash.WithSuccess(c, fiber.Map{"message": message, "error": false, "data": data}).Redirect(redirect)
}

func FlushCache(c *fiber.Ctx) error {
	cache.Flush()
	return c.Redirect("/")
}

func Forbiden(c *fiber.Ctx) error {
	log.Info("forbidden akses")
	return c.Status(fiber.StatusForbidden).Render("errors/error-403", fiber.Map{})
}
func GetFooterSocials() []map[string]interface{} {
	var res []map[string]interface{}

	// Static entries from AppSettings
	appSettings := services.GetSettings()
	if appSettings.FooterFacebook != "" && appSettings.FooterFacebook != "#" {
		res = append(res, map[string]interface{}{"label": "Facebook", "url": appSettings.FooterFacebook, "icon": "facebook"})
	}
	if appSettings.FooterInstagram != "" && appSettings.FooterInstagram != "#" {
		res = append(res, map[string]interface{}{"label": "Instagram", "url": appSettings.FooterInstagram, "icon": "instagram"})
	}
	if appSettings.FooterTwitter != "" && appSettings.FooterTwitter != "#" {
		res = append(res, map[string]interface{}{"label": "Twitter", "url": appSettings.FooterTwitter, "icon": "twitter"})
	}

	// Dynamic entries from footer_social_links table
	var socialLinks []models.FooterSocialLink
	models.GetDB().Order("sort asc, id asc").Find(&socialLinks)
	for _, s := range socialLinks {
		res = append(res, map[string]interface{}{
			"label": s.Label,
			"url":   s.URL,
			"icon":  utils.ToWebPath(s.Icon),
		})
	}
	return res
}

func GetFooterQuicks() []map[string]interface{} {
	var quickLinks []models.FooterQuickLink
	models.GetDB().Order("sort asc, id asc").Find(&quickLinks)
	var res []map[string]interface{}
	for _, q := range quickLinks {
		res = append(res, map[string]interface{}{
			"label": q.Label,
			"url":   q.URL,
		})
	}
	return res
}

func GetFooterServices() []map[string]interface{} {
	var serviceLinks []models.FooterService
	models.GetDB().Order("sort asc, id asc").Find(&serviceLinks)
	var res []map[string]interface{}
	for _, s := range serviceLinks {
		res = append(res, map[string]interface{}{
			"label": s.Label,
			"url":   s.URL,
		})
	}
	return res
}

func UploadAdminDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}
	userid := mp["id"].(uint)
	_, err := models.SaveDocument(c, userid, models.ADMIN_DOK, "file")
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal mengupload dokumen: "+err.Error(), "/admin-document")
	}
	return flashSuccess(c, "Dokumen berhasil diupload", "/admin-document")
}

func GetAllAdminDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}
	return c.Render("document/admin-index", mp)
}

func GetJsonAdminDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}
	return services.GetDataTableAdminDocument(c)
}

func DeleteAdminDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}
	docId, _ := c.ParamsInt("id")
	doc := services.GetDocument(uint(docId))
	if doc.ID == 0 || doc.Jenis != models.ADMIN_DOK {
		return flashError(c, "Dokumen tidak ditemukan", "/admin-document")
	}
	err := services.DeleteDocument(doc)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal menghapus dokumen: "+err.Error(), "/admin-document")
	}
	return flashSuccess(c, "Dokumen berhasil dihapus", "/admin-document")
}

func DownloadAllAdminDocument(c *fiber.Ctx) error {
	log.Info("DownloadAllAdminDocument hit")
	mp := currentMap(c)
	if !mp["isAdmin"].(bool) {
		return Forbiden(c)
	}
	log.Info("create zip file for admin documents")
	var files []string
	docs := services.GetAdminDocuments()
	for _, v := range docs {
		files = append(files, v.Filepath)
	}

	if len(files) == 0 {
		return flashError(c, "Tidak ada dokumen untuk didownload", "/admin-document")
	}

	zipFile, err := utils.CreateZip(files, "dokumen-admin.zip")
	if err != nil {
		log.Error("Error creating ", zipFile, " : ", err)
		return flashError(c, "Gagal membuat file zip: "+err.Error(), "/admin-document")
	}
	return c.SendFile(zipFile)
}
