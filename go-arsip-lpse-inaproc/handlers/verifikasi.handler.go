package handlers

import (
	"arsip/services"
	"arsip/utils"
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
	if action == "approve" {
		if group == "" {
			return flashError(c, "Gagal: Anda harus memilih Pengangkatan Menjadi terlebih dahulu!", "/verifikasi/" + utils.UintToString(uint(id)))
		}
		user.Usrgroup = group
	}
	err := services.VerifikasiAkun(user, action, usrgroup)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tindakan pada Akun " + user.PegNama + " Gagal", "/verifikasi")
	}
	message := "Tindakan berhasil dialakukan."
	switch action {
	case "approve":
		message = "Verifikasi Akun " + user.PegNama + " Sukses! Data sekarang dapat dilihat di menu Pegawai."
	case "reject":
		message = "Penolakan Akun " + user.PegNama + " berhasil. Data tetap terekam sebagai Ditolak."
	case "delete":
		message = "Data pendaftar " + user.PegNama + " telah Dihapus Permanen dari sistem."
	}
	return flashSuccess(c, message, "/verifikasi")
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
	docMap := make(map[string]models.Document)
	for _, d := range documents {
		docMap[d.Jenis] = d
	}
	mp := currentMap(c)
	mp["pegawai"] = user
	mp["documents"] = documents
	mp["docMap"] = docMap
	return c.Render("verifikasi/verifikasi-detil", mp)
}

func GetVerifikasiView(c *fiber.Ctx) error {
	sess := getSession(c)
	id := utils.InterfaceToUint(sess.Get("id"))
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

func DeleteVerifikasi(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user := services.GetPegawai(uint(id))
	if user.ID == 0 {
		return c.SendStatus(404)
	}
	
	// Restriction: Only WAIT (0) or REJECT (3) accounts can be deleted
	if user.PegStatus != 0 && user.PegStatus != 3 {
		return flashError(c, "Hanya akun yang belum diverifikasi atau ditolak yang dapat dihapus.", "/verifikasi")
	}

	err := services.DeletePegawai(uint(id))
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Akun Pendaftar Gagal", "/verifikasi")
	}
	return flashSuccess(c, "Data pendaftar telah berhasil dihapus.", "/verifikasi")
}
