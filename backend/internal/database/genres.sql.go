// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: genres.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createGenre = `-- name: CreateGenre :one
INSERT INTO genres (id, name)
VALUES (
  gen_random_uuid(),
  $1
)
RETURNING id, name
`

func (q *Queries) CreateGenre(ctx context.Context, name string) (Genre, error) {
	row := q.db.QueryRow(ctx, createGenre, name)
	var i Genre
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteGenreById = `-- name: DeleteGenreById :exec
DELETE FROM genres WHERE id = $1
`

func (q *Queries) DeleteGenreById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteGenreById, id)
	return err
}

const getAllGenres = `-- name: GetAllGenres :many
SELECT id, name from genres
`

// maybe add functionality to pass in what to sort by?
func (q *Queries) GetAllGenres(ctx context.Context) ([]Genre, error) {
	rows, err := q.db.Query(ctx, getAllGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Genre
	for rows.Next() {
		var i Genre
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGenreById = `-- name: GetGenreById :one
SELECT id, name FROM genres where id = $1
`

func (q *Queries) GetGenreById(ctx context.Context, id uuid.UUID) (Genre, error) {
	row := q.db.QueryRow(ctx, getGenreById, id)
	var i Genre
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateGenre = `-- name: UpdateGenre :one
UPDATE genres 
SET name = COALESCE(NULLIF($1, ''), name)
WHERE id = $2
RETURNING id, name
`

type UpdateGenreParams struct {
	Column1 interface{}
	ID      uuid.UUID
}

func (q *Queries) UpdateGenre(ctx context.Context, arg UpdateGenreParams) (Genre, error) {
	row := q.db.QueryRow(ctx, updateGenre, arg.Column1, arg.ID)
	var i Genre
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateGenreCopies = `-- name: UpdateGenreCopies :one
UPDATE genres 
SET name = $1
WHERE id = $2
RETURNING id, name
`

type UpdateGenreCopiesParams struct {
	Name string
	ID   uuid.UUID
}

func (q *Queries) UpdateGenreCopies(ctx context.Context, arg UpdateGenreCopiesParams) (Genre, error) {
	row := q.db.QueryRow(ctx, updateGenreCopies, arg.Name, arg.ID)
	var i Genre
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateGenrename = `-- name: UpdateGenrename :one
UPDATE genres 
SET name = $1
WHERE id = $2
RETURNING id, name
`

type UpdateGenrenameParams struct {
	Name string
	ID   uuid.UUID
}

func (q *Queries) UpdateGenrename(ctx context.Context, arg UpdateGenrenameParams) (Genre, error) {
	row := q.db.QueryRow(ctx, updateGenrename, arg.Name, arg.ID)
	var i Genre
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
