package services

import (
	"arsip/models"
	"github.com/gofiber/fiber/v2"
)


func GetDataTableAgency(c *fiber.Ctx) error {
	return models.GetDataTableAgency(c)
}

func GetDataTableUkpbj(c *fiber.Ctx) error {
	return models.GetDataTableUkpbj(c)
}

func GetDataTablePerubahanData(c *fiber.Ctx, usrsession models.UserSession) error {
	return models.GetDataTablePerubahanData(c, usrsession)
}

func GetDataTableVerifikasi(c *fiber.Ctx) error {
	return models.GetDataTableVerifikasi(c)
}

func GetDataTablePaketSirup(c *fiber.Ctx,  tahun int, satker string, metode string, jenis string) error {
	return models.GetDataTablePaketSirup(c, tahun, satker, metode, jenis)
}

func GetDataTableSwakelolaSirup(c *fiber.Ctx) error {
	return models.GetDataTableSwakelolaSirup(c)
}

func GetDataTablePegawai(c *fiber.Ctx, usrgroup string) error {
	return models.GetDataTablePegawai(c, usrgroup)
}

func GetDataTablePaket(c *fiber.Ctx, id uint, isPPK, isUkpbj, isPokja, isPp bool) error {
	return models.GetDataTablePaket(c, id, isPPK, isUkpbj, isPokja, isPp)
}

func GetDataTableTemplates(c *fiber.Ctx) error {
	return models.GetDataTableTemplates(c)
}

func GetDataTableReviu(c *fiber.Ctx) error {
	return models.GetDataTableReviu(c)
}

func GetDataTablePanitia(c *fiber.Ctx) error {
	return models.GetDataTablePanitia(c)
}

func GetDataTablePp(c *fiber.Ctx) error {
	return models.GetDataTablePp(c)
}

func GetDataTableInbox(c *fiber.Ctx, id uint) error {
	return models.GetDataTableInbox(c, id)
}

func GetDataTableDocTemplate(c *fiber.Ctx) error {
	return models.GetDataTableDocTemplate(c)
}

func GetDataTableBukuTamu(c *fiber.Ctx, isUkpbj bool) error {
	return models.GetDataTableBukuTamu(c, isUkpbj)
}

func GetDataTableFeedback(c *fiber.Ctx) error {
	return models.GetDataTableFeedback(c)
}

func GetDataTableDocument(c *fiber.Ctx) error {
	return models.GetDataTableDocument(c)
}

func GetDataTableAdminDocument(c *fiber.Ctx) error {
	return models.GetDataTableAdminDocument(c)
}

func GetDataTableChecklist(c *fiber.Ctx) error {
	return models.GetDataTableChecklist(c)
}
