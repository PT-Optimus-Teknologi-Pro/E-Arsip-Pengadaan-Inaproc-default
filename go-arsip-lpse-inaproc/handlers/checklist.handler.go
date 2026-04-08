package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditChecklist(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := strconv.Atoi(c.Params("id"))
	mp["url"] = "/checklist"
	mp["jenisPengadaan"] = services.GetJenisPengadaan()
	mp["metodePengadaan"] = services.GetMetodePengadaan()
	mp["templateList"] = services.GetAllDokTemplate()
	if id > 0 {
		checklist := services.GetChecklist(uint(id))
		if checklist.ID == 0 {
			return c.SendStatus(404)
		}
		mp["checklist"] = checklist
		mp["checklistDok"] = checklist.ChecklistDok()
	}
	return c.Render("checklist/form-checklist", mp)
}

type SyaratChecklist struct {
	Syarat []uint `form:"syarat"`
}

func CreateChecklist(c *fiber.Ctx) error {
	checklist := new(models.Checklist)
	err := c.BodyParser(checklist)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah checklist Gagal", "/checklist/edit")
	}
	checklist.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	checklist.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	err = models.SaveChecklist(checklist)
	syarat := new(SyaratChecklist)
	err = c.BodyParser(syarat)
	res := []models.ChecklistDok{}
	for _,v := range syarat.Syarat {
	// 	checklist.Jenis = jenis
	// 	checklist.DokTemplate = v
		// checklist.Required = true
		res = append(res, models.ChecklistDok{
			ChkId: checklist.ID,
			DokId: v,
			Status: 1,
		})
	}
	err = services.SimpanChecklistPersyaratan(res)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah checklist Gagal","/checklist/edit")
	}

	return flashSuccess(c, "Tambah checklist Sukses","/checklist")
}

// Get All Users from db
func GetAllChecklist(c *fiber.Ctx) error {
	mp := currentMap(c)
	// mp["checklists"] = services.GetChecklists()
	return c.Render("checklist/checklist", mp)
}

// GetSingleUser from db
func GetChecklist(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	mp := currentMap(c)
	checklist := services.GetChecklist(uint(id))
	// mp["checklists"] = checklists
	mp["checklist"] =  checklist
	return c.Render("checklist/checklist-detil", mp)
}

// delete user in db by ID
func DeleteChecklist(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	err := services.HapusChecklist(id)
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus checklist Gagal","/checklist")
	}
	return flashSuccess(c, "Hapus checklist Sukses","/checklist")
}

func GetJsonChecklist(c *fiber.Ctx) error {
	return services.GetDataTableChecklist(c)
}
