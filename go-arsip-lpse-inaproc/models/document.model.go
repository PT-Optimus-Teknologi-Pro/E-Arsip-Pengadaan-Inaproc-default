package models

import (
	"arsip/config"
	"arsip/utils"
	"errors"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const (
	KTP          	string = "KTP"
	SK           	string = "SK"
	TTD          	string = "TTD"
	SERTIFIKAT   	string = "SERTIFIKAT"
	DOK_TEMPLATE 	string = "DOK_TEMPLATE"
	CHECKLIST	 	string = "CHECKLIST"
	KAJIULANG	 	string = "KAJIULANG"
	DOKFINAL	 	string = "DOKFINAL"
	BUKUTAMU	 	string = "BUKUTAMU"
	PERUBAHAN_DATA 	string = "PERUBAHAN_DATA"
	PENDUKUNG	 	string = "PENDUKUNG"
	HASIL_PENGADAAN string = "HASIL_PENGADAAN"
	HASIL_PEKERJAAN string = "HASIL_PEKERJAAN"
	KONTRAK			string = "KONTRAK"
)

type Document struct {
	gorm.Model
	Versi int    `gorm:"primaryKey" json:"versi"`
	Jenis string `json:"jenis"`
	PegId uint   `json:"peg_id"`
	//UpdatedBy uint
	Filename string    `json:"filename"`
	Filesize int64     `json:"filesize"`
	Filepath string    `json:"filepath"`
	Filehash string    `json:"filehash"`
	Filedate time.Time `json:"filedate"`
}

func (Document) TableName() string {
	return "document"
}

func (c Document) Label() string {
	if c.Jenis == KTP {
		return "Scan KTP"
	} else if c.Jenis == SK {
		return "Scan SK"
	} else if c.Jenis == TTD {
		return "Tanda Tangan"
	} else if c.Jenis == SERTIFIKAT {
		return "Sertifikat Pengadaan"
	} else if c.Jenis == CHECKLIST {
		return "Checklist"
	}
	return ""
}

func GetDocument(id uint) Document {
	var rest Document
	db.First(&rest, id)
	return rest
}

func GetDocumentByJenis(id uint, jenis string) Document {
	var document Document
	db.Where("id = ? and jenis = ? and deleted_at IS NULL", id, jenis).First(&document)
	return document
}

func GetDocumentPegawai(pegId uint) []Document {
	var rest []Document
	db.Find(&rest, "peg_id = ?", pegId)
	return rest
}

func DeleteDocument(document *Document) error {
	return db.Delete(document).Error
}

// func SaveDocument(document *Document) error {
// 	return db.Save(document).Error
// }

func GetNextSequenceDokumen() uint {
 var res uint
  db.Raw("SELECT nextval('document_id_seq')").Scan(&res)
  return res
}

func SaveDocument(c *fiber.Ctx, id uint, jenis string, name string) (uint, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return 0, err
	}
	if file.Size == 0 {
		return 0, errors.New("file is empty")
	}
	documentID := GetNextSequenceDokumen()
	document := Document{
		Model: gorm.Model{
			ID: documentID,
		},
		Jenis: jenis,
		PegId: id,
		Versi: 1,
		Filename: file.Filename,
		Filesize: file.Size,
		Filedate: time.Now(),
	}
	// err = models.SaveDocument(&document)
	// if err != nil {
	// 	return 0, err
	// }
	destination := config.UploadPath() + "/" + utils.UintToString(document.ID) +"/"+utils.IntToString(document.Versi)+ "/" + file.Filename
	err = saveFile(c, destination, file)
	if err != nil {
		return 0, err
	}
	document.Filepath = destination
	document.Filehash = utils.HashFile(destination)
	err = db.Save(&document).Error
	if err != nil {
		return 0, err
	}
	return document.ID, nil
}

func saveFile(c *fiber.Ctx, filename string, file *multipart.FileHeader) error {
	dir := path.Dir(filename)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Error("create dir failed ", err)
		return err
	}
	return c.SaveFile(file, filename)
}


func SaveDocumentByJenis(id uint, jenis string, filepath string) error {
	document := GetDocumentByJenis(id, jenis)
	if document.ID == 0 {
		document.Jenis = jenis
		document.PegId = id
	}
	fileInfo, _ := os.Stat(filepath)
	document.Versi += 1
	document.Filename = fileInfo.Name()
	document.Filedate = time.Now()
	document.Filesize = fileInfo.Size()
	document.Filepath = filepath
	document.Filehash = utils.HashFile(filepath)
	return db.Save(&document).Error
}
