package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPGUUID(str string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID

	uuid, err := uuid.Parse(str)
	if err != nil {
		return pgUUID, err
	}

	var byteArray [16]byte
	copy(byteArray[:], uuid[:])

	pgUUID.Bytes = byteArray 
	pgUUID.Valid = true

	return pgUUID, nil
}
