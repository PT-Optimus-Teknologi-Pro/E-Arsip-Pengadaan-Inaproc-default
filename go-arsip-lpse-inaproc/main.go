package main

import (
	"arsip/config"
	_ "arsip/config"
	_ "arsip/cache"
	"arsip/handlers"
	"arsip/models"
	"arsip/routers"
	"arsip/services"
	"arsip/utils"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/django/v3"
	"github.com/robfig/cron/v3"
)

func main() {
	go services.SyncSirup()
	go services.SyncLpse()
	services.AutoCreateAdminIfNoExist()
	c := cron.New()
	cronjob := config.CronJob()
	c.AddFunc(cronjob, func() {
		services.SyncSirup()
		services.SyncLpse()
	})
	c.Start()
	engine := django.New("./views", ".html")
	engine.SetAutoEscape(false)
	utils.Setup(engine)
	engine.AddFunc("metode", func(id interface{}) string {
		return models.GetMetodeLabel(id)
	})
	engine.AddFunc("GetFooterSocials", handlers.GetFooterSocials)
	engine.AddFunc("GetFooterQuicks", handlers.GetFooterQuicks)
	engine.AddFunc("GetFooterServices", handlers.GetFooterServices)
	engine.AddFunc("parseApprovals", func(s string) string {
		type ApprovalState struct {
			Name   string `json:"name"`
			Status bool   `json:"status"`
		}
		var state []ApprovalState
		json.Unmarshal([]byte(s), &state)
		var res []string
		for _, v := range state {
			status := "❌"
			if v.Status {
				status = "✅"
			}
			res = append(res, fmt.Sprintf("%s %s", v.Name, status))
		}
		return strings.Join(res, ", ")
	})

	if config.IsModeDev(){
		log.Info("application mode development")
		engine.Reload(true)	
	} else {
		log.Info("application mode production")
	}
	app := fiber.New(fiber.Config{
		Views: engine,
		StreamRequestBody: false,
	})
	handlers.Sessions = session.New()
	routers.SetupRoutes(app)
	app.Listen(":" + config.Port())

}
