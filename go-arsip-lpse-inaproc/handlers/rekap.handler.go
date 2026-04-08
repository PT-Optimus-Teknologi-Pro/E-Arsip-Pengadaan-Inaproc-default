package handlers

import (
	"arsip/services"
	"github.com/gofiber/fiber/v2"
)


func RekapRencanaSatker(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun")
	satker := uint(c.QueryInt("satker"))
	bulan := c.QueryInt("bulan")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	paket := services.GetRekapRencanaPaketSatkerBulan(tahun, satker, bulan)
	if len(paket) > 0 {
		mp["paket"] = paket
	}
	mp["tahun"] = tahun
	mp["bulan"] = bulan
	return c.Render("beranda/rekap-paket-rencana-skpd", mp)
}


func RekapRealisasiSatker(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun")
	satker := uint(c.QueryInt("satker"))
	bulan := c.QueryInt("bulan")
	objSatker := services.GetSatkerSirup(uint(satker))
	mp["namaSatker"] = objSatker.Nama
	mp["tahun"] = tahun
	mp["bulan"] = bulan
	mp["tender"] = services.GetRekapRealisasiPaketTenderSatkerBulan(tahun, satker, bulan)
	mp["nontender"] = services.GetRekapRealisasiPaketNontenderSatkerBulan(tahun, satker, bulan)
	mp["pencatatan"] = services.GetRekapRealisasiPencatatanSatkerBulan(tahun, satker, bulan)
	mp["purchase"] = services.GetRekapRealisasiPurchaseSatkerBulan(tahun, satker, bulan)
	return c.Render("beranda/rekap-paket-realisasi-skpd", mp)
}


func RekapPaketPPK(c *fiber.Ctx) error {
	mp := currentMap(c)
	tahun := c.QueryInt("tahun")
	nip := c.Query("nip")
	nama := c.QueryInt("nama")
	mp["nama"] = nama
	mp["nip"] = nip
	mp["tahun"] = tahun
	mp["tender"] = services.GetRekapTenderPPk(nip, tahun)
	mp["nontender"] = services.GetRekapNontenderPPk(nip, tahun)
	mp["pencatatan"] = services.GetRekapPencacatanPPk(nip, tahun)
	mp["purchase"] = services.GetRekapPurchasePPK(nip, tahun)
	return c.Render("beranda/rekap-paket-ppk-detil", mp)
}
