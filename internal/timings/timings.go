package timings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/lanrey-waju/prayertimes/internal/cache"
	"github.com/spf13/viper"
)

func (p *PrayerTimes) String() string {
	var (
		purple = lipgloss.Color("99")
		white  = lipgloss.Color("#FFF")
		green  = lipgloss.Color("#04B575")
		red    = lipgloss.Color("160")

		headerStyle      = lipgloss.NewStyle().Foreground(white).Bold(true).Align(lipgloss.Center)
		cellStyle        = lipgloss.NewStyle().Padding(0, 1).Width(14)
		notOverTimeStyle = cellStyle.Foreground(green)
		overTimeStyle    = cellStyle.Foreground(red)
		infoStyle        = lipgloss.NewStyle().
					Width(70).
					Bold(true).
					Foreground(white).
					Align(lipgloss.Center)
	)

	prayerTimes := []string{
		p.Data.Timings.Fajr,
		p.Data.Timings.Dhuhr,
		p.Data.Timings.Asr,
		p.Data.Timings.Maghrib,
		p.Data.Timings.Isha,
	}

	currentTime := time.Now().Format("15:04")
	dateToday := time.Now().
		Format("02-01-2006") +
		fmt.Sprintf(
			" (Gregorian) | %s (Hijri)",
			p.Data.Date.Hijri.Date,
		)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}

			// check if time of prayer is over
			if col < len(prayerTimes) {
				if isPrayerTimeOver(prayerTimes[col], currentTime) {
					return overTimeStyle
				}
			}
			return notOverTimeStyle
		}).
		Headers("Fajr", "Dhuhr", "'Asr", "Maghrib", "'Ishaa").
		Rows(prayerTimes)
	fmt.Println(infoStyle.Render("Date:", dateToday, "Time:", currentTime))
	return t.Render()
}

// check if time of prayer is over
func isPrayerTimeOver(prayerTime, currentTime string) bool {
	var err error
	var prayerTimeParsed, currentTimeParsed time.Time

	if prayerTimeParsed, err = time.Parse("15:04", prayerTime); err != nil {
		fmt.Println("error parsing prayer time: ", err)
	}
	if currentTimeParsed, err = time.Parse("15:04", currentTime); err != nil {
		fmt.Println("error parsing current time: ", err)
	}
	return currentTimeParsed.After(prayerTimeParsed)
}

// RetrievePrayerTimes retrieves prayer times from the cache
func RetrievePrayerTimes(db *cache.Queries, city string) (*PrayerTimes, error) {
	prayerTimes, err := db.GetPrayerTimeForCity(
		context.Background(),
		cache.GetPrayerTimeForCityParams{
			City: city,
			Date: time.Now().Format("02-01-2006"),
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetPrayerTimes(db, city)
		}
		return &PrayerTimes{}, err
	}
	return databasePrayertimesToPrayerTimes(prayerTimes), nil
}

// GetPrayerTimes calls the aladhan API for the prayer times
func GetPrayerTimes(db *cache.Queries, city string) (*PrayerTimes, error) {
	c := http.Client{Timeout: 5 * time.Second}
	var prayertimes PrayerTimes

	baseURL := "https://api.aladhan.com/v1/timingsByAddress/"
	date := getDate()
	otherParams := "method=3&shafaq=general&tune=5%2C3%2C5%2C7%2C9%2C-1%2C0%2C8%2C-6&calendarMethod=UAQ"
	requestURL := fmt.Sprintf("%s/%s?address=%s&%s", baseURL, date, city, otherParams)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("could not create request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("could not make http request: %v", err)
		os.Exit(1)
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %v", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(dat, &prayertimes); err != nil {
		fmt.Printf("error unmarshalling response: %v", err)
		os.Exit(1)
	}

	if err := db.SavePrayerTimes(context.Background(), cache.SavePrayerTimesParams{
		City:      viper.GetString("location.city"),
		Date:      time.Now().Format("02-01-2006"),
		Fajr:      prayertimes.Data.Timings.Fajr,
		Dhuhr:     prayertimes.Data.Timings.Dhuhr,
		Asr:       prayertimes.Data.Timings.Asr,
		Maghrib:   prayertimes.Data.Timings.Maghrib,
		Isha:      prayertimes.Data.Timings.Isha,
		HijriDate: prayertimes.Data.Date.Hijri.Date,
		HijriDay:  prayertimes.Data.Date.Hijri.Day,
	}); err != nil {
		return &PrayerTimes{}, err
	}

	return &prayertimes, nil
}

// getDate returns today's date in dd-mm-yyyy format
func getDate() string {
	y, m, d := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", d, m, y)
	return date
}
