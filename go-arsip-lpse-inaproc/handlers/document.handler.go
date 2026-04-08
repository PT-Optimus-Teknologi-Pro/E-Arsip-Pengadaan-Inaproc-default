package handlers

import (
	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
)

func GetAllDocument(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("document/document", mp)
}

func GetDocument(c *fiber.Ctx) error {
	// get id params
	id := utils.StringToUint(c.Params("id"))
	document := services.GetDocument(id)
	if document.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Document not found", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Document Found", "data": document})
}

func DeleteDocument(c *fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	document := services.GetDocument(id)
	if document.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	err := services.DeleteDocument(document)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User deleted"})
}

func GetJsonDocument(c *fiber.Ctx) error {
	return services.GetDataTableDocument(c)
}
