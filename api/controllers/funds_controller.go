package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
)

type FundController struct {
	service *services.FundService
}

func NewFuncController() *FundController {
	fundService := services.NewFundService()
	return &FundController{
		service: fundService,
	}
}

func (controller *FundController) GetFunds(ctx *fiber.Ctx) error {
	return nil
}

func (controller *FundController) GetManagers(ctx *fiber.Ctx) error {
	return nil
}

func (controller *FundController) GetFundInfo(ctx *fiber.Ctx) error {
	return nil
}
