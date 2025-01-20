// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reservations.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkGameReservation = `-- name: CheckGameReservation :one
SELECT
  COALESCE((SELECT copies FROM games WHERE game_id = $1), 0) - COALESCE(COUNT(*), 0) AS available_copies
FROM reservations
WHERE reservations.game_id = $1
AND reservations.start_time < $2
AND reservations.end_time > $3
GROUP BY reservations.game_id
`

type CheckGameReservationParams struct {
	GameID    pgtype.UUID
	StartTime pgtype.Time
	EndTime   pgtype.Time
}

func (q *Queries) CheckGameReservation(ctx context.Context, arg CheckGameReservationParams) (int32, error) {
	row := q.db.QueryRow(ctx, checkGameReservation, arg.GameID, arg.StartTime, arg.EndTime)
	var available_copies int32
	err := row.Scan(&available_copies)
	return available_copies, err
}

const createReservation = `-- name: CreateReservation :one
INSERT INTO reservations (
  id, 
  res_date,
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
  $4,
  $5
)
RETURNING id, res_date, start_time, end_time, user_id, game_id, active
`

type CreateReservationParams struct {
	ResDate   pgtype.Date
	StartTime pgtype.Time
	EndTime   pgtype.Time
	UserID    pgtype.UUID
	GameID    pgtype.UUID
}

func (q *Queries) CreateReservation(ctx context.Context, arg CreateReservationParams) (Reservation, error) {
	row := q.db.QueryRow(ctx, createReservation,
		arg.ResDate,
		arg.StartTime,
		arg.EndTime,
		arg.UserID,
		arg.GameID,
	)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.ResDate,
		&i.StartTime,
		&i.EndTime,
		&i.UserID,
		&i.GameID,
		&i.Active,
	)
	return i, err
}

const deleteReservation = `-- name: DeleteReservation :exec
DELETE FROM reservations where id = $1
`

func (q *Queries) DeleteReservation(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteReservation, id)
	return err
}

const getAllActiveReservations = `-- name: GetAllActiveReservations :many
SELECT id, res_date, start_time, end_time, user_id, game_id, active from reservations where end_time > NOW()
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
			&i.ResDate,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
			&i.Active,
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
SELECT id, res_date, start_time, end_time, user_id, game_id, active from reservations
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
			&i.ResDate,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
			&i.Active,
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

const getExpiredReservations = `-- name: GetExpiredReservations :many
SELECT id, res_date, start_time, end_time, user_id, game_id, active from reservations where end_time < NOW()::TIME AND active = true
`

func (q *Queries) GetExpiredReservations(ctx context.Context) ([]Reservation, error) {
	rows, err := q.db.Query(ctx, getExpiredReservations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reservation
	for rows.Next() {
		var i Reservation
		if err := rows.Scan(
			&i.ID,
			&i.ResDate,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
			&i.Active,
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

const getReservationById = `-- name: GetReservationById :one
SELECT id, res_date, start_time, end_time, user_id, game_id, active from reservations where id = $1
`

func (q *Queries) GetReservationById(ctx context.Context, id uuid.UUID) (Reservation, error) {
	row := q.db.QueryRow(ctx, getReservationById, id)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.ResDate,
		&i.StartTime,
		&i.EndTime,
		&i.UserID,
		&i.GameID,
		&i.Active,
	)
	return i, err
}

const getReservationsForDay = `-- name: GetReservationsForDay :many
SELECT reservations.id, reservations.res_date, reservations.start_time, reservations.end_time, reservations.user_id, reservations.game_id, reservations.active, games.title 
from reservations 
JOIN games on reservations.game_id = games.id
where res_date = $1
`

type GetReservationsForDayRow struct {
	ID        uuid.UUID
	ResDate   pgtype.Date
	StartTime pgtype.Time
	EndTime   pgtype.Time
	UserID    pgtype.UUID
	GameID    pgtype.UUID
	Active    bool
	Title     string
}

func (q *Queries) GetReservationsForDay(ctx context.Context, resDate pgtype.Date) ([]GetReservationsForDayRow, error) {
	rows, err := q.db.Query(ctx, getReservationsForDay, resDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReservationsForDayRow
	for rows.Next() {
		var i GetReservationsForDayRow
		if err := rows.Scan(
			&i.ID,
			&i.ResDate,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
			&i.Active,
			&i.Title,
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
SELECT reservations.id, reservations.res_date, reservations.start_time, reservations.end_time, reservations.user_id, reservations.game_id, reservations.active,
       games.title as game_name
FROM reservations
JOIN
  games
ON
  reservations.game_id = games.id
where user_id = $1
ORDER BY res_date asc
`

type GetReservationsForUserRow struct {
	ID        uuid.UUID
	ResDate   pgtype.Date
	StartTime pgtype.Time
	EndTime   pgtype.Time
	UserID    pgtype.UUID
	GameID    pgtype.UUID
	Active    bool
	GameName  string
}

func (q *Queries) GetReservationsForUser(ctx context.Context, userID pgtype.UUID) ([]GetReservationsForUserRow, error) {
	rows, err := q.db.Query(ctx, getReservationsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReservationsForUserRow
	for rows.Next() {
		var i GetReservationsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.ResDate,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
			&i.GameID,
			&i.Active,
			&i.GameName,
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

const setReservationInactive = `-- name: SetReservationInactive :one
UPDATE reservations
SET active = false
WHERE id = $1 AND
end_time < NOW()::TIME
RETURNING id, res_date, start_time, end_time, user_id, game_id, active
`

func (q *Queries) SetReservationInactive(ctx context.Context, id uuid.UUID) (Reservation, error) {
	row := q.db.QueryRow(ctx, setReservationInactive, id)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.ResDate,
		&i.StartTime,
		&i.EndTime,
		&i.UserID,
		&i.GameID,
		&i.Active,
	)
	return i, err
}
