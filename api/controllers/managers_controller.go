package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
)

type ManagersController struct {
	service *services.ManagersService
}

func NewManagersController() *ManagersController {
	service := services.NewManagersService()
	return &ManagersController{
		service,
	}
}

func (controller *ManagersController) GetManagersList(ctx *fiber.Ctx) error {
	return nil
}
