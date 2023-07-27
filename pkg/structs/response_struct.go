package structs

type JsonResponse struct {
	StatusCode int              `json:"status_code"`
	Body       JsonResponseBody `json:"body"`
}
type JsonResponseBody struct {
	StatusCode int  `json:"status_code"`
	Success    bool `json:"success"`
	Data       any  `json:"data"`
}

func NewJsonResponse(status int, success bool, data any) *JsonResponse {
	return &JsonResponse{
		StatusCode: status,
		Body: JsonResponseBody{
			StatusCode: status,
			Success:    success,
			Data:       data,
		},
	}
}
