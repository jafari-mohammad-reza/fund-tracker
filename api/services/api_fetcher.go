package services

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiFetcherService struct {
	httpClient *http.Client
}

func NewApiFetcher() *ApiFetcherService {
	httpClient := &http.Client{
		Timeout: 6 * time.Second,
	}
	return &ApiFetcherService{
		httpClient: httpClient,
	}
}

type FetchResult struct {
	Result []byte
	Error  error
}

func (service *ApiFetcherService) FetchApiBytes(api string, params *string, queries *string) <-chan FetchResult {
	ch := make(chan FetchResult)

	go func() {
		defer close(ch)

		newUrl, err := url.ParseRequestURI(api)
		if err != nil {
			ch <- FetchResult{nil, err}
			return
		}
		if params != nil {
			newUrl.JoinPath(*params)
		}
		if queries != nil {
			queryList := strings.Split(*queries, "&")
			for _, queryItem := range queryList {
				query := strings.Split(queryItem, "=")
				newUrl.Query().Set(query[0], query[1])
			}
		}

		req, err := http.NewRequest("GET", newUrl.RequestURI(), nil)
		if err != nil {
			ch <- FetchResult{nil, err}
			return
		}
		resp, err := service.httpClient.Do(req)
		if err != nil {
			ch <- FetchResult{nil, err}
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ch <- FetchResult{nil, err}
			return
		}

		ch <- FetchResult{body, nil}
		defer resp.Body.Close()

	}()

	return ch
}
