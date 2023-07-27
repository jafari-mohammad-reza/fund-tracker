package structs

type ManagerListItems struct {
	NetAssetSum float64 `json:"netAssetSum"`
	AumDiffSum  float64 `json:"aumDiffSum"`
	RankDiffSum int     `json:"rankDiffSum"`
	Rank        int     `json:"rank"`
	IssuedCount int     `json:"issuedCount"`
	CancelCount int     `json:"cancelCount"`
}

type ManagerListResponse struct {
	Manager string            `json:"manager"`
	Funds   *[]CalculatedFund `json:"funds"`
	Items   ManagerListItems  `json:"Items"`
}
