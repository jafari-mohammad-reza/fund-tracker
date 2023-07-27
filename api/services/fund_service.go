package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"time"
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
	fundsChannel := service.apiFetcherService.FetchApiBytes(url, nil, nil)
	var responseData structs.FipIranResponse
	for res := range fundsChannel {
		if res.Error != nil {
			log.Printf("Error: %v\n", res.Error)
			return nil, res.Error
		}

		err := json.Unmarshal(res.Result, &responseData)
		if err != nil {
			log.Printf("Error unmarshalling data: %v\n", err)
			return nil, err
		}
	}

	return &responseData, nil
}
func findRank(items []structs.Fund, regNo string) int {
	for i, v := range items {
		if v.RegNo == regNo {
			return i
		}
	}
	return len(items)
}
func (service *FundService) GetFunds(queryList *dto.FundListQuery) (*[]structs.CalculatedFund, error) {

	baseUrl, err := url.Parse("https://fund.fipiran.ir/api/v1/fund/fundcompare")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	if err != nil {
		log.Println("Failed to parse URL: ", err.Error())
		return nil, err
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
	)

	if err != nil {
		return nil, err
	}
	var comparisionDate int
	if queryList.CompareDate != nil {
		comparisionDate = *queryList.CompareDate
	} else {
		comparisionDate = 7
	}
	date := url.QueryEscape(time.Now().AddDate(0, 0, -comparisionDate).Format("2006-01-02"))
	params := url.Values{}
	params.Add("date", date)
	baseUrl.RawQuery = params.Encode()
	previousDayResponseData, err := data.GetDataFromCacheOrFetch(
		fetchFuncWrapper,
		"fipiran-funds"+"-"+strconv.Itoa(*queryList.CompareDate)+"-"+*queryList.RankBy,
		ctx,
		service.redisClient,
	)
	go service.redisClient.Set(ctx, fmt.Sprintf("%s-%d", "fipiran-funds", queryList.CompareDate), previousDayResponseData, time.Hour*3)
	sortResponseDataItems(responseData.Items, *queryList.RankBy)
	sortResponseDataItems(previousDayResponseData.Items, *queryList.RankBy)
	if err != nil {
		return nil, err
	}
	previousDayFundsMap := make(map[string]int)
	for i, previousDayFunds := range previousDayResponseData.Items {
		previousDayFundsMap[previousDayFunds.RegNo] = i
	}

	// Calculate the funds and ranks using the map instead of linear search
	calculatedFunds := make([]structs.CalculatedFund, 0, len(responseData.Items))
	for fundsIndex, fund := range responseData.Items {
		previousDayIndex, ok := previousDayFundsMap[fund.RegNo]
		if ok {
			rank := previousDayIndex
			rankDiff := fundsIndex - rank
			netAssetDiff := math.Ceil(float64(fund.NetAsset - previousDayResponseData.Items[previousDayIndex].NetAsset))
			netAssetDiffPercent := float64((fund.NetAsset / previousDayResponseData.Items[previousDayIndex].NetAsset) * 100)
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

func sortResponseDataItems(responseData []structs.Fund, fieldName string) {
	sort.Slice(responseData, func(i, j int) bool {
		// Use reflection to get the field value based on the fieldName

		fieldI := reflect.ValueOf(responseData[i]).FieldByName(fieldName)
		fieldJ := reflect.ValueOf(responseData[j]).FieldByName(fieldName)

		// Compare the field values based on their types
		switch fieldI.Kind() {
		case reflect.String:
			return fieldI.String() < fieldJ.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldI.Int() < fieldJ.Int()
		// Add more cases for other types if needed
		default:
			// If the field type is not supported for comparison, you might want to handle it accordingly.
			// For example, return false to maintain the original order.
			return false
		}
	})
}
