package utils

import (
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPGTimestamp(str string) (pgtype.Timestamp, error) {
	location := os.Getenv("TIMEZONE")
	var timestamp pgtype.Timestamp
	if location == "" {
		location = "Australia/Perth"
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return timestamp, err
	}

	perthTime, err := time.ParseInLocation("2006-01-02T15:04", str, loc)
	if err != nil {
		return timestamp, err
	}

	utcTime := perthTime.UTC()


	timestamp.Time = utcTime
	timestamp.Valid = true

	return timestamp, nil
}