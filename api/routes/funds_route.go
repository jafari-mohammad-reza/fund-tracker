package routes

import "github.com/gofiber/fiber/v2"

func FundsRoute(group *fiber.Group) {
	// Return all funds with compare date of 1
	group.Get("/")
	// Return all funds managers with complete data
	group.Get("/managers")
	// Return given regNo fund with cancel and issue count and efficiency chart and portfo data
	group.Get("/fund/:regNo")
}
