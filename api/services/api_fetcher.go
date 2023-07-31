package services

import (
	"bytes"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"io"
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
			Timeout: 30 * time.Second,
		}
	})

	return &ApiFetcherService{
		httpClient: httpClient,
	}
}
func (service *ApiFetcherService) doRequest(method string, api string, headers *map[string]string, body io.Reader) <-chan structs.ApiFetchResult {
	ch := make(chan structs.ApiFetchResult)

	go func() {
		defer close(ch)

		newUrl, err := url.ParseRequestURI(api)
		if err != nil {
			ch <- structs.ApiFetchResult{Error: err}
			return
		}

		req, err := http.NewRequest(method, newUrl.String(), body)

		if headers != nil {
			for key, value := range *headers {
				req.Header.Set(key, value)
			}
		}

		if err != nil {
			ch <- structs.ApiFetchResult{Error: err}
			return
		}
		resp, err := service.httpClient.Do(req)
		if err != nil {
			ch <- structs.ApiFetchResult{Error: err}
			return
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- structs.ApiFetchResult{Error: err}
			return
		}

		ch <- structs.ApiFetchResult{Result: respBody}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

	}()

	return ch
}

func (service *ApiFetcherService) FetchApiBytes(api string, headers *map[string]string) <-chan structs.ApiFetchResult {
	return service.doRequest("GET", api, headers, nil)
}

func (service *ApiFetcherService) PostMultipartRequest(api string, headers *map[string]string, body *bytes.Buffer) <-chan structs.ApiFetchResult {
	return service.doRequest("POST", api, headers, body)
}
