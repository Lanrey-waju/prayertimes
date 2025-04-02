package timings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// getLoacation makes a GET request to ip-api.com and returns a Locadion
// data
func getLocation() *Location {
	var location Location
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		fmt.Printf("could not make request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not read response body: %v", err)
	}

	if err := json.Unmarshal(respBody, &location); err != nil {
		fmt.Printf("could not unmarshal response body into location struct: %v", err)
	}
	return &location
}

// GetCity uses getLocation function to return the City of the client
func GetLocationParams() (string, float64, float64) {
	loc := getLocation()
	return loc.City, loc.Latitude, loc.Longitude
}
