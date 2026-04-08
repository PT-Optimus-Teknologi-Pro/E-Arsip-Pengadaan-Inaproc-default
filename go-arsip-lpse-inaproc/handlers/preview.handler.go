package handlers

import (
	"arsip/config"
	"arsip/services"
	"arsip/utils"
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func PreviewImage(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	document := services.GetDocument(uint(id))
	c.Type("png")
	return c.SendFile(document.Filepath)
}

func PreviewTemplates(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	templates := services.GetTemplates(uint(id))
	if templates.ID == 0 {
		return c.SendStatus(404)
	}
	mp["templates"] = templates
	return c.Render("templates/templates-preview", mp)
}

func PreviewSkPp(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	log.Info("priview sk pp")
	mp := currentMap(c)
	paket := services.GetPaket(uint(id))
	mp["paket"] = paket
	mp["pegawai"] = paket.Pp()
	return c.Render("preview/surat-penunjukan-pp", mp)
}

func CetakSkPp(c *fiber.Ctx) error {
	log.Info("cetak sk pp")
	url := fmt.Sprintf("http://localhost:%s/preview/sk-pp/%s", config.Port(), c.Params("id"))
	return print(c, url, "SK-pejabat-pengadaan.pdf")
}

func PreviewSkPokja(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	log.Info("priview sk pokja")
	mp := currentMap(c)
	paket := services.GetPaket(uint(id))
	mp["paket"] = paket
	mp["pokja"] = paket.Pokja()
	return c.Render("preview/surat-penunjukan-pokja", mp)
}

func CetakSkPokja(c *fiber.Ctx) error {
	log.Info("cetak sk pokja")
	url := fmt.Sprintf("http://localhost:%s/preview/sk-pokja/%s",config.Port(), c.Params("id"))
	return print(c, url, "SK-pokja.pdf")
}

func PreviewBAKajiUlang(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	log.Info("priview BA Kaji Ulang")
	mp := currentMap(c)
	paket := services.GetPaket(uint(id))
	mp["paket"] = paket
	if mp["isPP"].(bool) {
		mp["pp"] = paket.Pp()
	}
	if mp["isPokja"].(bool) {
		mp["pokja"] = paket.Pokja()
	}
	return c.Render("preview/ba-kajiulang", mp)
}

func CetakBAKajiUlang(c *fiber.Ctx) error {
	log.Info("cetak BA Kaji ulang")
	url := fmt.Sprintf("http://localhost:%s/preview/ba-kajiulang/%s", config.Port(),c.Params("id"))
	return print(c, url, "BA-kajiulang.pdf")
}


func PreviewBANego(c *fiber.Ctx)  error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	return c.Render("preview/ba-nego", mp)
}


func CetakBANego(c *fiber.Ctx) error {
	log.Info("cetak BA Nego")
	url := fmt.Sprintf("http://localhost:%s/preview/ba-nego/%s", config.Port(), c.Params("id"))
	return print(c, url, "BA-Nego.pdf")
}


func PreviewBAPenetapan(c *fiber.Ctx)  error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	return c.Render("preview/ba-penetapan", mp)
}

func CetakBAPenetapan(c *fiber.Ctx) error {
	log.Info("cetak BA Penetapan")
	url := fmt.Sprintf("http://localhost:%s/preview/ba-penetapan/%s", config.Port(), c.Params("id"))
	return print(c, url, "BA-Penetapan-Pemenang.pdf")
}

func print(c *fiber.Ctx, url string, filename string) error {
	result := utils.ExportToPdf(url)
	reader := bytes.NewReader(result)
	// Set the Content-Type header for PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.SendStream(reader)
}
