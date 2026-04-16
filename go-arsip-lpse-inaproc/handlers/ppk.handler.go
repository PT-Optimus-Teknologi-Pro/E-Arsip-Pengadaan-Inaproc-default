package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// GetPrivateDocumentPage renders the private document management page for any authenticated role
func GetPrivateDocumentPage(c *fiber.Ctx) error {
	mp := currentMap(c)
	// All authenticated users can access private docs

	roleLabel := "Pengguna"
	rolePath := "/ppk" // default fallback

	if mp["isPPK"].(bool) {
		roleLabel = "PPK"
		rolePath = "/ppk"
	} else if mp["isPokja"].(bool) {
		roleLabel = "Pokja"
		rolePath = "/pokja"
	} else if mp["isPP"].(bool) {
		roleLabel = "Pejabat Pengadaan"
		rolePath = "/pp"
	} else if mp["isUkpbj"].(bool) {
		roleLabel = "UKPBJ"
		rolePath = "/ukpbj"
	} else if mp["isAdmin"].(bool) {
		roleLabel = "Admin"
		rolePath = "/admin"
	} else if mp["isArsiparis"].(bool) {
		roleLabel = "Arsiparis"
		rolePath = "/ukpbj"
	} else if mp["isPegawai"].(bool) {
		roleLabel = "Pegawai"
		rolePath = "/pegawai"
	}
	mp["roleLabel"] = roleLabel
	mp["rolePath"] = rolePath

	return c.Render("ppk/dokumen-privat", mp)
}

// GetJsonPaketPrivate returns JSON data of packages owned/assigned to the logged-in user (PPK/Pokja/PP)
func GetJsonPaketPrivate(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	
	if mp["isPPK"].(bool) {
		return models.GetDataTablePaketPPK(c, userid)
	} else if mp["isPokja"].(bool) {
		return models.GetDataTablePaketPokja(c, userid)
	} else if mp["isPP"].(bool) {
		return models.GetDataTablePaketPP(c, userid)
	}
	
	return Forbiden(c)
}

// GetJsonPaketPrivateSirup returns only SiRUP-sourced packages for the PPK (rup_id > 0)
func GetJsonPaketPrivateSirup(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	if !mp["isPPK"].(bool) {
		return Forbiden(c)
	}
	return models.GetDataTablePaketPPKSirup(c, userid)
}

// GetJsonPaketPrivateMandiri returns manually-entered packages - filters by ppk_id for PPK, created_by for all other roles
func GetJsonPaketPrivateMandiri(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])
	if mp["isPPK"].(bool) {
		return models.GetDataTablePaketPPKMandiri(c, userid)
	}
	// For all other roles (Pokja, PP, UKPBJ, Admin, Arsiparis, Pegawai) filter by created_by
	return models.GetDataTablePaketMandiriByCreator(c, userid)
}

// SimpanPrivateDocument handles the upload of private documents
func SimpanPrivateDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	// All authenticated users may upload private documents
	
	pktId := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(pktId)
	userid := utils.InterfaceToUint(mp["id"])

	// Authorization Check: Check if user is associated with the package
	isAuthorized := false
	if mp["isPPK"].(bool) && paket.PpkId == userid {
		isAuthorized = true
	} else if mp["isPP"].(bool) && paket.PpId == userid {
		isAuthorized = true
	} else if mp["isPokja"].(bool) {
		// For Pokja, check if they are in the committee (panitia)
		if models.IsPegawaiInPanitia(userid, paket.PntId) {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		return Forbiden(c)
	}

	err := services.SimpanDokTambahanPrivate(c, pktId, userid)
	if err != nil {
		log.Error(err)
		return flashError(c, err.Error(), c.Get("Referer"))
	}
	return flashSuccess(c, "Simpan Dokumen Privat Berhasil", c.Get("Referer"))
}

// GetJsonPrivateDocuments returns a list of private documents for a package (strictly isolated for the uploader)
func GetJsonPrivateDocuments(c *fiber.Ctx) error {
	mp := currentMap(c)
	if !mp["isPPK"].(bool) && !mp["isPokja"].(bool) && !mp["isPP"].(bool) {
		return Forbiden(c)
	}
	
	pktId := utils.StringToUint(c.Params("id"))
	userid := utils.InterfaceToUint(mp["id"])

	// Authorization is strictly checked in the query: users only see their OWN private docs
	docs := models.GetDokTambahanPrivateListByUser(pktId, userid)
	
	var res []interface{}
	for _, d := range docs {
		res = append(res, fiber.Map{
			"id":        d.ID,
			"dok_id":    d.DokId,
			"filename":  d.Document().Filename,
			"filesize":  d.Document().Filesize,
			"CreatedAt": d.CreatedAt,
		})
	}

	return c.JSON(res)
}

// DeletePrivateDocument handles deletion of private documents (strictly for the owner)
func DeletePrivateDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	userid := utils.InterfaceToUint(mp["id"])

	id := utils.StringToUint(c.Params("id"))
	dokPaket := models.GetDokPaket(id)

	// Authorization check: Only the owner of the private document can delete it
	if dokPaket.PegId != userid || dokPaket.Jenis != models.TAMBAHAN_PRIVATE {
		return Forbiden(c)
	}

	err := models.DeleteDokPaket(&dokPaket)
	if err != nil {
		log.Error(err)
		return flashError(c, "Gagal menghapus dokumen", c.Get("Referer"))
	}
	
	return flashSuccess(c, "Dokumen privat berhasil dihapus", c.Get("Referer"))
}
