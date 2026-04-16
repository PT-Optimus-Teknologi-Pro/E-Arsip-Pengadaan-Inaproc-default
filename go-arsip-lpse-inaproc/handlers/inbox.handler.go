package handlers

import (

	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
)


func GetAllInbox(c *fiber.Ctx) error {
	mp := currentMap(c)
	return c.Render("inbox/inbox", mp)
}

func GetInbox(c *fiber.Ctx) error {
	mp := currentMap(c)
	id, _ := c.ParamsInt("id")
	inbox := services.GetInbox(uint(id))
	if inbox.ID == 0 {
		return c.SendStatus(404)
	}
	mp["inbox"] = inbox
	return c.Render("inbox/inbox-view", mp)
}

func GetJsonInbox(c *fiber.Ctx) error {
	mp := currentMap(c)
	id := utils.InterfaceToUint(mp["id"])
	return services.GetDataTableInbox(c, id)
}
