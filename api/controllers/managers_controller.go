package controllers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
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
	bgCtx := context.Background()
	fundQuery := GetQueryListQueries(ctx)
	managersList, err := controller.service.GetManagersListWithFunds(bgCtx, fundQuery)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch managers"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, managersList))
	return nil
}

func (controller *ManagersController) GetManagerInfo(ctx *fiber.Ctx) error {
	bgCtx := context.Background()
	query, err := GetManagerInfoQueryListQueries(ctx)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch managers"))
		return err
	}
	managerInfo, err := controller.service.GetManagerInfo(bgCtx, query)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch managers"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, managerInfo))
	return nil
}

func GetManagerInfoQueryListQueries(ctx *fiber.Ctx) (*dto.ManagerInfoQuery, error) {
	fundQuery := GetQueryListQueries(ctx)
	managerName := ctx.Query("managerName")
	if managerName == "" {
		return nil, errors.New("insert manager name")
	}
	return &dto.ManagerInfoQuery{
		ManagerName:   managerName,
		FundListQuery: fundQuery,
	}, nil
}
