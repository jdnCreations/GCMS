// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reservations.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createReservation = `-- name: CreateReservation :one
INSERT INTO reservations (
  id, 
  start_time, 
  end_time, 
  user_id, 
  game_id
  )
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4
)
RETURNING id, start_time, end_time, user_id, game_id
`

type CreateReservationParams struct {
	StartTime pgtype.Timestamp
	EndTime   pgtype.Timestamp
	UserID    pgtype.UUID
	GameID    pgtype.UUID
}

func (q *Queries) CreateReservation(ctx context.Context, arg CreateReservationParams) (Reservation, error) {
	row := q.db.QueryRow(ctx, createReservation,
		arg.StartTime,
		arg.EndTime,
		arg.UserID,
		arg.GameID,
	)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.EndTime,
		&i.UserID,
		&i.GameID,
	)
	return i, err
}

const getAllActiveReservations = `-- name: GetAllActiveReservations :many
SELECT id, start_time, end_time, user_id, game_id from reservations where end_time > NOW()
`

func (q *Queries) GetAllActiveReservations(ctx context.Context) ([]Reservation, error) {
	rows, err := q.db.Query(ctx, getAllActiveReservations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reservation
	for rows.Next() {
		var i Reservation
		if err := rows.Scan(
			&i.ID,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllReservations = `-- name: GetAllReservations :many
SELECT id, start_time, end_time, user_id, game_id from reservations
`

func (q *Queries) GetAllReservations(ctx context.Context) ([]Reservation, error) {
	rows, err := q.db.Query(ctx, getAllReservations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reservation
	for rows.Next() {
		var i Reservation
		if err := rows.Scan(
			&i.ID,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReservationsForUser = `-- name: GetReservationsForUser :many
SELECT id, start_time, end_time, user_id, game_id from reservations where user_id = $1
`

func (q *Queries) GetReservationsForUser(ctx context.Context, userID pgtype.UUID) ([]Reservation, error) {
	rows, err := q.db.Query(ctx, getReservationsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reservation
	for rows.Next() {
		var i Reservation
		if err := rows.Scan(
			&i.ID,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
