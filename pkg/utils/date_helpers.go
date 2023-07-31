package utils

import (
	"fmt"
	persian "github.com/yaa110/go-persian-calendar"
	"strconv"
)

func ShamsiStringToGeoDate(date string) string {

	year, _ := strconv.Atoi(date[0:4])
	month, _ := strconv.Atoi(date[4:6])
	day, _ := strconv.Atoi(date[6:8])
	persianDate := persian.Date(year, persian.Month(month), day, 12, 1, 1, 0, persian.Iran())
	gregorianDate := persianDate.Time()

	// Format the Gregorian date as "YYYY-MM-DD"
	formattedDate := gregorianDate.Format("2006-01-02")

	return formattedDate
}
func CurrentShamsiDateNum() (*string, error) {
	time := persian.Now()
	year, month, day := time.Date()
	date := fmt.Sprintf("%d%02d%02d", year, month, day-1)
	return &date, nil
}
