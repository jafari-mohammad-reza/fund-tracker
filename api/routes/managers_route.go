package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/controllers"
)

func ManagersRouter(router fiber.Router) {
	controller := controllers.NewManagersController()
	managers := router.Group("/managers")
	managers.Get("/", controller.GetManagersList)
}
