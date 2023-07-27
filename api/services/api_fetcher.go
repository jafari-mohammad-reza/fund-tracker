package services

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	once       sync.Once
	httpClient *http.Client
)

type ApiFetcherService struct {
	httpClient *http.Client
}

func NewApiFetcher() *ApiFetcherService {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: 6 * time.Second,
		}
	})

	return &ApiFetcherService{
		httpClient: httpClient,
	}
}

type FetchResult struct {
	Result []byte
	Error  error
}

func (service *ApiFetcherService) FetchApiBytes(api string, headers *map[string]string) <-chan FetchResult {
	ch := make(chan FetchResult)

	go func() {
		defer close(ch)

		newUrl, err := url.ParseRequestURI(api)
		if err != nil {
			ch <- FetchResult{nil, err}
			return
		}

		req, err := http.NewRequest("GET", newUrl.String(), nil)

		if headers != nil {
			for key, value := range *headers {

				req.Header.Set(key, value)
			}
		}

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
