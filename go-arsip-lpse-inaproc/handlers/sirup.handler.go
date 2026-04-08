package handlers

import (
	"arsip/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetJsonPaketSirup(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.Query("satker")
	metode := c.Query("metode")
	jenis := c.Query("jenis")
	return services.GetDataTablePaketSirup(c, tahun, satker, metode, jenis)
}

func GetJsonSwakelolaSirup(c *fiber.Ctx) error {
	return services.GetDataTableSwakelolaSirup(c)
}

func GetPaketSirup(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	paket := services.GetPaketSirup(uint(id))
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	mp["paket"] = paket
	mp["lokasi"] = paket.Lokasi()
	mp["metode"] = paket.Metode()
	mp["anggaran"] = paket.Anggaran()
	return c.Render("paket/rencana-detil", mp)
}

func GetSwakelolaSirup(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	paket := services.GetSwakelolaSirup(uint(id))
	if paket.ID == 0 {
		return c.SendStatus(404)
	}
	mp["paket"] = paket
	return c.Render("paket/rencana-swakelola-detil", mp)
}
