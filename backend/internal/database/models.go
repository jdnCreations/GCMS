// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
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
	StartTime time.Time
	EndTime   time.Time
	UserID    uuid.NullUUID
	GameID    uuid.NullUUID
}

type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}
