package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func EditPp(c *fiber.Ctx) error {
	mp := currentMap(c)
	mp["url"] = "/pp"
	id, _ := c.ParamsInt("id")
	tahun := c.QueryInt("tahun", time.Now().Year())
	mp["satkers"] = services.GetSatkerAPI(tahun)
	mp["anggotas"] = services.GetListAnggotaPp()
	if id != 0 {
		pp := services.GetPejabatPengadaan(uint(id))
		if pp.ID == 0 {
			return c.SendStatus(404)
		}
		mp["pp"] = pp
		mp["url"] = "/pp/" + utils.IntToString(id)
	} else {
		mp["pp"] = models.PejabatPengadaan{Tahun: tahun}
	}
	return c.Render("pp/form-pp", mp)
}

func CreatePp(c *fiber.Ctx) error {
	pp := new(models.PejabatPengadaanDTO)
	err := c.BodyParser(pp)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Pejabat Pengadaan Gagal","/pp/edit")
	}
	pp.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	pp.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	err = services.SavePejabatPengadaan(pp)
	if err != nil {
		log.Error(err)
		return flashError(c, "Tambah Pejabat Pengadaan Gagal", "/pp/edit")
	}
	return flashSuccess(c, "Tambah Pejabat Pengadaan Sukses","/pp")
}

// Get All Users from db
func GetAllPp(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("pp/pp", mp)
}

// GetSingleUser from db
func GetPp(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	pp := services.GetPejabatPengadaan(uint(id))
	if pp.ID == 0 {
		return c.SendStatus(404)
	}
	mp["pp"] = pp
	return c.Render("pp/pp-detil", mp)
}

// update a user in db
func UpdatePp(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var pp models.PejabatPengadaanDTO
	err := c.BodyParser(&pp)
	if err != nil {
		log.Error(err)
		return flashError(c, "Edit Pejabat Pengadaan Gagal", "/pp/edit/" + utils.IntToString(id))
	}
	pp.ID = uint(id)
	pp.PeriodeAwal, _ = time.Parse("2006-01-02", c.FormValue("periode_awal"))
	pp.PeriodeAkhir, _ = time.Parse("2006-01-02", c.FormValue("periode_akhir"))
	err = services.SavePejabatPengadaan(&pp)
	if err != nil {
		return flashError(c, err.Error(), "/pp/edit/" + utils.IntToString(id))
	}
	return flashSuccess(c, "Edit Pejabat Pengadaan Sukses","/pp")
}

// delete user in db by ID
func DeletePp(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	err := services.DeletePejabatPengadaan(uint(id))
	if err != nil {
		log.Error(err)
		return flashError(c, "Hapus Pejabat Pengadaan Gagal","/pp")
	}
	return flashSuccess(c, "Hapus Pejabat Pengadaan Sukses","/pp")
}

func GetJsonPp(c *fiber.Ctx) error {
	return services.GetDataTablePp(c)
}
