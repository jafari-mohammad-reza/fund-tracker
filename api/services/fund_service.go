package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/utils"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL           = "https://fund.fipiran.ir/api/v1/fund/fundcompare"
	fundAssetChartUrl = "https://fund.fipiran.ir/api/v1/chart/getfundnetassetchart"
	refererURL        = "https://fund.fipiran.ir/mf/profile"
)

type FundService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
}

func NewFundService() *FundService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	return &FundService{
		redisClient:       redisClient,
		apiFetcherService: apiFetcherService,
	}
}
func (service *FundService) fetchFunds(url string) (*structs.FipIranResponse, error) {
	fundsChannel := service.apiFetcherService.FetchApiBytes(url, nil)
	var responseData structs.FipIranResponse
	for res := range fundsChannel {
		if res.Error != nil {
			log.Printf("Error: %v\n", res.Error)
			return nil, res.Error
		}

		err := json.NewDecoder(bytes.NewBuffer(res.Result)).Decode(&responseData)

		if err != nil {
			log.Printf("Error unmarshalling data: %v\n", err)
			return nil, err
		}
	}

	return &responseData, nil
}

func (service *FundService) GetFundsIssueAndCancelData(comparisonDays *int, regNo string) (issueAndCancel *[]structs.IssueAndCancelData, err error) {
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

	return &issueAndCancelData, nil
}
func getComparisonFunds(service *FundService, queryList *dto.FundListQuery) (currentDateFunds *[]structs.Fund, compareDateFunds *[]structs.Fund, err error) {
	baseUrl, err := url.Parse(baseURL)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	if err != nil {
		log.Println("Failed to parse URL: ", err.Error())
		return nil, nil, err
	}
	fetchFuncWrapper := func() (*structs.FipIranResponse, error) {
		responseData, err := service.fetchFunds(baseUrl.String())
		if err != nil {
			return nil, err
		}
		return responseData, nil
	}

	responseData, err := data.GetDataFromCacheOrFetch(
		fetchFuncWrapper,
		"fipiran-funds",
		ctx,
		service.redisClient,
		time.Hour*3,
	)
	if err != nil {
		return nil, nil, err
	}
	var comparisonDate int
	if queryList.CompareDate != nil {
		comparisonDate = *queryList.CompareDate
	} else {
		comparisonDate = 7
	}
	date := url.QueryEscape(time.Now().AddDate(0, 0, -comparisonDate).Format("2006-01-02"))
	params := url.Values{}
	params.Add("date", date)
	baseUrl.RawQuery = params.Encode()
	previousDayResponseData, err := data.GetDataFromCacheOrFetch(
		fetchFuncWrapper,
		"fipiran-funds"+"-"+strconv.Itoa(*queryList.CompareDate)+"-"+*queryList.RankBy,
		ctx,
		service.redisClient,
		time.Hour*3,
	)
	utils.SortResponseDataItems(responseData.Items, *queryList.RankBy)
	utils.SortResponseDataItems(previousDayResponseData.Items, *queryList.RankBy)
	if err != nil {
		return nil, nil, err
	}
	return &responseData.Items, &previousDayResponseData.Items, nil
}
func (service *FundService) GetFunds(queryList *dto.FundListQuery) (*[]structs.CalculatedFund, error) {

	responseData, previousDayResponseData, err := getComparisonFunds(service, queryList)
	if err != nil {
		return nil, err
	}
	previousDayFundsMap := make(map[string]int)
	for i, previousDayFunds := range *previousDayResponseData {
		previousDayFundsMap[previousDayFunds.RegNo] = i
	}

	calculatedFunds := make([]structs.CalculatedFund, 0, len(*responseData))
	for fundsIndex, fund := range *responseData {
		previousDayIndex, ok := previousDayFundsMap[fund.RegNo]
		if ok {
			rank := previousDayIndex
			rankDiff := fundsIndex - rank
			netAssetDiff := math.Ceil(float64(fund.NetAsset - (*previousDayResponseData)[previousDayIndex].NetAsset))
			netAssetDiffPercent := float64((fund.NetAsset / (*previousDayResponseData)[previousDayIndex].NetAsset) * 100)
			calculatedFunds = append(calculatedFunds, structs.CalculatedFund{
				Fund:                fund,
				Rank:                rank,
				RankDiff:            rankDiff,
				NetAssetDiff:        netAssetDiff,
				NetAssetDiffPercent: netAssetDiffPercent,
				//ProfitDiff: fund.CancelNav /
			})
		}
	}

	return &calculatedFunds, nil
}
