package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/api/services"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
)

type FundController struct {
	fundService     *services.FundService
	fundInfoService *services.FundInfoService
}

func NewFuncController() *FundController {
	fundService := services.NewFundService()
	fundInfoService := services.NewFundInfoService()
	return &FundController{
		fundService,
		fundInfoService,
	}
}

// GetFunds godoc
// @Summary Get all funds
// @Description get all funds with compare date of 1 with ranking and complete data like the count of cancel and issues
// @Tags funds
// @Accept */*
// @Produce json
// @Param compareDate query int false "Comparison date for funds data"
// @Param rankBy query string false "Ranking criteria for funds data"
// @Success 200 {object} map[string]interface{}
// @Router /funds/ [get]
func (controller *FundController) GetFunds(ctx *fiber.Ctx) error {
	queryList := GetQueryListQueries(ctx)
	funds, err := controller.fundService.GetFunds(queryList)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch funds"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, funds))
	return nil

}

// GetFundsIssueAndCancelData godoc
// @Summary Get asset chart data for a fund
// @Description get given regNo fund with cancel and issue count and efficiency chart and portfo data
// @Tags funds
// @Accept */*
// @Produce json
// @Param regNo path int true "Fund Registration Number"
// @Param compareDate query int false "Comparison date for funds data"
// @Param rankBy query string false "Ranking criteria for funds data"
// @Success 200 {object} map[string]interface{}
// @Router /funds/asset-chart/{regNo} [get]
func (controller FundController) GetFundsIssueAndCancelData(ctx *fiber.Ctx) error {
	regNo := ctx.Params("regNo")
	queryList := GetQueryListQueries(ctx)
	if regNo == "" {
		ctx.Status(400).JSON(structs.NewJsonResponse(400, false, "Insert regNo"))
	}
	issueAndCancelData, err := controller.fundInfoService.GetFundsIssueAndCancelData(queryList.CompareDate, regNo)
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

// GetFundInfo godoc
// @Summary Get information for a specific fund
// @Description get detailed information for a specific fund
// @Tags funds
// @Accept */*
// @Produce json
// @Param regNo path int true "Fund Registration Number"
// @Success 200 {object} map[string]interface{}
// @Router /funds/info/{regNo} [get]
func (controller *FundController) GetFundInfo(ctx *fiber.Ctx) error {
	regNo := ctx.Params("regNo")
	if regNo == "" {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "invalid regno"))
		return errors.New("invalid regno")
	}
	info, err := controller.fundInfoService.GetFundInfo(regNo)
	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch fund info"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, info))
	return nil
}

// GetNavPerYear godoc
// @Summary Get list of each year with that year nav
// @Description Get list of each year with that year nav
// @Tags funds
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /funds/nav-per-year [get]

func (controller *FundController) GetNavPerYear(ctx *fiber.Ctx) error {
	data, err := controller.fundService.CalculateEachYearTotalNav()

	if err != nil {
		ctx.Status(500).JSON(structs.NewJsonResponse(500, false, "failed to fetch nav per year"))
		return err
	}
	ctx.Status(200).JSON(structs.NewJsonResponse(200, true, data))
	return nil
}
