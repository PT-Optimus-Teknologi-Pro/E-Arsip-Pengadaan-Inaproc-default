package handlers

import (
	"arsip/config"
	"arsip/services"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Itkp(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	log.Info(tahun)
	mp["tahun"] =tahun
	services.UpdateItkp(tahun)
	mp["tahunlist"] = config.GetTahunList()
	return c.Render("itkp/itkp", mp)
}

func GetJsonItkpLokal(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkp(c, tahun)
}

func GetJsonItkpRup(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkpRup(c, tahun)
}

func GetItkpRupDetail(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.QueryInt("satker")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["paket"] = services.GetDetilRupSatker(uint(satker), objSatker.IdSatker, tahun)
	return c.Render("itkp/itkp-rup", mp)
}

func GetJsonItkpPemilihan(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkpPemilihan(c, tahun)
}

func GetItkpPemilihanDetil(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.QueryInt("satker")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["paket"] = services.GetDetilEkontrakSatker(c.Query("satker"), objSatker.IdSatker, tahun)
	return c.Render("itkp/itkp-pemilihan", mp)
}

func GetJsonItkpTender(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkpTender(c, tahun)
}

func GetItkpTenderDetil(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.QueryInt("satker")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["paket"] = services.GetDetilTenderSatker(c.Query("satker"), objSatker.IdSatker, tahun)
	return c.Render("itkp/itkp-tender", mp)
}

func GetJsonItkpPurchase(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkpPurchase(c, tahun)
}

func GetItkpPurchaseDetil(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.QueryInt("satker")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["paket"] = services.GetDetilPurchaseSatker(uint(satker), tahun)
	return c.Render("itkp/itkp-purchase", mp)
}


func GetJsonItkpNonTender(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	return services.GetDataTableItkpNontender(c, tahun)
}

func GetItkpNonTenderDetil(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun", time.Now().Year())
	satker := c.QueryInt("satker")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["paket"] = services.GetDetilNontenderSatker(c.Query("satker"), objSatker.IdSatker, tahun)
	return c.Render("itkp/itkp-nontender", mp)
}
