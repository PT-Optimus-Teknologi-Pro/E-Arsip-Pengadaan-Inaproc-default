package handlers

import (
	"arsip/cache"
	"arsip/models"
	"arsip/services"
	"strconv"

	// "encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetJsonProgressRup(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "progressRup-"+strconv.Itoa(tahun)
	progressRup, found := cache.Get(key)
	if !found {
		progressRup = services.GetRupProgress(tahun)
		cache.Set(key, progressRup)
	}
	responseData := fiber.Map{"data" : []interface{}{}}
	if len(progressRup.([]models.RupProgress)) > 0 {
		responseData = fiber.Map{"data" : progressRup}
	}
	return c.JSON(responseData)
}

func GetJsonPaketPrioritas(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "paketPrioritas-"+strconv.Itoa(tahun)
	paketPrioritas, found := cache.Get(key)
	if !found {
		paketPrioritas = services.GetPaketPrioritas(tahun)
		cache.Set(key, paketPrioritas)
	}
	responseData := fiber.Map{"data" : []interface{}{}}
	if len(paketPrioritas.([]models.PaketPrioritas)) > 0 {
		responseData = fiber.Map{"data" : paketPrioritas}
	}
	return c.JSON(responseData)
}

func GetJsonRekapSatker(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "rekapSatker-"+strconv.Itoa(tahun)
	rekapSatker, found := cache.Get(key)
	if !found {
		rekapSatker = services.CalculateRekapSatker(tahun)
		cache.Set(key, rekapSatker)
	}
	responseData := fiber.Map{"data" : []interface{}{}}
	if len(rekapSatker.([]models.RekapSatkerDashboard)) > 0 {
		responseData = fiber.Map{"data" : rekapSatker}
	}
	return c.JSON(responseData)
}

func GetJsonBebanPersonel(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "bebanPersonil-"+strconv.Itoa(tahun)
	bebanPersonil, found := cache.Get(key)
	if !found {
		bebanPersonil = services.GetBebanPersonil(tahun)
		cache.Set(key, bebanPersonil)
	}
	responseData := fiber.Map{
		"data" : bebanPersonil,
	}
	return c.JSON(responseData)
}

func GetJsonRekapPaketPerSatker(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "rekapPerSatker-"+strconv.Itoa(tahun)
	rekapPerSatker, found := cache.Get(key)
	if !found {
		rekapPerSatker = services.GetRekapPaketSatker(tahun)
		if len(rekapPerSatker.([]models.RekapPaketSatker)) > 0 {
			cache.Set(key, rekapPerSatker)
		}
	}
	responseData := fiber.Map{"data" : []interface{}{}}
	if len(rekapPerSatker.([]models.RekapPaketSatker)) > 0 {
		responseData = fiber.Map{"data" : rekapPerSatker}
	}
	return c.JSON(responseData)
}

func GetJsonRekapPaketPPK(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	key := "rekapPpk-"+strconv.Itoa(tahun)
	rekapPpk, found := cache.Get(key)
	if !found {
		rekapPpk = services.GetRekapPaketPpk(tahun)
		cache.Set(key, rekapPpk)
	}
	responseData := fiber.Map{"data" : rekapPpk}
	return c.JSON(responseData)
}
