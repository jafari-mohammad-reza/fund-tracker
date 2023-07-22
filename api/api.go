package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jafari-mohammad-reza/fund-tracker/api/routes"
	"os"
	"time"
)

func NewServer() {
	var printRoute bool
	if os.Getenv("Stage") == "development" {
		printRoute = true
	} else {
		printRoute = false
	}
	app := fiber.New(fiber.Config{
		Prefork:           true,
		AppName:           "Fund tracker v1.0.0",
		EnablePrintRoutes: printRoute,
		ReadTimeout:       time.Second * 10,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})
	setupRoutes(app)
	setupSwagger(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Route not available!")
	})

	if err := app.Listen(":5000"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(200)
	})
	routes.FundsRoute(v1)
}

// @title Fund Tracker
// @version 1.0
// @description fund tracker api documentation
// @contact.name mohammadreza jafari
// @contact.email mohammadrezajafari.dev@gmail.com
// @host localhost:5000
// @BasePath /
func setupSwagger(app *fiber.App) {
	app.Get("/api-docs/*", swagger.New(swagger.Config{
		InstanceName: "Fund Tracker",
		Title:        "Fund Tracker",
		URL:          os.Getenv("APP_URL"),
	}))
}
