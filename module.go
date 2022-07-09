package main

import (
	"github.com/gofiber/fiber/v2"
	"l.hilmy.dev/backend/modules/url"
)

type module struct {
	appName *string
	app     *fiber.App
}

func (m *module) run() {
	urlRouterGroup := m.app.Group("/")
	urlModule := url.New(&urlRouterGroup)
	urlModule.Run()
}
