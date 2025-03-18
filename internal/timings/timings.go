package timings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (p *PrayerTimes) String() string {
	return fmt.Sprintf(
		"Fajr: %s\nDhuhr: %s\n'Asr: %s\nMaghrib: %s\n'Ishaa: %s",
		p.Data.Timings.Fajr,
		p.Data.Timings.Dhuhr,
		p.Data.Timings.Asr,
		p.Data.Timings.Maghrib,
		p.Data.Timings.Isha,
	)
}

// GetPrayerTimes calls the aladhan API for the prayer times
func GetPrayerTimes(city string) *PrayerTimes {
	var prayertimes *PrayerTimes

	baseURL := "https://api.aladhan.com/v1/timingsByAddress/"
	date := getDate()
	otherParams := "method=3&shafaq=general&tune=5%2C3%2C5%2C7%2C9%2C-1%2C0%2C8%2C-6&calendarMethod=UAQ"
	requestURL := fmt.Sprintf("%s/%s?address=%s&%s", baseURL, date, city, otherParams)
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

	if err := json.Unmarshal(dat, &prayertimes); err != nil {
		fmt.Printf("error unmarshalling response: %v", err)
	}

	fmt.Println(prayertimes)
	return prayertimes
}

// getDate returns today's date in dd-mm-yyyy format
func getDate() string {
	y, m, d := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", d, m, y)
	return date
}
