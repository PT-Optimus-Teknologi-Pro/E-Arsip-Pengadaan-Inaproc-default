package services

import (
	"arsip/config"
	"arsip/models"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetDocument(id uint) models.Document {
	return models.GetDocument(id)
}

func GetDocumentPegawai(pegId uint) []models.Document {
	return models.GetDocumentPegawai(pegId)
}

func GetAdminDocuments() []models.Document {
	return models.GetAllDocumentByJenis(models.ADMIN_DOK)
}

func DeleteDocument(document models.Document) error {
	err := os.Remove(document.Filepath)
	if err != nil {
		return err
	}
	return models.DeleteDocument(&document)
}

func GetAllDokTemplate() []models.DokTemplate {
	return models.GetAllDokTemplate()
}

func GetDokTemplate(id uint) models.DokTemplate {
	return models.GetDokTemplate(id)
}

func SaveDocTemplate(c *fiber.Ctx, template models.DokTemplate, userid uint) error {
	template.DokId, _ = models.SaveDocument(c, userid, models.DOK_TEMPLATE, "file")
	return models.SaveDocTemplate(template)
}

func DeleteDocTemplate(template models.DokTemplate) error {
	return models.DeleteDocTemplate(template)
}

func SaveTTD(c *fiber.Ctx, id uint) error {
	data := c.FormValue("ttd")
	filename := config.UploadPath() + "/" + strconv.Itoa(int(id)) + "/signed.png"
	return saveBase64Image(data, filename, id, models.TTD)
}

func SaveSignatureBA(c *fiber.Ctx, pktId uint, pegId uint, data string) error {
	// Unique filename for this packet's signature by this specific user
	filename := config.UploadPath() + "/ba_reviu/" + strconv.Itoa(int(pktId)) + "_" + strconv.Itoa(int(pegId)) + "_sig.png"
	return saveBase64Image(data, filename, pegId, models.TTD) // Reusing TTD type for now
}

func saveBase64Image(data string, filename string, pegId uint, jenis string) error {
	dir := path.Dir(filename)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Error("create dir failed ", err)
		return err
	}
	i := strings.Index(data, ",")
	if i < 0 {
		return errors.New("Invalid base64 data")
	}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[i+1:]))
	im, err := png.Decode(dec)
	if err != nil {
		log.Error(err)
		return errors.New("Bad png")
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
		return errors.New("Cannot open file")
	}
	defer f.Close()
	png.Encode(f, im)
	return models.SaveDocumentByJenis(pegId, jenis, filename)
}

func GetBase64FromFile(filepath string) string {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Error("Read file failed ", err)
		return ""
	}
	ext := strings.ToLower(path.Ext(filepath))
	mime := "image/png"
	if ext == ".jpg" || ext == ".jpeg" {
		mime = "image/jpeg"
	}
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(data))
}
