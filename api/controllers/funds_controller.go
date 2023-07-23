package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"strconv"
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
	compareDateString := ctx.Query("compareDate")
	var compareDate int
	var err error
	if compareDateString != "" {
		compareDate, err = strconv.Atoi(compareDateString)
		if err != nil {
			ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "invalid compare date"))
			return err
		}
	} else {
		compareDate = 7
	}
	funds, err := controller.service.GetFunds(&compareDate)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch funds"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, funds))
	return nil

}

func (controller *FundController) GetManagers(ctx *fiber.Ctx) error {
	return nil
}

func (controller *FundController) GetFundInfo(ctx *fiber.Ctx) error {
	return nil
}
