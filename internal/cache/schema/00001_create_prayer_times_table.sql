-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS prayer_times (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    city TEXT NOT NULL,
    date TEXT NOT NULL,
    fajr TEXT NOT NULL,
    dhuhr TEXT NOT NULL,
    asr TEXT NOT NULL,
    maghrib TEXT NOT NULL,
    isha TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(city, date)  -- Prevent duplicates per city per day
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE prayer_times;
-- +goose StatementEnd
