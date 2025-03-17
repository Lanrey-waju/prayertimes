package timings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GetPrayerTimes calls the aladhan API for the prayer times
func GetPrayerTimes() Response {
	var response Response

	baseURL := "https://api.aladhan.com/v1/timingsByAddress/"
	date := getDate()
	otherParams := "method=3&shafaq=general&tune=5%2C3%2C5%2C7%2C9%2C-1%2C0%2C8%2C-6&calendarMethod=UAQ"
	cityAddr := GetCity()
	requestURL := fmt.Sprintf("%s/%s?address=%s&%s", baseURL, date, cityAddr, otherParams)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("could not create request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("could not make http request: %v", err)
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %v", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(dat, &response); err != nil {
		fmt.Printf("error unmarshalling response: %v", err)
	}

	fmt.Println(response.Code, " ", response.Data.Date.Hijri.Date)
	return response
}

// getDate returns today's date in dd-mm-yyyy format
func getDate() string {
	y, m, d := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", d, m, y)
	return date
}
