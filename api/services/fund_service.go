package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	baseURL    = "https://fund.fipiran.ir/api/v1/fund/fundcompare"
	refererURL = "https://fund.fipiran.ir/mf/profile"
)

type FundService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
	fundInfoService   *FundInfoService
}

func NewFundService() *FundService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	fundInfoService := NewFundInfoService()
	return &FundService{
		redisClient,
		apiFetcherService,
		fundInfoService,
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

func (service *FundService) CalculateIssueAndCancelSum(issueAndCancelData *[]structs.IssueAndCancelData, issueNav float64, cancelNav float64) (*structs.IssueAndCancelSum, error) {
	if issueNav == 0 {
		return nil, errors.New("issueNav cannot be zero")
	}
	if cancelNav == 0 {
		return nil, errors.New("cancelNav cannot be zero")
	}
	var issueCountSum int
	var cancelCountSum int
	for _, dt := range *issueAndCancelData {
		issueCountSum = issueCountSum + dt.UnitsSubDAY
		cancelCountSum = cancelCountSum + dt.UnitsRedDAY
	}
	issueValue := float64(issueCountSum) * issueNav
	cancelValue := float64(cancelCountSum) * cancelNav
	issueAndCancelSum := structs.IssueAndCancelSum{
		UnitsSubDAYSum: issueCountSum,
		UnitsRedDAYSum: cancelCountSum,
		Profit:         issueValue - cancelValue,
	}
	return &issueAndCancelSum, nil
}

func getComparisonFunds(service *FundService, queryList *dto.FundListQuery) (currentDateFunds *[]structs.Fund, compareDateFunds *[]structs.Fund, err error) {
	baseUrl, err := url.Parse(baseURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err != nil {
		log.Println("Failed to parse URL: ", err.Error())
		cancel()
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
	if err != nil {
		log.Printf("Error fetching previous day's data: %v\n", err)
	}

	if responseData == nil || previousDayResponseData == nil {
		return nil, nil, errors.New("Can not sort empty data")
	}
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

	var calculatedFunds []structs.CalculatedFund
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
			})
		}
	}
	return &calculatedFunds, nil
}

func (service *FundService) GetNavPerYear() {

}

func (service *FundService) GetEachYearFunds() (*[]structs.EachYearFunds, error) {
	startYear := 2008
	currentYear := time.Now().Year()
	yearDiff := (currentYear - startYear) + 1

	eachYearDataMap := make(map[int][]structs.Fund, yearDiff)
	eachYearFunds := make([]structs.EachYearFunds, 0, yearDiff) // Initialize with length 0, capacity yearDiff

	// Create a channel to communicate the results from the goroutines
	resChan := make(chan structs.EachYearFunds, yearDiff)
	// Create a channel to communicate any errors from the goroutines
	errChan := make(chan error, yearDiff)
	// Create a WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup

	for i := 0; i < yearDiff; i++ {
		year := startYear + i

		wg.Add(1)
		go func(year int) {
			defer wg.Done()

			data, err := service.fetchFunds(fmt.Sprintf("%s?%s=%s", baseURL, "date", year))
			if err != nil {
				errChan <- err
				return
			}
			res := structs.EachYearFunds{
				Year:  year,
				Funds: data.Items,
			}
			resChan <- res
		}(year)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	// Close the channels
	close(resChan)
	close(errChan)

	// Check for any errors
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	// Process results
	for res := range resChan {
		eachYearDataMap[res.Year] = append(eachYearDataMap[res.Year], res.Funds...)
	}

	for year, funds := range eachYearDataMap {
		yearFunds := structs.EachYearFunds{
			Year:  year,
			Funds: funds,
		}
		eachYearFunds = append(eachYearFunds, yearFunds)
	}

	return &eachYearFunds, nil
}

func (service *FundService) CalculateEachYearTotalNav() {

}
