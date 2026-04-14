package handlers

import (
	"arsip/config"
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

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

	// Get SK data from PejabatPengadaan Master
	ppSatker := models.GetPejabatPengadaanSatker(paket.SatkerId)
	sk := models.GetPejabatPengadaan(ppSatker.PpId)
	mp["sk"] = sk

	// Full URL for QR Code Verification
	mp["full_url"] = fmt.Sprintf("%s://%s/preview/sk-pp/%d", c.Protocol(), c.Hostname(), id)

	return c.Render("preview/surat-penunjukan-pp", mp)
}

func CetakSkPp(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Infof("cetak sk pp id: %s", id)

	// Persiapkan semua data
	paket := services.GetPaket(utils.StringToUint(id))
	pegawai := paket.Pp()
	ppSatker := models.GetPejabatPengadaanSatker(paket.SatkerId)
	sk := models.GetPejabatPengadaan(ppSatker.PpId)
	appSettings := services.GetSettings()
	satker := paket.Satker()

	fullUrl := fmt.Sprintf("%s://%s/preview/sk-pp/%s", c.Protocol(), c.Hostname(), id)

	// Build HTML langsung di Go (tidak pakai Django template engine untuk menghindari masalah layout)
	htmlContent := buildSkPpHtml(sk, pegawai, paket, satker, appSettings, fullUrl)

	result := utils.ExportHtmlToPdf(htmlContent, "")
	if result == nil || len(result) == 0 {
		log.Error("PDF generation returned empty/nil result")
		return flashError(c, "Gagal membuat PDF", "/paket/"+id)
	}

	reader := bytes.NewReader(result)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"SK-pejabat-pengadaan.pdf\"")
	return c.SendStream(reader)
}

func renderToString(templateName string, data fiber.Map) (string, error) {
	var buf bytes.Buffer
	binding := map[string]interface{}{}
	for k, v := range data {
		binding[k] = v
	}
	err := utils.Engine.Render(&buf, templateName, binding)
	if err != nil {
		log.Errorf("renderToString error for template '%s': %v", templateName, err)
	}
	return buf.String(), err
}

func buildSkPpHtml(sk models.PejabatPengadaan, pegawai models.Pegawai, paket models.Paket, satker models.SatkerSirup, settings models.AppSettings, fullUrl string) string {
	docInstansi := settings.DocInstansi
	if docInstansi == "" { docInstansi = "PEMERINTAH KOTA BANJARMASIN" }
	docSub := settings.DocSubInstansi
	if docSub == "" { docSub = "SEKRETARIAT DAERAH" }
	docAddr := settings.DocAddress
	if docAddr == "" { docAddr = "Jl. RE. Martadinata No.1 Banjarmasin" }
	docPejabatNama := settings.DocPejabatNama
	if docPejabatNama == "" { docPejabatNama = "AHSAN BUDIMAN" }
	docPejabatJabatan := settings.DocPejabatJabata
	if docPejabatJabatan == "" { docPejabatJabatan = "SEKRETARIS DAERAH" }

	tahunAng := sk.Tahun + 1

	// Generate QR Code sebagai base64 PNG (tanpa internet)
	qrPng, qrErr := qrcode.Encode(fullUrl, qrcode.Medium, 128)
	qrCell := `<div style="font-size:7pt; border:1px solid #ccc; padding:3px; text-align:center;">` +
		"<small>Scan Verifikasi</small></div>"
	if qrErr == nil {
		qrB64 := base64.StdEncoding.EncodeToString(qrPng)
		qrCell = fmt.Sprintf(
			`<img src="data:image/png;base64,%s" style="width:85px;"><div style="font-size:6pt;text-align:center;">Verifikasi Dokumen</div>`,
			qrB64,
		)
	}

	// Logo: baca dari disk, embed sebagai base64
	logoHtml := ""
	if settings.DocLogoPath != "" {
		cwd, _ := os.Getwd()
		// DocLogoPath berformat: /uploads/settings/filename.png
		// Disk path: {cwd}/public/uploads/settings/filename.png
		logoWebPath := strings.TrimPrefix(settings.DocLogoPath, "/")
		logoDiskPath := filepath.Join(cwd, "public", logoWebPath)

		logoData, err := os.ReadFile(logoDiskPath)
		if err == nil {
			ext := strings.ToLower(filepath.Ext(logoDiskPath))
			mime := "image/png"
			if ext == ".jpg" || ext == ".jpeg" {
				mime = "image/jpeg"
			} else if ext == ".svg" {
				mime = "image/svg+xml"
			} else if ext == ".webp" {
				mime = "image/webp"
			}
			b64 := base64.StdEncoding.EncodeToString(logoData)
			logoHtml = fmt.Sprintf(`<img src="data:%s;base64,%s" alt="Logo" style="width:80px;">`, mime, b64)
		} else {
			log.Warnf("Logo tidak ditemukan di: %s", logoDiskPath)
		}
	}

	golStr := ""
	if pegawai.PegGolongan != "" {
		golStr = fmt.Sprintf("(%s)", pegawai.PegGolongan)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { font-family: 'Times New Roman', Times, serif; color: black; font-size: 12pt; }
  table { width: 100%%; border-collapse: collapse; }
  ol { margin-top: 0; padding-left: 20px; }
  li { text-align: justify; margin-bottom: 5px; }
</style>
</head>
<body style="margin: 40px 50px;">

  <!-- Kop Surat -->
  <table style="border-bottom: 3px solid black; padding-bottom: 10px; margin-bottom: 20px;">
    <tr>
      <td style="width:15%%; text-align:center;">%s</td>
      <td style="width:70%%; text-align:center;">
        <div style="font-size:14pt; font-weight:bold; text-transform:uppercase;">%s</div>
        <div style="font-size:16pt; font-weight:bold; text-transform:uppercase;">%s</div>
        <div style="font-size:9pt;">%s %s</div>
        <div style="font-size:9pt;">%s</div>
      </td>
      <td style="width:15%%; text-align:right; vertical-align:top;">
        %s
      </td>
    </tr>
  </table>

  <!-- Judul -->
  <div style="text-align:center; margin-bottom:20px;">
    <div style="font-weight:bold; text-transform:uppercase;">KEPUTUSAN %s KOTA BANJARMASIN</div>
    <div style="font-weight:bold; text-transform:uppercase;">NOMOR %s TAHUN %d</div>
    <div style="margin:5px 0; font-weight:bold;">TENTANG</div>
    <div style="font-weight:bold; text-transform:uppercase;">PENETAPAN PEJABAT PENGADAAN BARANG/JASA</div>
    <div style="font-weight:bold; text-transform:uppercase;">PADA %s</div>
    <div style="font-weight:bold; text-transform:uppercase;">%s</div>
    <div style="font-weight:bold; text-transform:uppercase;">TAHUN ANGGARAN %d</div>
    <div style="margin-top:10px; font-weight:bold; text-transform:uppercase;">%s KOTA BANJARMASIN</div>
  </div>

  <!-- Menimbang / Mengingat -->
  <table style="margin-bottom:15px;">
    <tr>
      <td style="width:20%%; vertical-align:top; font-weight:bold;">Menimbang</td>
      <td style="width:3%%; vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        <ol type="a">
          <li>Bahwa untuk tertib dan lancarnya pelaksanaan pengadaan barang/jasa di %s %s Tahun %d dipandang perlu menetapkan Pejabat Pengadaan Barang/Jasa.</li>
          <li>bahwa berdasarkan pertimbangan sebagaimana dimaksud dalam huruf a, perlu menetapkan Keputusan %s tentang Penetapan Pejabat Pengadaan Barang/Jasa pada %s %s Tahun Anggaran %d;</li>
        </ol>
      </td>
    </tr>
    <tr><td colspan="3" style="height:10px;"></td></tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">Mengingat</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        <ol>
          <li>Undang-Undang Nomor 27 Tahun 1959 tentang Penetapan Undang-Undang Darurat Nomor 3 Tahun 1953 tentang Perpanjangan Pembentukan Daerah Tingkat II di Kalimantan;</li>
          <li>Undang-Undang Nomor 23 Tahun 2014 tentang Pemerintahan Daerah sebagaimana telah diubah dengan Undang-Undang Nomor 6 Tahun 2023;</li>
          <li>Undang-Undang Nomor 1 Tahun 2022 tentang Hubungan Keuangan Antara Pemerintah Pusat dan Pemerintahan Daerah;</li>
          <li>Undang-Undang Nomor 20 Tahun 2023 tentang Aparatur Sipil Negara;</li>
          <li>Peraturan Pemerintah Nomor 12 Tahun 2019 tentang Pengelolaan Keuangan Daerah;</li>
          <li>Peraturan Presiden Nomor 16 Tahun 2018 tentang Pengadaan Barang/Jasa Pemerintah sebagaimana telah diubah dengan Peraturan Presiden Nomor 12 Tahun 2021;</li>
          <li>Peraturan Menteri Dalam Negeri Nomor 77 Tahun 2020 tentang Pedoman Teknis Pengelolaan Keuangan Daerah;</li>
          <li>Peraturan Lembaga Kebijakan Pengadaan Barang/Jasa Pemerintah Nomor 11 Tahun 2021 tentang Pedoman Perencanaan Pengadaan Barang/Jasa Pemerintah;</li>
          <li>Peraturan Daerah Kota Banjarmasin Nomor 7 Tahun 2016 tentang Pembentukan dan Susunan Perangkat Daerah Kota Banjarmasin sebagaimana telah diubah dengan Peraturan Daerah Kota Banjarmasin Nomor 3 Tahun 2021;</li>
        </ol>
      </td>
    </tr>
    <tr><td colspan="3" style="height:10px;"></td></tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">Memperhatikan</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Dokumen Pelaksanaan Anggaran Satuan Kerja Perangkat Daerah %s %s Tahun Anggaran %d</td>
    </tr>
  </table>

  <div style="text-align:center; margin-bottom:15px; font-weight:bold; text-transform:uppercase;">MEMUTUSKAN:</div>

  <table style="margin-bottom:20px;">
    <tr>
      <td style="width:20%%; vertical-align:top; font-weight:bold;">MENETAPKAN</td>
      <td style="width:3%%; vertical-align:top;">:</td>
      <td style="vertical-align:top;"></td>
    </tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">KESATU</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        Pejabat Pengadaan Barang/Jasa Pada %s %s Tahun Anggaran %d sebagai berikut:<br>
        <table style="margin-top:10px; border-collapse:collapse;">
          <tr><td style="width:100px;">Nama</td><td>: %s</td></tr>
          <tr><td>NIP</td><td>: %s</td></tr>
          <tr><td>Pangkat/Gol</td><td>: %s %s</td></tr>
          <tr><td>Unit Kerja</td><td>: %s</td></tr>
        </table>
      </td>
    </tr>
    <tr><td colspan="3" style="height:10px;"></td></tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">KEDUA</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        Tugas Pejabat Pengadaan Barang/Jasa sebagaimana dimaksud dalam Diktum KESATU adalah melaksanakan persiapan dan pelaksanaan pengadaan langsung, penunjukan langsung, dan e-purchasing sesuai ketentuan yang berlaku; menyerahkan dokumen hasil pemilihan kepada PA/KPA; serta memberikan pertanggungjawaban kepada PA/KPA.
      </td>
    </tr>
    <tr><td colspan="3" style="height:10px;"></td></tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">KETIGA</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Segala biaya yang timbul akibat ditetapkannya Keputusan %s ini dibebankan pada APBD Banjarmasin Tahun Anggaran %d pada %s %s.</td>
    </tr>
    <tr><td colspan="3" style="height:10px;"></td></tr>
    <tr>
      <td style="vertical-align:top; font-weight:bold;">KEEMPAT</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Keputusan %s ini mulai berlaku pada tanggal ditetapkan.</td>
    </tr>
  </table>

  <!-- Tanda Tangan -->
  <div style="float:right; width:40%%; margin-top:30px; text-align:left;">
    <p style="margin:0;">Ditetapkan di Banjarmasin</p>
    <p style="margin:0;">pada tanggal 9 Desember 2024</p>
    <p style="margin-top:10px; font-weight:bold; text-transform:uppercase;">%s</p>
    <p style="margin:0; font-weight:bold; text-transform:uppercase;">%s</p>
    <br><br><br><br>
    <p style="margin:0; font-weight:bold; text-decoration:underline; text-transform:uppercase;">%s</p>
  </div>
  <div style="clear:both;"></div>

</body>
</html>`,
		// Kop: logo, instansi, sub, addr, telp+fax, website
		logoHtml, docInstansi, docSub, docAddr,
		func() string {
			s := ""
			if settings.DocPhone != "" { s += "("+settings.DocPhone+") " }
			if settings.DocFax != "" { s += settings.DocFax }
			return s
		}(),
		settings.DocWebsite,
		// Verifikasi box (kanan atas)
		qrCell,
		// Judul
		docPejabatJabatan, sk.NoSk, sk.Tahun,
		docSub, docInstansi, tahunAng, docPejabatJabatan,
		// Menimbang
		docSub, docInstansi, tahunAng,
		docPejabatJabatan, docSub, docInstansi, tahunAng,
		// Memperhatikan
		docSub, docInstansi, tahunAng,
		// Kesatu
		docSub, docInstansi, tahunAng,
		pegawai.PegNama, pegawai.PegNip, pegawai.PegPangkat, golStr, satker.Nama,
		// Ketiga
		docPejabatJabatan, tahunAng, docSub, docInstansi,
		// Keempat
		docPejabatJabatan,
		// TTD
		docPejabatJabatan, docInstansi, docPejabatNama,
	)
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
