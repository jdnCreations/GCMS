package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPGTimestamp(str string) (pgtype.Timestamp, error) {
	var timestamp pgtype.Timestamp

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return timestamp, err
	}

	timestamp.Time = t
	timestamp.Valid = true

	return timestamp, nil
}