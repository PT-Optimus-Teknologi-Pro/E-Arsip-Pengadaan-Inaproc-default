package main

import (
	"arsip/config"
	_ "arsip/config"
	_ "arsip/cache"
	"arsip/handlers"
	_ "arsip/models"
	"arsip/routers"
	"arsip/services"
	"arsip/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/django/v3"
	"github.com/robfig/cron/v3"
)

func main() {
	services.SyncSirup()
	services.AutoCreateAdminIfNoExist();
	c := cron.New()
	cronjob := config.CronJob()
	c.AddFunc(cronjob, func() {
		services.SyncSirup()
	})
	c.Start()
	engine := django.New("./views", ".html")
	engine.SetAutoEscape(false)
	utils.Setup(engine)

	if config.IsModeDev(){
		log.Info("application mode development")
		engine.Reload(true)
	} else {
		log.Info("application mode production")
	}
	app := fiber.New(fiber.Config{
		Views: engine,
		StreamRequestBody: true,
	})
	handlers.Sessions = session.New()
	routers.SetupRoutes(app)
	app.Listen(":" + config.Port())

}
