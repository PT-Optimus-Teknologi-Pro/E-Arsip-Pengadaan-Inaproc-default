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
	htmlContent := buildSkPpHtml(sk, pegawai, satker, appSettings, fullUrl)

	result := utils.ExportHtmlToPdf(htmlContent, "")
	if len(result) == 0 {
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

func buildSkPpHtml(sk models.PejabatPengadaan, pegawai models.Pegawai, satker models.SatkerSirup, settings models.AppSettings, fullUrl string) string {
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
    <br><br><br><br>
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
		func() string { if sk.TempatSk != "" { return sk.TempatSk }; return "Banjarmasin" }(),
		func() string { if !sk.TglSk.IsZero() { return sk.TglSk.Format("02-01-2006") }; return "..." }(),
		docPejabatJabatan, docInstansi, docPejabatNama, settings.DocPejabatNip,
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
	log.Info("priview BA Reviu Dokumen")
	mp := currentMap(c)
	paket := services.GetPaket(uint(id))
	mp["paket"] = paket
	
	// Data Signers (TTE)
	type Signer struct {
		Nama     string
		Jabatan  string
		IsSigned bool
		QrCode   string
		SigImg   string // NEW: Base64 signature image
		Nip      string
	}

	var signersPPK []Signer
	var signersProses []Signer

	// 1. Ambil Data PPK
	ppk := paket.Ppk()
	if ppk.ID > 0 {
		isSigned := false
		// Cek apakah PPK sudah menyetujui dokumen persiapan
		dokPersiapans := paket.DokPersiapan()
		if len(dokPersiapans) > 0 {
			// Jika ada setidaknya satu dokumen dan sudah disetujui PPK
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

	// 3. New: Fetch BA Metadata and Checklist Results
	var ba models.BeritaAcara
	models.GetDB().Where("pkt_id = ? AND jenis = 'REVIU'", id).First(&ba)
	
	reviuMaster := services.GetAllReviu()
	var reviuResults []models.ReviuPaket
	models.GetDB().Where("pkt_id = ?", id).Find(&reviuResults)
	
	resMap := make(map[uint]models.ReviuPaket)
	for _, r := range reviuResults {
		resMap[r.RevId] = r
	}

	// 4. Fetch Actual Signature Images (PNG) for each signer
	for i := range signersPPK {
		doc := models.GetDocumentByJenis(ppk.ID, models.TTD) // Simplified: assuming one TTD per user
		if doc.ID > 0 {
			signersPPK[i].SigImg = services.GetBase64FromFile(doc.Filepath)
		}
	}
	
	// For process signers (PP/Pokja)
	for i := range signersProses {
		// We need to find the PegId for this signer to get their SigImg
		var pid uint
		if paket.PpId > 0 { pid = paket.PpId } else { pid = paket.Pokja().AnggotaList()[i].ID }
		
		doc := models.GetDocumentByJenis(pid, models.TTD)
		if doc.ID > 0 {
			signersProses[i].SigImg = services.GetBase64FromFile(doc.Filepath)
		}
	}

	mp["ba"] = ba
	mp["reviuMaster"] = reviuMaster
	mp["reviuResults"] = resMap
	mp["signersPPK"] = signersPPK
	mp["signersProses"] = signersProses

	return c.Render("preview/ba-kajiulang", mp)
}

func generateQrBase64(content string) string {
	qrPng, qrErr := qrcode.Encode(content, qrcode.Medium, 128)
	if qrErr != nil {
		log.Error("Generate QR Code failed: ", qrErr)
		return ""
	}
	return base64.StdEncoding.EncodeToString(qrPng)
}


func CetakBAKajiUlang(c *fiber.Ctx) error {
	log.Info("cetak BA Reviu Dokumen")
	id := utils.StringToUint(c.Params("id"))
	paket := services.GetPaket(id)
	if paket.ID > 0 {
		paket.IsLockedReview = true
		models.SavePaket(&paket)
	}
	url := fmt.Sprintf("http://localhost:%s/preview/ba-kajiulang/%s", config.Port(), c.Params("id"))
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
