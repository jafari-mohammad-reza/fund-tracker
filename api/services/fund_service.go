package services

import (
	"encoding/json"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"net/url"
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
func (service *FundService) GetFunds() (*[]structs.CalculatedFund, error) {
	baseUrl, err := url.Parse("https://fund.fipiran.ir/api/v1/fund/fundcompare")
	if err != nil {
		log.Println("Failed to parse URL: ", err.Error())
		return nil, err
	}

	responseData, err := service.fetchFunds(baseUrl.String())
	if err != nil {
		return nil, err
	}

	date := url.QueryEscape(time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	params := url.Values{}
	params.Add("date", date)
	baseUrl.RawQuery = params.Encode()

	previousDayResponseData, err := service.fetchFunds(baseUrl.String())
	if err != nil {
		return nil, err
	}
	previousDayFundsMap := make(map[string]structs.Fund)
	for _, previousDayFunds := range previousDayResponseData.Items {
		previousDayFundsMap[previousDayFunds.RegNo] = previousDayFunds
	}

	var calculatedFunds []structs.CalculatedFund
	for fundsIndex, fund := range responseData.Items {
		previousDayFunds, ok := previousDayFundsMap[fund.RegNo]
		if ok {
			rank := findRank(previousDayResponseData.Items, fund.RegNo)
			rankDiff := fundsIndex - rank
			netAssetDiff := math.Ceil(float64(fund.NetAsset - previousDayFunds.NetAsset))
			netAssetDiffPercent := math.Ceil(float64((fund.NetAsset - previousDayFunds.NetAsset) / 100))
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
