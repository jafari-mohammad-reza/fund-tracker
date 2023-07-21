package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/controllers"
)

func FundsRoute(router fiber.Router) {
	controller := controllers.NewFuncController()
	funds := router.Group("/funds")
	// Return all funds with compare date of 1
	funds.Get("/", controller.GetFunds)
	// Return all funds managers with complete data
	funds.Get("/managers", controller.GetManagers)
	// Return given regNo fund with cancel and issue count and efficiency chart and portfo data
	funds.Get("/fund/:regNo", controller.GetFundInfo)
}
