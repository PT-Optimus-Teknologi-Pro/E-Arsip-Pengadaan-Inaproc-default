package services

import (
	"arsip/config"
	"arsip/models"
	"encoding/base64"
	"errors"
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
	dir := path.Dir(filename)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Error("create dir failed ", err)
		return err
	}
	i := strings.Index(data, ",")
	if i < 0 {
		log.Fatal("no comma")
		return errors.New("no comma")
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
	png.Encode(f, im)
	return models.SaveDocumentByJenis(id, models.TTD, filename)
}
