package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sotowang/sunotoapi/cfg"
	"github.com/sotowang/sunotoapi/router"
	"github.com/sotowang/sunotoapi/serve"
)

func init() {
	cfg.ConfigInit()
	serve.Session = serve.GetSession(cfg.Config.App.Client)
}

func main() {
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})
	router.SetupRoutes(app)
	app.Listen(":" + cfg.Config.Server.Port)
}
