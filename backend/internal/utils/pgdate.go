package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPGDate(str string) (pgtype.Date, error) {
	parsed, err := time.Parse("2006-01-02", str)
	if err != nil {
		return pgtype.Date{}, err
	}

	var date pgtype.Date
	date.Scan(parsed)

	return date, nil
}