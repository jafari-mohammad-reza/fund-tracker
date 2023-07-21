package structs

type JsonResponse struct {
	StatusCode int              `json:"status_code"`
	Body       JsonResponseBody `json:"body"`
}
type JsonResponseBody struct {
	StatusCode int          `json:"status_code"`
	Success    bool         `json:"success"`
	data       *interface{} `json:"data"`
}
