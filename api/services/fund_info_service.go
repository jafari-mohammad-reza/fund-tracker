package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
	"net/url"
)

const (
	fundAssetChartUrl = "https://fund.fipiran.ir/api/v1/chart/getfundnetassetchart"
	fundPortfoUrl     = "https://fund.fipiran.ir/api/v1/chart/portfoliochart"
	fundEfficiencyUrl = "https://fund.fipiran.ir/api/v1/chart/fundefficiencychart"
	fundBasicInfoUrl  = "https://fund.fipiran.ir/api/v1/fund/getfund"
)

type FundInfoService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
}

func NewFundInfoService() *FundInfoService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	return &FundInfoService{redisClient, apiFetcherService}
}
func (service *FundInfoService) GetFundsIssueAndCancelData(comparisonDays *int, regNo string) (issueAndCancel *[]structs.IssueAndCancelData, err error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s?regno=%s", fundAssetChartUrl, regNo))
	headers := make(map[string]string)
	headers["Referer"] = fmt.Sprintf("%s/%s", refererURL, regNo)

	response := service.apiFetcherService.FetchApiBytes(baseUrl.String(), &headers)
	var issueAndCancelData []structs.IssueAndCancelData
	for res := range response {

		if res.Error != nil {
			return nil, res.Error
		}

		err := json.NewDecoder(bytes.NewBuffer(res.Result)).Decode(&issueAndCancelData)
		if err != nil {
			return nil, err
		}
	}
	if comparisonDays != nil {
		slicedData := issueAndCancelData[:*comparisonDays]
		return &slicedData, nil
	}
	return &issueAndCancelData, nil
}

func (service *FundInfoService) GetFundPortfo(regNo int) (*[]structs.FundPortfolio, error) {
	return nil, nil
}
func (service *FundInfoService) GetFundEfficiency(regNo int) (*[]structs.FundEfficiency, error) {
	return nil, nil
}

func (service *FundInfoService) GetFundBasicInfo(regNo int) (*structs.Fund, error) {
	return nil, nil
}

func (service *FundInfoService) GetFundInfo(date *string) (*structs.FundInfo, error) {
	return nil, nil
}
