package handlers

import (
	"arsip/services"
	"time"

	"github.com/gofiber/fiber/v2"
)



func GetJsonSatker(c *fiber.Ctx) error {
	tahun := c.QueryInt("tahun", time.Now().Year())
	satkers := services.GetSatkerAPI(tahun)
	return c.JSON(satkers)
}
