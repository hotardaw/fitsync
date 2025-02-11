// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package sqlc

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password_hash, username)
VALUES ($1, $2, $3)
RETURNING user_id, email, password_hash, username, active, created_at, updated_at, last_login
`

type CreateUserParams struct {
	Email        string
	PasswordHash string
	Username     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.PasswordHash, arg.Username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.PasswordHash,
		&i.Username,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const deleteUser = `-- name: DeleteUser :one
UPDATE users
SET active = false, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING user_id, email, password_hash, username, active, created_at, updated_at, last_login
`

// Soft delete only - too many headaches if this gets actual users.
func (q *Queries) DeleteUser(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, deleteUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.PasswordHash,
		&i.Username,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT user_id, email, password_hash, username, active, created_at, updated_at, last_login
FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Email,
			&i.PasswordHash,
			&i.Username,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastLogin,
		); err != nil {
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

const getUser = `-- name: GetUser :one

SELECT users.user_id, users.email, users.password_hash, users.username, users.active, users.created_at, users.updated_at, users.last_login
FROM users
WHERE users.user_id = $1
`

// The comments above each query are SQLc directives dictating the naming of the generated Go func & what type of query it is (:one, :many, or :exec).
//
//	:many returns a slice of records via QueryContext
//
// :one returns a single record via QueryRowContext
// :exec returns the error from ExecContext
// More: https://docs.sqlc.dev/en/latest/reference/query-annotations.html
func (q *Queries) GetUser(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.PasswordHash,
		&i.Username,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, email, password_hash, username, active, created_at, updated_at, last_login FROM users 
WHERE email = $1 AND active = true
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.PasswordHash,
		&i.Username,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateLastLogin = `-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login = CURRENT_TIMESTAMP
WHERE user_id = $1
`

func (q *Queries) UpdateLastLogin(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, updateLastLogin, userID)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET 
  email = COALESCE($2, email),
  password_hash = COALESCE($3, password_hash),
  username = COALESCE($4, username),
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING user_id, email, password_hash, username, active, created_at, updated_at, last_login
`

type UpdateUserParams struct {
	UserID       int32
	Email        string
	PasswordHash string
	Username     string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.UserID,
		arg.Email,
		arg.PasswordHash,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.PasswordHash,
		&i.Username,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}
