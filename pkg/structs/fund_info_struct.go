package structs

type FundPortfolio struct {
	Date      string  `json:"date"`
	FiveBest  float64 `json:"fiveBest"`
	Stock     float64 `json:"stock"`
	Bond      float64 `json:"bond"`
	Other     float64 `json:"other"`
	Cash      float64 `json:"cash"`
	Deposit   float64 `json:"deposit"`
	FundUnit  float64 `json:"fundUnit"`
	Commodity float64 `json:"commodity"`
}

type FundEfficiency struct {
	Date                string  `json:"date"`
	DailyEfficiency     float64 `json:"dailyEfficiency"`
	WeeklyEfficiency    float64 `json:"weeklyEfficiency"`
	MonthlyEfficiency   float64 `json:"monthlyEfficiency"`
	QuarterlyEfficiency float64 `json:"quarterlyEfficiency"`
	SixMonthEfficiency  float64 `json:"sixMonthEfficiency"`
	AnnualEfficiency    float64 `json:"annualEfficiency"`
}

type FundBasicInfo struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Item    Fund   `json:"item"`
}

type FundInfo struct {
	Fund           Fund             `json:"fund-basic-info"`
	FundEfficiency []FundEfficiency `json:"fund-efficiency"`
	FundPortfolio  []FundPortfolio  `json:"fund-portfolio"`
}

type IssueAndCancelData struct {
	Date        string `json:"date" time_format:"2006-01-02T15:04:05"`
	NetAsset    int64  `json:"netAsset"`
	UnitsSubDAY int    `json:"unitsSubDAY"`
	UnitsRedDAY int    `json:"unitsRedDAY"`
}

type IssueAndCancelSum struct {
	UnitsSubDAYSum int     `json:"unitsSubDAYSum"`
	UnitsRedDAYSum int     `json:"unitsRedDAYSum"`
	Profit         float64 `json:"profit"`
}
