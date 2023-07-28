package dto

type ManagerInfoQuery struct {
	ManagerName   string         `json:"managerName"`
	FundListQuery *FundListQuery `json:"fundListQuery"`
}
