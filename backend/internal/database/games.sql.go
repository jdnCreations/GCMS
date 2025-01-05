// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: games.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (id, title, copies)
VALUES (
  gen_random_uuid(),
  $1,
  $2
)
RETURNING id, title, copies
`

type CreateGameParams struct {
	Title  string
	Copies int16
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, createGame, arg.Title, arg.Copies)
	var i Game
	err := row.Scan(&i.ID, &i.Title, &i.Copies)
	return i, err
}

const deleteGameById = `-- name: DeleteGameById :exec
DELETE FROM games WHERE id = $1
`

func (q *Queries) DeleteGameById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteGameById, id)
	return err
}

const getAllGames = `-- name: GetAllGames :many
SELECT id, title, copies from games
`

// maybe add functionality to pass in what to sort by?
func (q *Queries) GetAllGames(ctx context.Context) ([]Game, error) {
	rows, err := q.db.QueryContext(ctx, getAllGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Game
	for rows.Next() {
		var i Game
		if err := rows.Scan(&i.ID, &i.Title, &i.Copies); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGameById = `-- name: GetGameById :one
SELECT id, title, copies FROM games where id = $1
`

func (q *Queries) GetGameById(ctx context.Context, id uuid.UUID) (Game, error) {
	row := q.db.QueryRowContext(ctx, getGameById, id)
	var i Game
	err := row.Scan(&i.ID, &i.Title, &i.Copies)
	return i, err
}

const updateGame = `-- name: UpdateGame :one
UPDATE games 
SET title = COALESCE(NULLIF($1, ''), title),
    copies = CASE WHEN $2::SMALLINT IS NOT NULL THEN $2 ELSE copies END
WHERE id = $3
RETURNING id, title, copies
`

type UpdateGameParams struct {
	Column1 interface{}
	Column2 sql.NullInt16
	ID      uuid.UUID
}

func (q *Queries) UpdateGame(ctx context.Context, arg UpdateGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, updateGame, arg.Column1, arg.Column2, arg.ID)
	var i Game
	err := row.Scan(&i.ID, &i.Title, &i.Copies)
	return i, err
}

const updateGameCopies = `-- name: UpdateGameCopies :one
UPDATE games 
SET copies = $1
WHERE id = $2
RETURNING id, title, copies
`

type UpdateGameCopiesParams struct {
	Copies int16
	ID     uuid.UUID
}

func (q *Queries) UpdateGameCopies(ctx context.Context, arg UpdateGameCopiesParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, updateGameCopies, arg.Copies, arg.ID)
	var i Game
	err := row.Scan(&i.ID, &i.Title, &i.Copies)
	return i, err
}

const updateGameTitle = `-- name: UpdateGameTitle :one
UPDATE games 
SET title = $1
WHERE id = $2
RETURNING id, title, copies
`

type UpdateGameTitleParams struct {
	Title string
	ID    uuid.UUID
}

func (q *Queries) UpdateGameTitle(ctx context.Context, arg UpdateGameTitleParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, updateGameTitle, arg.Title, arg.ID)
	var i Game
	err := row.Scan(&i.ID, &i.Title, &i.Copies)
	return i, err
}
