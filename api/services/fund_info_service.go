package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"
)

const (
	fundAssetChartUrl = "https://fund.fipiran.ir/api/v1/chart/getfundnetassetchart"
	fundPortfoUrl     = "https://fund.fipiran.ir/api/v1/chart/portfoliochart"
	fundEfficiencyUrl = "https://fund.fipiran.ir/api/v1/chart/fundefficiencychart"
	fundBasicInfoUrl  = "https://fund.fipiran.ir/api/v1/fund/getfund"
	marketIndexUrl    = "https://fipiran.ir/DataService/Exportindex"
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
	if comparisonDays != nil && *comparisonDays != -1 {
		slicedData := issueAndCancelData[:*comparisonDays]
		return &slicedData, nil
	}
	return &issueAndCancelData, nil
}

func (service *FundInfoService) GetFundPortfo(regNo string) (*[]structs.FundPortfolio, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s?regno=%s", fundPortfoUrl, regNo))
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	headers["Referer"] = fmt.Sprintf("%s/%s", refererURL, regNo)

	response := service.apiFetcherService.FetchApiBytes(baseUrl.String(), &headers)
	var fundPortfo []structs.FundPortfolio
	for res := range response {

		if res.Error != nil {
			return nil, res.Error
		}

		err := json.NewDecoder(bytes.NewBuffer(res.Result)).Decode(&fundPortfo)
		if err != nil {
			return nil, err
		}
	}

	return &fundPortfo, nil
}
func (service *FundInfoService) GetFundEfficiency(regNo string) (*[]structs.FundEfficiency, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s?regno=%s", fundEfficiencyUrl, regNo))
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	headers["Referer"] = fmt.Sprintf("%s/%s", refererURL, regNo)

	response := service.apiFetcherService.FetchApiBytes(baseUrl.String(), &headers)
	var fundEfficiency []structs.FundEfficiency
	for res := range response {

		if res.Error != nil {
			return nil, res.Error
		}

		err := json.NewDecoder(bytes.NewBuffer(res.Result)).Decode(&fundEfficiency)
		if err != nil {
			return nil, err
		}
	}

	return &fundEfficiency, nil
}

func (service *FundInfoService) GetFundBasicInfo(regNo string) (*structs.FundBasicInfo, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s?regno=%s", fundBasicInfoUrl, regNo))
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	headers["Referer"] = fmt.Sprintf("%s/%s", refererURL, regNo)

	response := service.apiFetcherService.FetchApiBytes(baseUrl.String(), &headers)
	var fundBasicInfo structs.FundBasicInfo
	for res := range response {

		if res.Error != nil {
			return nil, res.Error
		}

		err := json.NewDecoder(bytes.NewBuffer(res.Result)).Decode(&fundBasicInfo)
		if err != nil {
			return nil, err
		}
	}

	return &fundBasicInfo, nil
}

func (service *FundInfoService) GetFundInfo(regNo string) (*structs.FundInfo, error) {
	fundPortfoCh := make(chan *[]structs.FundPortfolio, 1)
	fundEfficiencyCh := make(chan *[]structs.FundEfficiency, 1)
	fundBasicInfoCh := make(chan *structs.FundBasicInfo, 1)
	errorCh := make(chan error, 3)

	// define a comparisonDays value
	comparisonDays := new(int)
	*comparisonDays = -1

	go func() {
		res, err := service.GetFundPortfo(regNo)
		if err != nil {
			errorCh <- err
			return
		}
		fundPortfoCh <- res
	}()

	go func() {
		res, err := service.GetFundEfficiency(regNo)
		if err != nil {
			errorCh <- err
			return
		}
		fundEfficiencyCh <- res
	}()

	go func() {
		res, err := service.GetFundBasicInfo(regNo)
		if err != nil {
			errorCh <- err
			return
		}
		fundBasicInfoCh <- res
	}()

	// create a FundInfo struct to store the results
	fundInfo := &structs.FundInfo{}

	// wait for all the goroutines to finish
	for i := 0; i < 3; i++ {
		select {
		case res := <-fundPortfoCh:
			fundInfo.FundPortfolio = *res
		case res := <-fundEfficiencyCh:
			fundInfo.FundEfficiency = *res
		case res := <-fundBasicInfoCh:
			fundInfo.Fund = res.Item
		case err := <-errorCh:
			return nil, err
		}
	}

	return fundInfo, nil
}

func (service *FundInfoService) GetMarketIndexPerYear() (*[]structs.TransformedMarketIndex, error) {
	body, contentType, err := getMarketIndexBody()
	fmt.Println(body)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string, 1)
	headers["Content-Type"] = *contentType
	response := service.apiFetcherService.PostMultipartRequest(marketIndexUrl, &headers, body)
	var result []byte
	for res := range response {
		if res.Error != nil {
			return nil, res.Error
		}
		result = res.Result
	}
	fmt.Println(string(result))
	err = ioutil.WriteFile("test.xlsx", result, 0644)
	if err != nil {
		log.Fatal(err)
	}

	workbook, err := xlsx.OpenBinary(result)
	if err != nil {
		return nil, err
	}

	sheet := workbook.Sheets[0]
	items := make([]structs.MarketIndex, 0, len(sheet.Rows))
	for _, row := range sheet.Rows {
		var item structs.MarketIndex
		err = row.ReadStruct(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	transformedDataList := make([]structs.TransformedMarketIndex, 0, len(items))
	for _, data := range items {
		transformedData := structs.TransformedMarketIndex{
			Date:  utils.ShamsiStringToGeoDate(strconv.Itoa(data.DateIssue)),
			Value: data.Value,
		}
		transformedDataList = append(transformedDataList, transformedData)
	}
	return &transformedDataList, nil
}

func getMarketIndexBody() (*bytes.Buffer, *string, error) {
	body := &bytes.Buffer{}
	indexEnd, err := utils.CurrentShamsiDateNum()
	if err != nil {
		return nil, nil, err
	}
	writer := multipart.NewWriter(body)
	err = writer.WriteField("indexstartDate", "13870914")
	if err != nil {
		return nil, nil, err
	}
	err = writer.WriteField("indexEnd", *indexEnd)
	if err != nil {
		return nil, nil, err
	}
	err = writer.WriteField("indexpara", "شاخص+كل")
	if err != nil {
		return nil, nil, err
	}
	err = writer.WriteField("inscodeindex", "IRX6XTPI0006")
	if err != nil {
		return nil, nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	contentType := writer.FormDataContentType()
	fmt.Println(contentType)
	return body, &contentType, nil
}
