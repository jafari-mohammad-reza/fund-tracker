package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
	"regexp"
	"strings"
	"time"
)

type ManagersService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
	fundService       *FundService
}

func NewManagersService() *ManagersService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	fundService := NewFundService()
	return &ManagersService{
		redisClient,
		apiFetcherService,
		fundService,
	}
}

func (service *ManagersService) GetManagersListWithFunds(ctx context.Context, queryList *dto.FundListQuery) (*[]structs.ManagerListResponse, error) {
	var response []structs.ManagerListResponse
	existData, err := data.GetValue(ctx, service.redisClient, "managers")
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if existData != "" {
		err := json.Unmarshal([]byte(existData), &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	calculatedFunds, err := service.fundService.GetFunds(queryList)
	if err != nil {
		return nil, err
	}
	managers := make(map[string]*structs.ManagerListResponse)

	wordsToRemove := []string{
		"تامین سرمایه",
		"سبدگردان",
		"سرمایه گذاری",
		"مشاور سرمایه گذاری",
		"کارگزاری",
		"تامین سرمایه",
		"مشاور",
		"سهامی خاص",
	}

	for _, fund := range *calculatedFunds {
		cleanedManager := fund.Manager
		for _, word := range wordsToRemove {
			cleanedManager = strings.Replace(cleanedManager, word, "", -1)
		}
		cleanedManager = strings.TrimSpace(cleanedManager) // Normalize the manager name

		if manager, exists := managers[cleanedManager]; exists {
			funds := append(*manager.Funds, fund)
			manager.Funds = &funds
			// Update ManagerListItems
			manager.Items.NetAssetSum += fund.NetAsset
			manager.Items.AumDiffSum += fund.NetAssetDiff
			manager.Items.RankDiffSum += fund.RankDiff
			manager.Items.Rank += fund.Rank
		} else {
			funds := []structs.CalculatedFund{fund}
			items := structs.ManagerListItems{
				NetAssetSum: fund.NetAsset,
				AumDiffSum:  fund.NetAssetDiff,
				RankDiffSum: fund.RankDiff,
				Rank:        fund.Rank,
			}
			managers[cleanedManager] = &structs.ManagerListResponse{
				Manager: cleanedManager,
				Funds:   &funds,
				Items:   items,
			}
		}
	}

	// Convert map to slice
	managerList := make([]structs.ManagerListResponse, 0, len(managers))
	for _, manager := range managers {
		// Compute averages
		numFunds := len(*manager.Funds)
		manager.Items.RankDiffSum /= numFunds
		manager.Items.Rank /= numFunds
		managerList = append(managerList, *manager)
	}
	dataJSON, err := json.Marshal(managerList)
	err = data.SetValue(ctx, service.redisClient, "managers", dataJSON, time.Hour*12)
	if err != nil {
		return nil, err
	}
	return &managerList, nil
}

func (service *ManagersService) GetManagerInfo(ctx context.Context, query *dto.ManagerInfoQuery) (*structs.ManagerInfoResponse, error) {
	managersList, err := service.GetManagersListWithFunds(ctx, query.FundListQuery)
	//var response structs.ManagerInfoResponse
	if err != nil {
		return nil, err
	}
	var manager structs.ManagerListResponse
	for _, mn := range *managersList {
		pattern := regexp.MustCompile(query.ManagerName)
		// Check if mn.Manager matches the pattern
		if pattern.MatchString(mn.Manager) {
			manager = mn
		}
	}

	if manager.Funds == nil {
		return nil, errors.New("No funds found for the manager")
	}
	issueAndCancelSum := structs.IssueAndCancelSum{
		UnitsSubDAYSum: 0,
		UnitsRedDAYSum: 0,
		Profit:         0,
	}

	for _, fund := range *manager.Funds {
		issueAndCancelData, err := service.fundService.fundInfoService.GetFundsIssueAndCancelData(query.FundListQuery.CompareDate, fund.RegNo)
		issueAndCancelDataSum, err := service.fundService.CalculateIssueAndCancelSum(issueAndCancelData, fund.IssueNav, fund.CancelNav)
		if err != nil {
			continue
		}
		issueAndCancelSum.UnitsSubDAYSum = issueAndCancelSum.UnitsSubDAYSum + issueAndCancelDataSum.UnitsSubDAYSum
		issueAndCancelSum.UnitsRedDAYSum = issueAndCancelSum.UnitsRedDAYSum + issueAndCancelDataSum.UnitsRedDAYSum
		issueAndCancelSum.Profit = issueAndCancelSum.Profit + issueAndCancelDataSum.Profit
	}

	return &structs.ManagerInfoResponse{
		ManagerListResponse: manager,
		IssueAndCancelSum:   issueAndCancelSum,
	}, nil
}
