package services

import (
	"github.com/jafari-mohammad-reza/fund-tracker/api/dto"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
	"strings"
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

func (service *ManagersService) GetManagersListWithFunds(queryList *dto.FundListQuery) (*[]structs.ManagerListResponse, error) {
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

	return &managerList, nil
}
