package timings

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings struct {
			Fajr       string `json:"Fajr"`
			Sunrise    string `json:"Sunrise"`
			Dhuhr      string `json:"Dhuhr"`
			Asr        string `json:"Asr"`
			Sunset     string `json:"Sunset"`
			Maghrib    string `json:"Maghrib"`
			Isha       string `json:"Isha"`
			Imsak      string `json:"Imsak"`
			Midnight   string `json:"Midnight"`
			Firstthird string `json:"Firstthird"`
			Lastthird  string `json:"Lastthird"`
		} `json:"timings"`
		Date struct {
			Readable  string `json:"readable"`
			Timestamp string `json:"timestamp"`
			Hijri     struct {
				Date string `json:"date"`
				Day  string `json:"day"`
			} `json:"hijri"`
			Month struct {
				Number int    `json:"number"`
				En     string `json:"en"`
				Days   int    `json:"days"`
			} `json:"month"`
			Year string `json:"year"`
		} `json:"date"`
		Gregorian struct {
			Date    string `json:"date"`
			Day     string `json:"day"`
			Weekday struct {
				En string `json:"en"`
			} `json:"weekday"`
			Month struct {
				Number int    `json:"number"`
				En     string `json:"en"`
			} `json:"month"`
			Year string `json:"year"`
		} `json:"gregorian"`
	} `json:"data"`
	Meta struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
		Method    struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Params struct {
				Fajr int `json:"Fajr"`
				Isha int `json:"Isha"`
			} `json:"params"`
			Location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
		} `json:"method"`
	} `json:"meta"`
}
