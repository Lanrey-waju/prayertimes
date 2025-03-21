-- name: GetPrayerTimeForCity :one
SELECT fajr, dhuhr, asr, maghrib, isha FROM prayer_times
WHERE city = ? AND date = ?;

-- name: SavePrayerTimes :exec
INSERT INTO prayer_times (city, date, fajr, dhuhr, asr, maghrib, isha)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(city, date) DO UPDATE SET
    fajr = excluded.fajr,
    dhuhr = excluded.dhuhr,
    asr = excluded.asr,
    maghrib = excluded.maghrib,
    isha = excluded.isha;
