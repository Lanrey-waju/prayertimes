package timings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/lanrey-waju/prayer-times/internal/cache"
	"github.com/lanrey-waju/prayer-times/internal/config"
	"github.com/spf13/viper"
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

func RetrievePrayerTimes(db *cache.Queries, city string) (*PrayerTimes, error) {
	defer config.TimeTrack(time.Now(), "RetrievePrayerTimes")
	prayerTimes, err := db.GetPrayerTimeForCity(
		context.Background(),
		cache.GetPrayerTimeForCityParams{
			City: city,
			Date: time.Now().Format("02-01-2006"),
		},
	)
	if err != nil {
		return &PrayerTimes{}, err
	}
	return databasePrayertimesToPrayerTimes(prayerTimes), nil
}

// GetPrayerTimes calls the aladhan API for the prayer times
func GetPrayerTimes(db *cache.Queries, city string) (*PrayerTimes, error) {
	var prayertimes *PrayerTimes

	if cache.DBExists() {
		return RetrievePrayerTimes(db, city)
	}

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

	if err := db.SavePrayerTimes(context.Background(), cache.SavePrayerTimesParams{
		City:    viper.GetString("location.city"),
		Date:    time.Now().Format("02-01-2006"),
		Fajr:    prayertimes.Data.Timings.Fajr,
		Dhuhr:   prayertimes.Data.Timings.Dhuhr,
		Asr:     prayertimes.Data.Timings.Asr,
		Maghrib: prayertimes.Data.Timings.Maghrib,
		Isha:    prayertimes.Data.Timings.Isha,
	}); err != nil {
		return &PrayerTimes{}, err
	}

	fmt.Println("Prayer times saved to cache")
	return prayertimes, nil
}

// getDate returns today's date in dd-mm-yyyy format
func getDate() string {
	y, m, d := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", d, m, y)
	return date
}
