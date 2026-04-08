package handlers

import (
	"arsip/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateVerifikasi(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return c.SendStatus(404)
	}
	sess := getSession(c)
	usrgroup := sess.Get("group").(string)
	action := c.FormValue("action")
	group := c.FormValue("usrgroup")
	user.Usrgroup = group
	err := services.VerifikasiAkun(user, action, usrgroup)
	if err != nil {
		log.Error(err)
		return flashError(c, "Verifikasi Akun " + user.PegNama + " Gagal", "/verifikasi")
	}
	return flashSuccess(c, "Verifikasi Akun " + user.PegNama + " Sukses","/verifikasi")
}

func GetAllVerifikasi(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("verifikasi/verifikasi", mp)
}

func GetJsonVerifikasi(c *fiber.Ctx) error {
	return services.GetDataTableVerifikasi(c)
}

func GetVerifikasi(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return c.SendStatus(404)
	}
	documents := services.GetDocumentPegawai(uint(id))
	mp := currentMap(c)
	mp["pegawai"] = user
	mp["documents"] = documents
	return c.Render("verifikasi/verifikasi-detil", mp)
}

func GetVerifikasiView(c *fiber.Ctx) error {
	sess := getSession(c)
	id := sess.Get("id").(uint)
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return c.SendStatus(404)
	}
	documents := services.GetDocumentPegawai(uint(id))
	mp := currentMap(c)
	mp["pegawai"] = user
	mp["documents"] = documents
	return c.Render("verifikasi/verifikasi-view", mp)
}
