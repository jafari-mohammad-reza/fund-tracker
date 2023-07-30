package structs

type Fund struct {
	RegNo                     string   `json:"regNo"`
	Name                      string   `json:"name"`
	RankOf12Month             *int     `json:"rankOf12Month"`
	RankOf24Month             *int     `json:"rankOf24Month"`
	RankOf36Month             *int     `json:"rankOf36Month"`
	RankOf48Month             *int     `json:"rankOf48Month"`
	RankOf60Month             *int     `json:"rankOf60Month"`
	RankLastUpdate            string   `json:"rankLastUpdate"`
	FundType                  int      `json:"fundType"`
	TypeOfInvest              string   `json:"typeOfInvest"`
	FundSize                  int64    `json:"fundSize"`
	InitiationDate            string   `json:"initiationDate"`
	DailyEfficiency           float64  `json:"dailyEfficiency"`
	WeeklyEfficiency          float64  `json:"weeklyEfficiency"`
	MonthlyEfficiency         float64  `json:"monthlyEfficiency"`
	QuarterlyEfficiency       float64  `json:"quarterlyEfficiency"`
	SixMonthEfficiency        float64  `json:"sixMonthEfficiency"`
	AnnualEfficiency          float64  `json:"annualEfficiency"`
	StatisticalNav            float64  `json:"statisticalNav"`
	Efficiency                float64  `json:"efficiency"`
	CancelNav                 float64  `json:"cancelNav"`
	IssueNav                  float64  `json:"issueNav"`
	DividendIntervalPeriod    int      `json:"dividendIntervalPeriod"`
	GuaranteedEarningRate     *float64 `json:"guaranteedEarningRate"`
	Date                      string   `json:"date"`
	NetAsset                  int64    `json:"netAsset"`
	EstimatedEarningRate      *float64 `json:"estimatedEarningRate"`
	InvestedUnits             int      `json:"investedUnits"`
	ArticlesOfAssociationLink *string  `json:"articlesOfAssociationLink"`
	ProsoectusLink            *string  `json:"prosoectusLink"`
	WebsiteAddress            []string `json:"websiteAddress"`
	Manager                   string   `json:"manager"`
	ManagerSeoRegisterNo      string   `json:"managerSeoRegisterNo"`
	GuarantorSeoRegisterNo    *string  `json:"guarantorSeoRegisterNo"`
	Auditor                   string   `json:"auditor"`
	Custodian                 string   `json:"custodian"`
	Guarantor                 string   `json:"guarantor"`
	Beta                      *float64 `json:"beta"`
	Alpha                     *float64 `json:"alpha"`
	IsCompleted               bool     `json:"isCompleted"`
	FiveBest                  float64  `json:"fiveBest"`
	Stock                     float64  `json:"stock"`
	Bond                      float64  `json:"bond"`
	Other                     float64  `json:"other"`
	Cash                      float64  `json:"cash"`
	Deposit                   float64  `json:"deposit"`
	FundUnit                  *float64 `json:"fundUnit"`
	Commodity                 *float64 `json:"commodity"`
	FundPublisher             int      `json:"fundPublisher"`
	FundWatch                 *bool    `json:"fundWatch"`
}

type FipIranResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
	TotalCount int    `json:"totalCount"`
	Items      []Fund `json:"items"`
}

type CalculatedFund struct {
	Fund
	Rank                int     `json:"rank"`
	RankDiff            int     `json:"rankDiff"`
	NetAssetDiff        float64 `json:"netAssetDiff"`
	NetAssetDiffPercent float64 `json:"netAssetDiffPercent"`
	IssueAndCancelSum
}

type EachYearFunds struct {
	Year  int    `json:"year"`
	Funds []Fund `json:"funds"`
}

type EachYearNav struct {
	Year  int   `json:"year"`
	Nav   int64 `json:"nav`
	Count int   `json:"count"`
}
