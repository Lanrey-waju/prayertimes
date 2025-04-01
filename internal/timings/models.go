package timings

import "github.com/lanrey-waju/prayer-times/internal/cache"

// databasePrayertimesToPrayerTimes converts a database prayer times row to a PrayerTimes struct
func databasePrayertimesToPrayerTimes(dbPrayerTimes cache.GetPrayerTimeForCityRow) *PrayerTimes {
	var prayerTimes PrayerTimes

	prayerTimes.Data.Timings.Fajr = dbPrayerTimes.Fajr
	prayerTimes.Data.Timings.Dhuhr = dbPrayerTimes.Dhuhr
	prayerTimes.Data.Timings.Asr = dbPrayerTimes.Asr
	prayerTimes.Data.Timings.Maghrib = dbPrayerTimes.Maghrib
	prayerTimes.Data.Timings.Isha = dbPrayerTimes.Isha
	prayerTimes.Data.Date.Hijri.Date = dbPrayerTimes.HijriDate
	prayerTimes.Data.Date.Hijri.Day = dbPrayerTimes.HijriDay

	return &prayerTimes
}
