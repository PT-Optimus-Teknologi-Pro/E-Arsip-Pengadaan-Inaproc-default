package handlers

import "github.com/gofiber/fiber/v2"

// Get All Users from db
func GetAllPerencanaan(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("perencanaan/index", mp)
}
