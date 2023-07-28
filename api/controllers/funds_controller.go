package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"strconv"
	"strings"
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
	queryList := GetQueryListQueries(ctx)
	funds, err := controller.service.GetFunds(queryList)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch funds"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, funds))
	return nil

}

func (controller FundController) GetFundsIssueAndCancelData(ctx *fiber.Ctx) error {
	regNo := ctx.Params("regNo")
	queryList := GetQueryListQueries(ctx)
	if regNo == "" {
		ctx.Status(400).JSON(structs.NewJsonResponse(400, false, "Insert regNo"))
	}
	issueAndCancelData, err := controller.service.GetFundsIssueAndCancelData(queryList.CompareDate, regNo)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch fund issue and cancel data"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, issueAndCancelData))

	return nil
}

func GetQueryListQueries(ctx *fiber.Ctx) *dto.FundListQuery {
	queryString := string(ctx.Request().URI().QueryString())
	queryList := make(map[string]string)

	for _, queryItem := range strings.Split(queryString, "&") {
		query := strings.Split(queryItem, "=")
		if len(query) == 2 {
			queryList[query[0]] = query[1]
		}
	}

	// Now, create a FundListQuery instance and populate its fields
	fundListQuery := dto.FundListQuery{}

	if compareDate, ok := queryList["compareDate"]; ok {
		compareDateInt, err := strconv.Atoi(compareDate)
		if err == nil {
			fundListQuery.CompareDate = &compareDateInt
		}
	} else {
		// Set default value for CompareDate (7 in this case)
		defaultValue := 7
		fundListQuery.CompareDate = &defaultValue
	}

	if rankBy, ok := queryList["rankBy"]; ok {
		fundListQuery.RankBy = &rankBy
	} else {
		// Set default value for RankBy (empty string in this case)
		defaultValue := "monthlyEfficiency"
		fundListQuery.RankBy = &defaultValue
	}

	return &fundListQuery
}

func (controller *FundController) GetManagers(ctx *fiber.Ctx) error {
	return nil
}

func (controller *FundController) GetFundInfo(ctx *fiber.Ctx) error {
	return nil
}
