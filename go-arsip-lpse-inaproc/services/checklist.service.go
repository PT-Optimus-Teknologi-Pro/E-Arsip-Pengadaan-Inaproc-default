package services

import (
	"arsip/models"
)

func GetJenisPengadaan() []string {
	return models.GetAllJenisPengadaan()
}

func GetMetodePengadaan() []string {
	return models.GetAllMetodePengadaan()
}

func GetChecklists() []models.Checklist {
	return models.GetChecklists()
}

func GetChecklistsBYjenis(jenis int) []models.ChecklistDok {
	return models.GetChecklistsBYjenis(jenis)
}

func GetChecklist(id uint) models.Checklist {
	return models.GetChecklist(id)
}

func HapusChecklist(id uint) error {
	return models.HapusChecklist(id)
}

func SimpanChecklist(checklist models.Checklist) error {
	return models.SaveChecklist(&checklist)
}

func SimpanChecklistPersyaratan(checklist []models.ChecklistDok) error {
	return models.SimpanChecklist(checklist)
}

func GetTemplates(id uint) models.Templates {
	return models.GetTemplates(id)
}

func CreateTemplates(template models.Templates) error {
	return models.CreateTemplates(template)
}

func SaveTemplates(template models.Templates) error {
	return models.SaveTemplates(template)
}

func DeleteTemplates(template models.Templates) error {
	return models.DeleteTemplates(template)
}
