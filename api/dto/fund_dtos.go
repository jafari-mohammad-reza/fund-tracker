package dto

type FundListQuery struct {
	CompareDate *int    `json:"compareDate"`
	RankBy      *string `json:"rankBy"`
	FundType    *[]int  `json:"fundType"`
}
