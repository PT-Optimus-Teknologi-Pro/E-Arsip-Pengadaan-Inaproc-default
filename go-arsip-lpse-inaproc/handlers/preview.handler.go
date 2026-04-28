package handlers

import (
	"arsip/models"
	"arsip/services"
	"arsip/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	hash := c.Query("hash")
	if hash != "" {
		var dok models.DokumenTercetak
		models.GetDB().Where("md5_hash = ?", hash).First(&dok)
		mp["dokTercetak"] = dok

		valUrl := fmt.Sprintf("%s://%s/validasi/dokumen/%s", c.Protocol(), c.Hostname(), hash)
		mp["qrValidasi"] = generateQrBase64(valUrl)
        if !dok.TanggalPenetapan.IsZero() {
            mp["tglPenetapan"] = dok.TanggalPenetapan.Format("02-01-2006")
        }
	}

	appSettings := services.GetSettings()
	mp["appSettings"] = appSettings

	return c.Render("preview/surat-penunjukan-pp", mp)
}

func CetakSkPp(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	
	// Check if already printed, if not create record with defaults
	var dok models.DokumenTercetak
	models.GetDB().Where("paket_id = ? AND jenis_dokumen = ?", id, "SK_PENUNJUKAN_PP").First(&dok)
	if dok.ID == 0 {
		paket := services.GetPaket(id)
		ppSatker := models.GetPejabatPengadaanSatker(paket.SatkerId)
		skMaster := models.GetPejabatPengadaan(ppSatker.PpId)
		
		dok = models.DokumenTercetak{
			PaketID:      id,
			JenisDokumen: "SK_PENUNJUKAN_PP",
			NomorSurat:   skMaster.NoSk,
			TahunSurat:   fmt.Sprintf("%d", skMaster.Tahun),
			TempatPenetapan: skMaster.TempatSk,
			TanggalPenetapan: skMaster.TglSk,
		}
		session := getUserSession(c)
		if session.Pegawai().ID > 0 { dok.PembuatPegawaiID = session.Pegawai().ID }
		models.GetDB().Create(&dok)
	}

	paket := services.GetPaket(id)
	pegawai := paket.Pp()
	ppSatker := models.GetPejabatPengadaanSatker(paket.SatkerId)
	sk := models.GetPejabatPengadaan(ppSatker.PpId)
	appSettings := services.GetSettings()
	satker := paket.Satker()
	fullUrl := fmt.Sprintf("%s://%s/preview/sk-pp/%d?hash=%s", c.Protocol(), c.Hostname(), id, dok.Md5Hash)

	htmlContent := buildSkPpHtml(sk, pegawai, satker, appSettings, dok, fullUrl)
	result := utils.ExportHtmlToPdf(htmlContent, "")
	reader := bytes.NewReader(result)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"SK-pejabat-pengadaan.pdf\"")
	return c.SendStream(reader)
}

func CetakSkPpProcess(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	dok := processFormDokumen(c, id, "SK_PENUNJUKAN_PP")

	paket := services.GetPaket(id)
	pegawai := paket.Pp()
	ppSatker := models.GetPejabatPengadaanSatker(paket.SatkerId)
	sk := models.GetPejabatPengadaan(ppSatker.PpId)
	appSettings := services.GetSettings()
	satker := paket.Satker()
	fullUrl := fmt.Sprintf("%s://%s/preview/sk-pp/%d?hash=%s", c.Protocol(), c.Hostname(), id, dok.Md5Hash)

	htmlContent := buildSkPpHtml(sk, pegawai, satker, appSettings, dok, fullUrl)
	result := utils.ExportHtmlToPdf(htmlContent, "")
	reader := bytes.NewReader(result)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"SK-pejabat-pengadaan.pdf\"")
	return c.SendStream(reader)
}

func processFormDokumen(c *fiber.Ctx, paketID uint, jenis string) models.DokumenTercetak {
	dok := models.DokumenTercetak{
		PaketID:      paketID,
		JenisDokumen: jenis,
		NomorSurat:   c.FormValue("nomor_surat"),
		TentangSurat: c.FormValue("tentang_surat"),
		TahunSurat:   c.FormValue("tahun_surat"),
		TempatPenetapan: c.FormValue("tempat_penetapan"),
		NomorKeputusanSekda: c.FormValue("nomor_kep_sekda"),
	}

	tglPenetapan := c.FormValue("tanggal_penetapan")
	if tglPenetapan != "" { pt, _ := time.Parse("2006-01-02", tglPenetapan); dok.TanggalPenetapan = pt }
	tglTerbit := c.FormValue("tanggal_terbit_kep")
	if tglTerbit != "" { tt, _ := time.Parse("2006-01-02", tglTerbit); dok.TanggalTerbitKeputusan = tt }

	session := getUserSession(c)
	if session.Pegawai().ID > 0 { dok.PembuatPegawaiID = session.Pegawai().ID }
	models.GetDB().Create(&dok)
	return dok
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

func buildSkPpHtml(sk models.PejabatPengadaan, pegawai models.Pegawai, satker models.SatkerSirup, settings models.AppSettings, dokTercetak models.DokumenTercetak, fullUrl string) string {
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
	noSurat := sk.NoSk
	if dokTercetak.NomorSurat != "" {
		noSurat = dokTercetak.NomorSurat
	}
	tahunSurat := sk.Tahun
	if dokTercetak.TahunSurat != "" {
		tahunSurat = int(utils.StringToUint(dokTercetak.TahunSurat))
	}
	
	tempatPenetapan := "Banjarmasin"
	if sk.TempatSk != "" { tempatPenetapan = sk.TempatSk }
	if dokTercetak.TempatPenetapan != "" { tempatPenetapan = dokTercetak.TempatPenetapan }
	
	tglPenetapan := "..."
	if !sk.TglSk.IsZero() { tglPenetapan = sk.TglSk.Format("02-01-2006") }
	if !dokTercetak.TanggalPenetapan.IsZero() { tglPenetapan = dokTercetak.TanggalPenetapan.Format("02-01-2006") }

	// Generate QR Code sebagai base64 PNG (tanpa internet)
	qrPng, qrErr := qrcode.Encode(fullUrl, qrcode.Medium, 128)
	qrCell := `<div style="font-size:7pt; border:1px solid #ccc; padding:3px; text-align:center;">` +
		"<small>Scan Verifikasi</small></div>"
	if qrErr == nil && dokTercetak.Md5Hash != "" {
		qrB64 := base64.StdEncoding.EncodeToString(qrPng)
		qrCell = fmt.Sprintf(
			`<img src="data:image/png;base64,%s" style="width:85px;"><div style="font-size:6pt;text-align:center;">Kode Verifikasi: %s</div>`,
			qrB64, dokTercetak.Md5Hash[:8],
		)
	}

	// Logo: baca dari disk, embed sebagai base64
	logoHtml := ""
	if settings.DocLogoPath != "" {
		cwd, _ := os.Getwd()
		logoWebPath := strings.TrimPrefix(settings.DocLogoPath, "/")
		// Coba beberapa variasi path untuk Windows
		pathsToTry := []string{
			filepath.Join(cwd, "public", logoWebPath),
			filepath.Join(cwd, logoWebPath),
			filepath.Join("public", logoWebPath),
		}

		var logoData []byte
		var err error
		var finalPath string

		for _, p := range pathsToTry {
			finalPath = p
			logoData, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}

		if err == nil {
			log.Infof("Logo ditemukan dan dibaca dari: %s", finalPath)
			ext := strings.ToLower(filepath.Ext(finalPath))
			mime := "image/png"
			switch ext {
			case ".jpg", ".jpeg":
				mime = "image/jpeg"
			case ".svg":
				mime = "image/svg+xml"
			case ".webp":
				mime = "image/webp"
			}
			b64 := base64.StdEncoding.EncodeToString(logoData)
			logoHtml = fmt.Sprintf(`<img src="data:%s;base64,%s" alt="Logo" style="width:80px; max-height:80px;">`, mime, b64)
		} else {
			log.Warnf("Logo TIDAK ditemukan. Percobaan terakhir di: %s. Error: %v", finalPath, err)
		}
	}

	docSignatureHtml := ""
	if (settings.DocSignatureMode == "digital" || settings.DocSignatureMode == "canvas") && settings.DocSignaturePath != "" {
		cwd, _ := os.Getwd()
		sigWebPath := strings.TrimPrefix(settings.DocSignaturePath, "/")
		pathsToTry := []string{
			filepath.Join(cwd, "public", sigWebPath),
			filepath.Join(cwd, sigWebPath),
			filepath.Join("public", sigWebPath),
		}
		for _, p := range pathsToTry {
			sigData, err := os.ReadFile(p)
			if err == nil {
				ext := strings.ToLower(filepath.Ext(p))
				mime := "image/png"
				if ext == ".webp" { mime = "image/webp" }
				b64 := base64.StdEncoding.EncodeToString(sigData)
				docSignatureHtml = fmt.Sprintf(`<img src="data:%s;base64,%s" style="max-height:60px; margin-top:10px;">`, mime, b64)
				break
			}
		}
	}
	if docSignatureHtml == "" {
		docSignatureHtml = "<br><br><br><br>" // Fallback if no signature
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
        <div style="font-size:9pt;">Website: %s | Email: %s</div>
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
    <p style="margin:0;">Ditetapkan di %s</p>
    <p style="margin:0;">pada tanggal %s</p>
    <p style="margin-top:10px; font-weight:bold; text-transform:uppercase;">%s</p>
    <p style="margin:0; font-weight:bold; text-transform:uppercase;">%s</p>
    %s
    <p style="margin:0; font-weight:bold; text-decoration:underline; text-transform:uppercase;">%s</p>
    <p style="margin:0;">NIP. %s</p>
  </div>
  <div style="clear:both;"></div>

</body>
</html>`,
		// Kop: logo, instansi, sub, addr, telp+fax, website
		logoHtml, docInstansi, docSub, docAddr,
		func() string {
			s := ""
			if settings.DocPhone != "" { s += "Telp: "+settings.DocPhone+" " }
			if settings.DocFax != "" { s += "| Fax: "+settings.DocFax }
			return s
		}(),
		settings.DocWebsite, settings.DocEmail,
		// Verifikasi box (kanan atas)
		qrCell,
		// Judul
		docPejabatJabatan, noSurat, tahunSurat,
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
		tempatPenetapan,
		tglPenetapan,
		docPejabatJabatan, func() string { if settings.DocRegion != "" { return settings.DocRegion } else { return "KOTA BANJARMASIN" } }(),
		docSignatureHtml,
		docPejabatNama, settings.DocPejabatNip,
	)
}

func PreviewSkPokja(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	mp := currentMap(c)
	paket := services.GetPaket(uint(id))
	mp["paket"] = paket
	mp["pokja"] = paket.Pokja()
	mp["appSettings"] = services.GetSettings()

	hash := c.Query("hash")
	if hash != "" {
		var dok models.DokumenTercetak
		models.GetDB().Where("md5_hash = ?", hash).First(&dok)
		mp["dokTercetak"] = dok
		valUrl := fmt.Sprintf("%s://%s/validasi/dokumen/%s", c.Protocol(), c.Hostname(), hash)
		mp["qrValidasi"] = generateQrBase64(valUrl)
        if !dok.TanggalPenetapan.IsZero() {
            mp["tglPenetapan"] = dok.TanggalPenetapan.Format("02-01-2006")
        }
	}
	return c.Render("preview/surat-penunjukan-pokja", mp)
}

func buildSkPokjaHtml(panitia models.Panitia, settings models.AppSettings, dokTercetak models.DokumenTercetak, fullUrl string) string {
	docInstansi := settings.DocInstansi
	if docInstansi == "" { docInstansi = "PEMERINTAH KOTA BANJARMASIN" }
	docSub := settings.DocSubInstansi
	if docSub == "" { docSub = "S E K R E T A R I A T  D A E R A H" }
	docAddr := settings.DocAddress
	if docAddr == "" { docAddr = "J. RE. Martadinata No.1 Banjarmasin 70111" }
	docRegion := settings.DocRegion
	if docRegion == "" { docRegion = "KOTA BANJARMASIN" }
	docPejabatNama := settings.DocPejabatNama
	if docPejabatNama == "" { docPejabatNama = "AHSAN BUDIMAN" }
	docPejabatJabatan := settings.DocPejabatJabata
	if docPejabatJabatan == "" { docPejabatJabatan = "SEKRETARIS DAERAH" }

	noSurat := panitia.NoSk
	if dokTercetak.NomorSurat != "" { noSurat = dokTercetak.NomorSurat }
	if noSurat == "" { noSurat = "................" }

	tahunSurat := panitia.Tahun
	if dokTercetak.TahunSurat != "" { tahunSurat = int(utils.StringToUint(dokTercetak.TahunSurat)) }

	tempatPenetapan := "Banjarmasin"
	if panitia.TempatSk != "" { tempatPenetapan = panitia.TempatSk }
	if dokTercetak.TempatPenetapan != "" { tempatPenetapan = dokTercetak.TempatPenetapan }

	tglPenetapan := "..."
	if !panitia.TglSk.IsZero() { tglPenetapan = panitia.TglSk.Format("02 January 2006") }
	if !dokTercetak.TanggalPenetapan.IsZero() { tglPenetapan = dokTercetak.TanggalPenetapan.Format("02 January 2006") }

	// QR & Logo
	qrPng, _ := qrcode.Encode(fullUrl, qrcode.Medium, 128)
	qrCell := ""
	if dokTercetak.Md5Hash != "" {
		qrB64 := base64.StdEncoding.EncodeToString(qrPng)
		// Menampilkan QR Code dengan label MD5 di bawahnya
		qrCell = fmt.Sprintf(`<div style="text-align:center;"><img src="data:image/png;base64,%s" style="width:80px;"><div style="font-size:7pt; font-family:monospace; font-weight:bold; margin-top:2px;">KODE: %s</div><div style="font-size:6pt;">Verifikasi Dokumen</div></div>`, qrB64, dokTercetak.Md5Hash[:8])
	}

	logoHtml := ""
	pathLogo := settings.DocLogoPath
	if pathLogo == "" { pathLogo = settings.LogoPath } // Fallback ke logo utama jika logo dokumen kosong

	if pathLogo != "" {
		cwd, _ := os.Getwd()
		p := filepath.Join(cwd, "public", strings.TrimPrefix(pathLogo, "/"))
		if logoData, err := os.ReadFile(p); err == nil {
			b64 := base64.StdEncoding.EncodeToString(logoData)
			logoHtml = fmt.Sprintf(`<img src="data:image/png;base64,%s" style="width:80px;">`, b64)
		}
	}

	// Logic Tanda Tangan Dinamis
	docSignatureHtml := ""
	mode := strings.ToLower(settings.DocSignatureMode)
	if mode == "digital" || mode == "" {
		// Jika pilih Barcode MD5 Otomatis
		if dokTercetak.Md5Hash != "" {
			qrPng, _ := qrcode.Encode(fullUrl, qrcode.Medium, 128)
			qrB64 := base64.StdEncoding.EncodeToString(qrPng)
			docSignatureHtml = fmt.Sprintf(`
				<div style="margin: 10px 0;">
				<div style="margin: 5px 0;">
					<img src="data:image/png;base64,%s" style="width:75px; height:75px;">
					<div style="font-size:7pt; font-family:monospace; font-weight:bold; margin-top:1px;">KODE: %s</div>
				</div>`, qrB64, dokTercetak.Md5Hash[:8])
		}
	} else if (mode == "canvas" || mode == "physical") && settings.DocSignaturePath != "" {
		cwd, _ := os.Getwd()
		p := filepath.Join(cwd, "public", strings.TrimPrefix(settings.DocSignaturePath, "/"))
		if sigData, err := os.ReadFile(p); err == nil {
			b64 := base64.StdEncoding.EncodeToString(sigData)
			docSignatureHtml = fmt.Sprintf(`<img src="data:image/png;base64,%s" style="max-height:90px; max-width:150px; margin: 5px 0;">`, b64)
		}
	}

	if docSignatureHtml == "" {
		docSignatureHtml = "<br><br><br><br><br>"
	}

	// Kop surat address line
	kopContact := ""
	if settings.DocPhone != "" { kopContact += settings.DocPhone }
	kopContact2 := ""
	if settings.DocFax != "" { kopContact2 += "Fax. " + settings.DocFax }
	if settings.DocWebsite != "" { kopContact2 += " " + settings.DocWebsite }

	// Build KESATU anggota list
	anggotaRows := ""
	for _, v := range panitia.AnggotaList() {
		golStr := v.PegGolongan
		if v.PegPangkat != "" && golStr != "" {
			golStr = v.PegPangkat + " (" + golStr + ")"
		} else if v.PegPangkat != "" {
			golStr = v.PegPangkat
		}
		unitKerja := "-"
		if v.AgcId.Valid {
			var agency models.Agency
			if err := models.GetDB().First(&agency, v.AgcId.Int64).Error; err == nil {
				unitKerja = agency.AgcNama
			}
		}
		if unitKerja == "-" && v.PegJabatan != "" {
			unitKerja = v.PegJabatan
		}
		
		anggotaRows += fmt.Sprintf(`
		<table style="margin-bottom:8px; font-size:11pt;">
		  <tr><td style="width:120px;">Nama</td><td style="width:10px;">:</td><td>%s</td></tr>
		  <tr><td>NIP</td><td>:</td><td>%s</td></tr>
		  <tr><td>Pangkat/Gol</td><td>:</td><td>%s</td></tr>
		  <tr><td>Unit Kerja</td><td>:</td><td>%s</td></tr>
		</table>`, v.PegNama, v.PegNip, golStr, unitKerja)
	}
	if anggotaRows == "" {
		anggotaRows = `<p style="font-style:italic; color:#666;">(Data anggota belum tersedia)</p>`
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { font-family: 'Times New Roman', Times, serif; color: black; font-size: 12pt; line-height: 1.4; margin: 40px 50px; }
  table { border-collapse: collapse; }
  table.full { width: 100%%; }
  ol { margin: 3px 0; padding-left: 20px; }
  li { text-align: justify; margin-bottom: 4px; }
  .bold { font-weight: bold; }
  .center { text-align: center; }
  .upper { text-transform: uppercase; }
  .underline { text-decoration: underline; }
  .small { font-size: 9pt; }
</style>
</head>
<body>
  <!-- KOP SURAT -->
  <table class="full" style="border-bottom: 3px solid black; padding-bottom: 8px; margin-bottom: 18px;">
    <tr>
      <td style="width:12%%; text-align:center; vertical-align:middle;">%s</td>
      <td style="width:76%%; text-align:center; vertical-align:middle;">
        <div class="bold upper" style="font-size:13pt;">%s</div>
        <div class="bold upper" style="font-size:15pt; letter-spacing:2px;">%s</div>
        <div class="small">%s %s</div>
        <div class="small">%s</div>
      </td>
      <td style="width:12%%; text-align:right; vertical-align:top;">%s</td>
    </tr>
  </table>

  <!-- JUDUL -->
  <div class="center" style="margin-bottom:16px;">
    <div class="bold upper">KEPUTUSAN %s %s</div>
    <div class="bold upper">NOMOR %s TAHUN %d</div>
    <div class="bold" style="margin:4px 0;">TENTANG</div>
    <div class="bold upper">PENETAPAN POKJA PENGADAAN BARANG/JASA</div>
    <div class="bold upper">PADA %s</div>
    <div class="bold upper">%s</div>
    <div class="bold upper">TAHUN ANGGARAN %d</div>
    <div class="bold upper" style="margin-top:8px;">%s %s</div>
  </div>

  <!-- MENIMBANG/MENGINGAT -->
  <table class="full" style="margin-bottom:8px;">
    <tr>
      <td style="width:20%%; vertical-align:top;" class="bold">Menimbang</td>
      <td style="width:2%%; vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        <ol type="a">
          <li>Bahwa untuk tertib dan lancarnya pelaksanaan pengadaan barang/jasa di %s %s %s Tahun %d dipandang perlu menetapkan Pokja Pengadaan Barang/Jasa;</li>
          <li>bahwa berdasarkan pertimbangan sebagaimana dimaksud dalam huruf a, perlu menetapkan Keputusan %s tentang Penetapan Pokja Pengadaan Barang/Jasa pada %s %s %s Tahun Anggaran %d;</li>
        </ol>
      </td>
    </tr>
    <tr><td colspan="3" style="height:8px;"></td></tr>
    <tr>
      <td style="vertical-align:top;" class="bold">Mengingat</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        <ol>
          <li>Undang-Undang Nomor 27 Tahun 1959 tentang Penetapan Undang-Undang Darurat Nomor 3 Tahun 1953 tentang Perpanjangan Pembentukan Daerah Tingkat II di Kalimantan (Lembaran Negara Republik Indonesia tahun 1953 Nomor 9) sebagai Undang-Undang (Lembaran Negara Republik Indonesia Tahun 1959 Nomor 72, Tambahan Lembaran Negara Republik Indonesia Nomor 1820);</li>
          <li>Undang-Undang Nomor 12 Tahun 2011 tentang Pembentukan Peraturan Perundang-Undangan (Lembaran Negara Republik Indonesia Tahun 2011 Nomor 82, Tambahan Lembaran Negara Republik Indonesia Nomor 5234) sebagaimana telah diubah beberapa kali terakhir dengan Undang-Undang Nomor 13 Tahun 2022 tentang Perubahan Kedua atas Undang-Undang Nomor 12 Tahun 2011 tentang Pembentukan Peraturan Perundang-Undangan (Lembaran Negara Republik Indonesia Tahun 2022 Nomor 143, Tambahan Lembaran Negara Republik Indonesia Nomor 6801);</li>
          <li>Undang-Undang Nomor 23 Tahun 2014 tentang Pemerintahan Daerah (Lembaran Negara Republik Indonesia Tahun 2014 Nomor 244, Tambahan Lembaran Negara Republik Indonesia Nomor 5587) sebagaimana telah diubah beberapa kali terakhir dengan Undang-Undang Nomor 6 Tahun 2023 tentang Penetapan Peraturan Pemerintah Pengganti Undang-Undang Nomor 2 Tahun 2022 tentang Cipta Kerja menjadi Undang-Undang (Lembaran Negara Republik Indonesia 2023 Nomor 41, Tambahan Lembaran Negara Republik Indonesia Nomor 6856);</li>
          <li>Undang-Undang Nomor 1 Tahun 2022 tentang Hubungan Keuangan Antara Pemerintah Pusat dan Pemerintahan Daerah (Lembaran Negara Republik Indonesia Tahun 2022 Nomor 4, Tambahan Lembaran Negara Republik Indonesia Nomor 6757);</li>
          <li>Undang-Undang Nomor 8 Tahun 2022 tentang Provinsi Kalimantan Selatan (Lembaran Negara Republik Indonesia Tahun 2022 Nomor 68, Tambahan Lembaran Negara Republik Indonesia Nomor 6779);</li>
          <li>Undang-Undang Nomor 20 Tahun 2023 tentang Aparatur Sipil Negara (Lembaran Negara Republik Indonesia Tahun Nomor 141, Tambahan Lembaran Negara Republik Indonesia Nomor 6897);</li>
          <li>Peraturan Pemerintah Nomor 12 Tahun 2019 tentang Pengelolaan Keuangan Daerah (Lembaran Negara Republik Indonesia Tahun 2019 Nomor 42, Tambahan Lembaran Negara Republik Indonesia Nomor 6322);</li>
          <li>Peraturan Presiden Nomor 16 Tahun 2018 tentang Pengadaan Barang/Jasa Pemerintah (Lembaran Negara Republik Indonesia Tahun 2018 Nomor 33) sebagaimana telah diubah dengan Peraturan Presiden Nomor 12 Tahun 2021 tentang Perubahan atas Peraturan Presiden Nomor 16 Tahun 2018 tentang Pengadaan Barang/Jasa Pemerintah (Lembaran Negara Republik Indonesia Tahun 2021 Nomor 63);</li>
          <li>Peraturan Menteri Dalam Negeri Nomor 80 Tahun 2015 tentang Pembentukan Produk Hukum Daerah (Berita Negara Republik Indonesia Tahun 2015 Nomor 2036) sebagaimana telah diubah dengan Peraturan Menteri Dalam Negeri Nomor 120 Tahun 2018 tentang Perubahan atas Peraturan Menteri Dalam Negeri Nomor 80 Tahun 2015 tentang Pembentukan Produk Hukum Daerah (Berita Negara Republik Indonesia Tahun 2019 Nomor 157);</li>
          <li>Peraturan Menteri Dalam Negeri Nomor 77 Tahun 2020 tentang Pedoman Teknis Pengelolaan Keuangan Daerah (Berita Negara Republik Indonesia Tahun 2020 Nomor 1781);</li>
          <li>Peraturan Lembaga Kebijakan Pengadaan Barang/Jasa Pemerintah Nomor 11 Tahun 2021 tentang Pedoman Perencanaan Pengadaan Barang/Jasa Pemerintah (Berita Negara Republik Indonesia Tahun 2021 Nomor 512);</li>
          <li>Peraturan Daerah %s Nomor 7 Tahun 2016 tentang Pembentukan dan Susunan Perangkat Daerah %s (Lembaran Daerah %s Tahun 2016 Nomor 7, Tambahan Lembaran Daerah %s Nomor 40) sebagaimana telah diubah dengan Peraturan Daerah %s Nomor 3 Tahun 2021 tentang Perubahan atas Peraturan Daerah %s Nomor 7 Tahun 2016 (Lembaran Daerah %s Tahun 2021 Nomor 3, Tambahan Lembaran Daerah %s Nomor 63);</li>
          <li>Peraturan Daerah %s Nomor 7 Tahun tentang Keuangan Daerah (Lembaran Daerah %s Tahun 2021 Nomor 7, Tambahan Lembaran Daerah %s Nomor 66);</li>
        </ol>
      </td>
    </tr>
    <tr><td colspan="3" style="height:8px;"></td></tr>
    <tr>
      <td style="vertical-align:top;" class="bold">Memperhatikan</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Dokumen Pelaksanaan Anggaran Satuan Kerja Perangkat Daerah %s %s %s Tahun Anggaran %d</td>
    </tr>
  </table>

  <div class="center bold upper" style="margin:15px 0;">Memutuskan</div>

  <table class="full">
    <tr>
      <td style="width:20%%; vertical-align:top;" class="bold">Menetapkan</td>
      <td style="width:2%%; vertical-align:top;">:</td>
      <td style="vertical-align:top;"></td>
    </tr>
    <tr>
      <td style="vertical-align:top;" class="bold">KESATU</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        Pokja Pengadaan Barang/Jasa Pada %s %s %s Tahun Anggaran %d sebagai berikut:<br><br>
        %s
      </td>
    </tr>
    <tr><td colspan="3" style="height:8px;"></td></tr>
    <tr>
      <td style="vertical-align:top;" class="bold">KEDUA</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">
        Tugas Pokja Pengadaan Barang/Jasa sebagaimana dimaksud dalam Diktum KESATU adalah sebagai berikut:
        <ol type="a">
          <li>melaksanakan persiapan dan pelaksanaan pemilihan Penyedia;</li>
          <li>melaksanakan persiapan dan pelaksanaan pemilihan Penyedia untuk katalog elektronik; dan</li>
          <li>menetapkan pemenang pemilihan/Penyedia untuk metode pemilihan:
            <ol>
              <li>Tender/Penunjukan Langsung untuk paket Pengadaan Barang/Pekerjaan Konstruksi/Jasa Lainnya dengan nilai Pagu Anggaran paling banyak Rp. 100.000.000.000,00 (seratus miliar rupiah); dan</li>
              <li>Seleksi/Penunjukan Langsung untuk paket Pengadaan Jasa Konsultansi dengan nilai Pagu Anggaran paling banyak Rp. 10.000.000.000,00 (sepuluh miliar rupiah).</li>
            </ol>
          </li>
          <li>menyerahkan dokumen hasil pemilihan Penyedia Barang/Jasa kepada Pengguna Anggaran/Kuasa Pengguna Anggaran;</li>
          <li>membuat laporan mengenai proses dan hasil pengadaan kepada Kepala Daerah/Pimpinan Institusi; dan</li>
          <li>memberikan pertanggungjawaban atas pelaksanaan kegiatan Pengadaan Barang/Jasa kepada Pengguna Anggaran/Kuasa Pengguna Anggaran</li>
        </ol>
      </td>
    </tr>
    <tr><td colspan="3" style="height:8px;"></td></tr>
    <tr>
      <td style="vertical-align:top;" class="bold">KETIGA</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Segala biaya yang timbul akibat ditetapkannya Keputusan %s ini dibebankan pada Anggaran Pendapatan dan Belanja Daerah %s Tahun Anggaran %d pada %s %s %s.</td>
    </tr>
    <tr><td colspan="3" style="height:8px;"></td></tr>
    <tr>
      <td style="vertical-align:top;" class="bold">KEEMPAT</td>
      <td style="vertical-align:top;">:</td>
      <td style="vertical-align:top;">Keputusan %s ini mulai berlaku pada tanggal ditetapkan.</td>
    </tr>
  </table>

  <!-- TTD -->
  <div style="float:right; width:35%%; margin-top:30px; text-align:left;">
    <table style="width:100%%; border:none;">
      <tr>
        <td style="padding:0; line-height:1.2;">
          <p style="margin:0;">Ditetapkan di %s</p>
          <p style="margin:0;">pada tanggal %s</p>
          <div style="margin-top:10px;" class="bold upper">%s</div>
          <div class="bold upper">%s</div>
          %s
          <div class="bold upper underline" style="margin-top:5px;">%s</div>
          <p style="margin:0;">NIP. %s</p>
        </td>
      </tr>
    </table>
  </div>
  <div style="clear:both;"></div>
  <div style="clear:both;"></div>
</body>
</html>`,
		// Kop
		logoHtml, docInstansi, docSub, docAddr, kopContact, kopContact2, qrCell,
		// Judul
		docPejabatJabatan, docRegion,
		noSurat, tahunSurat,
		panitia.Nama, docSub, tahunSurat,
		docPejabatJabatan, docRegion,
		// Menimbang a
		panitia.Nama, docSub, docRegion, tahunSurat,
		// Menimbang b
		docPejabatJabatan, panitia.Nama, docSub, docRegion, tahunSurat,
		// Mengingat 12 - Perda (8 placeholders)
		docRegion, docRegion, docRegion, docRegion, docRegion, docRegion, docRegion, docRegion,
		// Mengingat 13 (3 placeholders)
		docRegion, docRegion, docRegion,
		// Memperhatikan
		panitia.Nama, docSub, docRegion, tahunSurat,
		// KESATU
		panitia.Nama, docSub, docRegion, tahunSurat,
		anggotaRows,
		// KETIGA
		docPejabatJabatan, docRegion, tahunSurat, panitia.Nama, docSub, docRegion,
		// KEEMPAT
		docPejabatJabatan,
		// TTD
		tempatPenetapan, tglPenetapan,
		docPejabatJabatan, docRegion,
		docSignatureHtml,
		docPejabatNama, settings.DocPejabatNip,
	)
	return html
}

func CetakSkPokja(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	
	// Check if already printed, if not create record with defaults
	var dok models.DokumenTercetak
	models.GetDB().Where("paket_id = ? AND jenis_dokumen = ?", id, "SK_PENUNJUKAN_POKJA").First(&dok)
	
	paket := services.GetPaket(id)
	pokjaMaster := paket.Pokja()

	if dok.ID == 0 {
		dok = models.DokumenTercetak{
			PaketID:      id,
			JenisDokumen: "SK_PENUNJUKAN_POKJA",
			NomorSurat:   pokjaMaster.NoSk,
			TahunSurat:   fmt.Sprintf("%d", pokjaMaster.Tahun),
			TempatPenetapan: pokjaMaster.TempatSk,
			TanggalPenetapan: pokjaMaster.TglSk,
		}
		if dok.TempatPenetapan == "" { dok.TempatPenetapan = "Banjarmasin" }
		if dok.TanggalPenetapan.IsZero() { dok.TanggalPenetapan = time.Now() }
		
		session := getUserSession(c)
		if session.Pegawai().ID > 0 { dok.PembuatPegawaiID = session.Pegawai().ID }
		models.GetDB().Create(&dok)
	}

	appSettings := services.GetSettings()
	fullUrl := fmt.Sprintf("%s://%s/preview/sk-pokja/%d?hash=%s", c.Protocol(), c.Hostname(), id, dok.Md5Hash)

	html := buildSkPokjaHtml(pokjaMaster, appSettings, dok, fullUrl)
	log.Infof("Generating SK Pokja PDF. HTML length: %d", len(html))

	result := utils.ExportHtmlToPdf(html, "")
	if result == nil {
		log.Error("Failed to generate SK Pokja PDF result (is nil)")
		return c.Status(500).SendString("Gagal membuat PDF SK Pokja. Pastikan wkhtmltopdf terinstal di server.")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"SK-pokja.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}

func CetakSkPokjaProcess(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	dok := processFormDokumen(c, id, "SK_PENUNJUKAN_POKJA")
	mp := currentMap(c)
	paket := services.GetPaket(id)
	mp["paket"] = paket
	mp["pokja"] = paket.Pokja()
	mp["appSettings"] = services.GetSettings()
	mp["dokTercetak"] = dok
	valUrl := fmt.Sprintf("%s://%s/validasi/dokumen/%s", c.Protocol(), c.Hostname(), dok.Md5Hash)
	mp["qrValidasi"] = generateQrBase64(valUrl)
    if !dok.TanggalPenetapan.IsZero() {
        mp["tglPenetapan"] = dok.TanggalPenetapan.Format("02-01-2006")
    }

	html, _ := renderToString("preview/surat-penunjukan-pokja", mp)
	result := utils.ExportHtmlToPdf(html, "")
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"SK-pokja.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}
func prepareBAKajiUlangData(c *fiber.Ctx) fiber.Map {
	id := utils.StringToUint(c.Params("id"))
	mp := currentMap(c)
	paket := services.GetPaket(id)
	mp["paket"] = paket
	
	// Data Signers (TTE)
	type Signer struct {
		Nama     string
		Jabatan  string
		IsSigned bool
		QrCode   string
		SigImg   string 
		Nip      string
	}

	var signersPPK []Signer
	var signersProses []Signer

	// 1. Ambil Data PPK
	ppk := paket.Ppk()
	if ppk.ID > 0 {
		isSigned := false
		dokPersiapans := paket.DokPersiapan()
		if len(dokPersiapans) > 0 {
			p := dokPersiapans[0].PersetujuanPegawai(ppk.ID)
			isSigned = p.Status
		}

		s := Signer{
			Nama:     ppk.PegNama,
			Jabatan:  "Pejabat Pembuat Komitmen",
			IsSigned: isSigned,
			Nip:      ppk.PegNip,
		}
		if isSigned {
			fullUrl := fmt.Sprintf("%s://%s/preview/ba-kajiulang/%d/verify/%d", c.Protocol(), c.Hostname(), id, ppk.ID)
			s.QrCode = generateQrBase64(fullUrl)
		}
		signersPPK = append(signersPPK, s)
	}

	// 2. Ambil Data Pelaksana (PP atau Pokja)
	if paket.PpId > 0 {
		pp := paket.Pp()
		isSigned := false
		dokPersiapans := paket.DokPersiapan()
		if len(dokPersiapans) > 0 {
			p := dokPersiapans[0].PersetujuanPegawai(pp.ID)
			isSigned = p.Status
		}

		s := Signer{
			Nama:     pp.PegNama,
			Jabatan:  "Pejabat Pengadaan",
			IsSigned: isSigned,
			Nip:      pp.PegNip,
		}
		if isSigned {
			fullUrl := fmt.Sprintf("%s://%s/preview/ba-kajiulang/%d/verify/%d", c.Protocol(), c.Hostname(), id, pp.ID)
			s.QrCode = generateQrBase64(fullUrl)
		}
		signersProses = append(signersProses, s)
	} else if paket.PntId > 0 {
		panitia := paket.Pokja()
		anggota := panitia.AnggotaList()
		dokPersiapans := paket.DokPersiapan()
		
		for _, a := range anggota {
			isSigned := false
			if len(dokPersiapans) > 0 {
				p := dokPersiapans[0].PersetujuanPegawai(a.ID)
				isSigned = p.Status
			}
			s := Signer{
				Nama:     a.PegNama,
				Jabatan:  "Anggota Pokja",
				IsSigned: isSigned,
				Nip:      a.PegNip,
			}
			if isSigned {
				fullUrl := fmt.Sprintf("%s://%s/preview/ba-kajiulang/%d/verify/%d", c.Protocol(), c.Hostname(), id, a.ID)
				s.QrCode = generateQrBase64(fullUrl)
			}
			signersProses = append(signersProses, s)
		}
	}

	// 3. Metadata
	var ba models.BeritaAcara
	resultBA := models.GetDB().Where("pkt_id = ? AND jenis = 'REVIU'", id).First(&ba)
	if resultBA.Error != nil {
		// If not found, we still provide an empty object with basic info to avoid template errors
		ba = models.BeritaAcara{
			PktId: id,
			Jenis: "REVIU",
		}
	}
	
	reviuMaster := services.GetAllReviu()
	var reviuResults []models.ReviuPaket
	models.GetDB().Where("pkt_id = ?", id).Find(&reviuResults)
	
	resMap := make(map[uint]models.ReviuPaket)
	for _, r := range reviuResults {
		resMap[r.RevId] = r
	}

	// 4. Signatures
	for i := range signersPPK {
		doc := models.GetDocumentByJenis(ppk.ID, models.TTD)
		if doc.ID > 0 { signersPPK[i].SigImg = services.GetBase64FromFile(doc.Filepath) }
	}
	for i := range signersProses {
		var pid uint
		if paket.PpId > 0 { 
			pid = paket.PpId 
		} else if paket.PntId > 0 && len(paket.Pokja().AnggotaList()) > i {
			pid = paket.Pokja().AnggotaList()[i].ID 
		}
		
		if pid > 0 {
			doc := models.GetDocumentByJenis(pid, models.TTD)
			if doc.ID > 0 { signersProses[i].SigImg = services.GetBase64FromFile(doc.Filepath) }
		}
	}

	// 5. Foto Rapat
	fotoRapat := models.GetDokPaketJenis(id, models.FOTO_RAPAT)
	// Embed foto sebagai base64 agar tampil di PDF (tidak tergantung URL)
	fotoBase64 := ""
	if fotoRapat.ID > 0 {
		dok := fotoRapat.Document()
		if dok.ID > 0 {
			fotoBase64 = services.GetBase64FromFile(dok.Filepath)
		}
	}
	mp["fotoBase64"] = fotoBase64
	mp["fotoRapat"] = fotoRapat

	appSettings := services.GetSettings()
	mp["appSettings"] = appSettings
	mp["ba"] = ba
	mp["reviuMaster"] = reviuMaster
	mp["reviuResults"] = resMap
	mp["signersPPK"] = signersPPK
	mp["signersProses"] = signersProses
	return mp
}

func generateQrBase64(content string) string {
	qrPng, qrErr := qrcode.Encode(content, qrcode.Medium, 128)
	if qrErr != nil {
		log.Error("Generate QR Code failed: ", qrErr)
		return ""
	}
	b64 := base64.StdEncoding.EncodeToString(qrPng)
	log.Infof("QR Code generated for %s (length: %d)", content, len(b64))
	return b64
}

func PreviewBAKajiUlang(c *fiber.Ctx) error {
	mp := prepareBAKajiUlangData(c)
	hash := c.Query("hash")
	if hash != "" {
		var dok models.DokumenTercetak
		models.GetDB().Where("md5_hash = ?", hash).First(&dok)
		mp["dokTercetak"] = dok
		valUrl := fmt.Sprintf("%s://%s/validasi/dokumen/%s", c.Protocol(), c.Hostname(), hash)
		mp["qrValidasi"] = generateQrBase64(valUrl)
	} else {
		// Always provide a QR code even for draft preview
		valUrl := fmt.Sprintf("%s://%s/preview/ba-kajiulang/%s", c.Protocol(), c.Hostname(), c.Params("id"))
		qr := generateQrBase64(valUrl)
		mp["qrValidasi"] = qr
		mp["qr_validasi"] = qr
	}
	return c.Render("preview/ba-kajiulang", mp)
}

func CetakBAKajiUlang(c *fiber.Ctx) error {
	// Tidak mengunci di sini - paket dikunci otomatis setelah semua pihak menyetujui
	// Gunakan ExportToPdf dengan URL langsung agar gambar/logo ter-render dengan benar
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3002"
	}
	previewUrl := fmt.Sprintf("http://127.0.0.1:%s/preview/ba-kajiulang/%s", port, c.Params("id"))
	result := utils.ExportToPdf(previewUrl)
	if result == nil || len(result) == 0 {
		// Fallback ke renderToString jika ExportToPdf gagal
		mp := prepareBAKajiUlangData(c)
		valUrl := fmt.Sprintf("%s://%s/preview/ba-kajiulang/%s", c.Protocol(), c.Hostname(), c.Params("id"))
		qr := generateQrBase64(valUrl)
		mp["qrValidasi"] = qr
		mp["qr_validasi"] = qr
		html, err := renderToString("preview/ba-kajiulang", mp)
		if err != nil {
			log.Errorf("renderToString error: %v", err)
		}
		result = utils.ExportHtmlToPdf(html, "")
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"BA-kajiulang.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}

func CetakBAKajiUlangProcess(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	dok := processFormDokumen(c, id, "BA_KAJIULANG")
	// Tidak mengunci di sini - paket dikunci otomatis setelah semua pihak menyetujui
	mp := prepareBAKajiUlangData(c)
	mp["dokTercetak"] = dok
	valUrl := fmt.Sprintf("%s://%s/validasi/dokumen/%s", c.Protocol(), c.Hostname(), dok.Md5Hash)
	qr := generateQrBase64(valUrl)
	mp["qrValidasi"] = qr
	mp["qr_validasi"] = qr
	
	html, _ := renderToString("preview/ba-kajiulang", mp)
	result := utils.ExportHtmlToPdf(html, "")
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"BA-kajiulang.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}


func PreviewBANego(c *fiber.Ctx)  error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	mp["appSettings"] = services.GetSettings()
	return c.Render("preview/ba-nego", mp)
}


func CetakBANego(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	mp["appSettings"] = services.GetSettings()

	html, _ := renderToString("preview/ba-nego", mp)
	result := utils.ExportHtmlToPdf(html, "")
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"BA-Nego.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}


func PreviewBAPenetapan(c *fiber.Ctx)  error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	mp["appSettings"] = services.GetSettings()
	return c.Render("preview/ba-penetapan", mp)
}

func CetakBAPenetapan(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.StringToUint(c.Params("id"))
	berita_acara := services.GetBeritaAcara(id)
	mp["berita_acara"] = berita_acara
	mp["appSettings"] = services.GetSettings()

	html, _ := renderToString("preview/ba-penetapan", mp)
	result := utils.ExportHtmlToPdf(html, "")
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"BA-Penetapan-Pemenang.pdf\"")
	return c.SendStream(bytes.NewReader(result))
}

func VerifyTte(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	pegId := utils.StringToUint(c.Params("pegId"))
	mp := currentMap(c)

	paket := services.GetPaket(id)
	if paket.ID == 0 {
		return c.Status(404).SendString("Paket tidak ditemukan")
	}

	pegawai := services.GetPegawai(pegId)
	if pegawai.ID == 0 {
		return c.Status(404).SendString("Pegawai tidak ditemukan")
	}

	// Cek status persetujuan
	isSigned := false
	dokPersiapans := paket.DokPersiapan()
	if len(dokPersiapans) > 0 {
		p := dokPersiapans[0].PersetujuanPegawai(pegId)
		isSigned = p.Status
	}

	mp["paket"] = paket
	mp["pegawai"] = pegawai
	mp["isSigned"] = isSigned
	mp["tglVerifikasi"] = time.Now().Format("02-01-2006 15:04:05")

	return c.Render("preview/verify-tte", mp)
}

func print(c *fiber.Ctx, url string, filename string) error {
	result := utils.ExportToPdf(url)
	reader := bytes.NewReader(result)
	// Set the Content-Type header for PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.SendStream(reader)
}

func ValidasiDokumenTercetak(c *fiber.Ctx) error {
	hash := c.Params("hash")
	var dok models.DokumenTercetak
	models.GetDB().Where("md5_hash = ?", hash).First(&dok)

	if dok.ID == 0 {
		return c.Status(404).SendString("Peringatan: Dokumen tidak valid atau tidak terdeteksi di server kami! Hati-hati pemalsuan.")
	}

	paket := services.GetPaket(dok.PaketID)
	// Build simple JSON response or simple HTML
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head><title>Validasi Dokumen</title><style>body{font-family:Arial,sans-serif;background:#eee;margin:0;padding:50px;} .card{background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 15px rgba(0,0,0,0.1);max-width:600px;margin:auto;} h2{color:#28a745;} p{font-size:16px;} .badge{background:#007bff;color:white;padding:5px 10px;border-radius:20px;font-size:12px;}</style></head>
	<body>
		<div class="card">
			<center>
				<h2 style="margin-top:0;">DOKUMEN TERVALIDASI <i data-feather="check-circle"></i></h2>
				<p>Dokumen ini telah disahkan dan di-generate oleh Sistem Informasi Pengadaan Daerah pada tanggal <strong>%s</strong>.</p>
			</center>
			<hr>
			<table cellpadding="5">
				<tr><td><strong>Jenis Surat:</strong></td><td><span class="badge">%s</span></td></tr>
				<tr><td><strong>Nomor Surat:</strong></td><td>%s</td></tr>
				<tr><td><strong>Paket Terkait:</strong></td><td>%s (ID: %d)</td></tr>
				<tr><td><strong>Tahun/Tanggal:</strong></td><td>%s / %s</td></tr>
				<tr><td><strong>Barcode Hash:</strong></td><td><code>%s</code></td></tr>
			</table>
			<br>
			<div style="background:#e9f7ef;color:#155724;padding:10px;border-radius:5px;border:1px solid #c3e6cb;text-align:center;font-size:14px;">Checksum Valid. Integritas Data Dinyatakan Aman.</div>
		</div>
	</body>
	</html>
	`, dok.CreatedAt.Format("02 Jan 2006 15:04:05"), dok.JenisDokumen, dok.NomorSurat, paket.Nama, paket.ID, dok.TahunSurat, dok.TanggalPenetapan.Format("02-01-2006"), hash)

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
