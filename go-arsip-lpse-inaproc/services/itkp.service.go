package services

import (
	"arsip/cache"
	"arsip/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func UpdateItkp(tahun int) {
	key := "itkp-"+strconv.Itoa(tahun)
	itkp, found := cache.Get(key)
	if !found {
		models.UpdateItkpRup(tahun)
		itkp = "done"
		cache.Set(key, itkp)
	}
}

func GetItkp(tahun int) []models.Itkp {
	return models.GetItkp(tahun)
}

func GetDataTableItkp(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkp(c, tahun)
}

func GetDataTableItkpRup(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkpRup(c, tahun)
}

func GetDataTableItkpPemilihan(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkpPemilihan(c, tahun)
}

func GetDataTableItkpTender(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkpTender(c, tahun)
}

func GetDataTableItkpPurchase(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkpPurchase(c, tahun)
}

func GetDataTableItkpNontender(c *fiber.Ctx, tahun int) error {
	return models.GetDataTableItkpNontender(c, tahun)
}

func GetDetilRupSatker(satkerId uint, satkerStr string, tahun int) []models.ItkpRupSatker {
	return models.GetDetilRupSatker(satkerId, satkerStr, tahun)
}

func GetDetilTenderSatker(satkerId string, satkerStr string, tahun int) []models.ItkpTenderSatker {
	return models.GetDetilTenderSatker(satkerId, satkerStr, tahun)
}

func GetDetilNontenderSatker(satkerId string, satkerStr string, tahun int) []models.ItkpNontenderSatker {
	return models.GetDetilNontenderSatker(satkerId, satkerStr, tahun)
}

func GetDetilPurchaseSatker(satkerId uint, tahun int) []models.ItkpPurchaseSatker {
	return models.GetDetilPurchaseSatker(satkerId, tahun)
}

func GetDetilEkontrakSatker(satkerId string, satkerStr string, tahun int) []models.ItkpEkontrakSatker {
	return models.GetDetilEkontrakSatker(satkerId, satkerStr, tahun)
}
