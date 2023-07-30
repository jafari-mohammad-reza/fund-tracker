package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/controllers"
)

func FundsRoute(router fiber.Router) {
	controller := controllers.NewFuncController()

	funds := router.Group("/funds")
	// Return all funds with compare date of 1 with ranking and complete data like the count of cancel and issues
	funds.Get("/", controller.GetFunds)
	funds.Get("/asset-chart/:regNo", controller.GetFundsIssueAndCancelData)
	// Return given regNo fund with cancel and issue count and efficiency chart and portfo data
	funds.Get("/info/:regNo", controller.GetFundInfo)
	//nav-per-year
	funds.Get("/nav-per-year", controller.GetNavPerYear)
	// funds efficiency chart in compare of market index
	funds.Get("/market-efficiency", controller.GetFundEfficiencyBaseOnMarket)

}

// TODO: add below features as well
