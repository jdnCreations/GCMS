// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Game struct {
	ID     uuid.UUID
	Title  string
	Copies int16
}

type GameGenre struct {
	GameID  uuid.UUID
	GenreID uuid.UUID
}

type Genre struct {
	ID   uuid.UUID
	Name string
}

type Reservation struct {
	ID        uuid.UUID
	ResDate   pgtype.Date
	StartTime pgtype.Time
	EndTime   pgtype.Time
	UserID    pgtype.UUID
	GameID    pgtype.UUID
}

type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}
