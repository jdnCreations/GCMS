package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// func ConvertToPGTime(str string) (pgtype.Time, error) {
// 	log.Printf("trying to convert %s  to a pgtime", str)
// 	parsedTime, err := time.Parse("15:04", str)
// 	if err != nil {
// 		return pgtype.Time{}, err
// 	}

// 	var pgTime pgtype.Time
// 	pgTime.Microseconds = parsedTime.UnixMicro()
// 	pgTime.Valid = true

// 	log.Printf("pgTime: %s", parsedTime)

// 	return pgTime, nil
// }

func ConvertToPGTime(str string) (pgtype.Time, error) {
	// Parse the string into time.Time
	parsedTime, err := time.Parse("15:04", str)
	if err != nil {
		return pgtype.Time{}, err
	}

	// Create pgtype.Time and assign the microseconds
	var pgTime pgtype.Time
	pgTime.Microseconds = int64(parsedTime.Hour() *3600*1e6 + parsedTime.Minute()*60*1e6)
	pgTime.Valid = true

	return pgTime, nil
}