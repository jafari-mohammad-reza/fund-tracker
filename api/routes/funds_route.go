package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/controllers"
)

func FundsRoute(router fiber.Router) {
	controller := controllers.NewFuncController()

	funds := router.Group("/funds")

	funds.Get("/", controller.GetFunds)
	funds.Get("/asset-chart/:regNo", controller.GetFundsIssueAndCancelData)
	// Return given regNo fund with cancel and issue count and efficiency chart and portfo data
	funds.Get("/info/:regNo", controller.GetFundInfo)
	//nav-per-year
	funds.Get("/nav-per-year", controller.GetNavPerYear)
}
