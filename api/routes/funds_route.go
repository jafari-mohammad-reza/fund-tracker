package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/controllers"
)

func FundsRoute(group *fiber.Group) {
	controller := controllers.NewFuncController()
	// Return all funds with compare date of 1
	group.Get("/", controller.GetFunds)
	// Return all funds managers with complete data
	group.Get("/managers", controller.GetManagers)
	// Return given regNo fund with cancel and issue count and efficiency chart and portfo data
	group.Get("/fund/:regNo", controller.GetFundInfo)
}
