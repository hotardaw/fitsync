// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workouts.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createWorkout = `-- name: CreateWorkout :one
INSERT INTO workouts (user_id, workout_date, title)
VALUES ($1, $2, $3)
RETURNING workout_id, user_id, workout_date, title, created_at
`

type CreateWorkoutParams struct {
	UserID      sql.NullInt32
	WorkoutDate time.Time
	Title       sql.NullString
}

type CreateWorkoutRow struct {
	WorkoutID   int32
	UserID      sql.NullInt32
	WorkoutDate time.Time
	Title       sql.NullString
	CreatedAt   sql.NullTime
}

// CREATE: Insert a new workout
func (q *Queries) CreateWorkout(ctx context.Context, arg CreateWorkoutParams) (CreateWorkoutRow, error) {
	row := q.db.QueryRowContext(ctx, createWorkout, arg.UserID, arg.WorkoutDate, arg.Title)
	var i CreateWorkoutRow
	err := row.Scan(
		&i.WorkoutID,
		&i.UserID,
		&i.WorkoutDate,
		&i.Title,
		&i.CreatedAt,
	)
	return i, err
}

const deleteWorkout = `-- name: DeleteWorkout :one
DELETE FROM workouts
WHERE workout_id = $1 
AND user_id = $2
RETURNING workout_id
`

type DeleteWorkoutParams struct {
	WorkoutID int32
	UserID    sql.NullInt32
}

// DELETE: Remove a workout
func (q *Queries) DeleteWorkout(ctx context.Context, arg DeleteWorkoutParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, deleteWorkout, arg.WorkoutID, arg.UserID)
	var workout_id int32
	err := row.Scan(&workout_id)
	return workout_id, err
}

const getAllWorkoutsForUser = `-- name: GetAllWorkoutsForUser :many
SELECT workout_id, workout_date, title, created_at
FROM workouts
WHERE user_id = $1
ORDER BY workout_date DESC
`

type GetAllWorkoutsForUserRow struct {
	WorkoutID   int32
	WorkoutDate time.Time
	Title       sql.NullString
	CreatedAt   sql.NullTime
}

// READ: Get all workouts for a specific user
func (q *Queries) GetAllWorkoutsForUser(ctx context.Context, userID sql.NullInt32) ([]GetAllWorkoutsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllWorkoutsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllWorkoutsForUserRow
	for rows.Next() {
		var i GetAllWorkoutsForUserRow
		if err := rows.Scan(
			&i.WorkoutID,
			&i.WorkoutDate,
			&i.Title,
			&i.CreatedAt,
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

const getWorkoutByID = `-- name: GetWorkoutByID :one
SELECT workout_id, user_id, workout_date, title, created_at
FROM workouts
WHERE workout_id = $1
`

type GetWorkoutByIDRow struct {
	WorkoutID   int32
	UserID      sql.NullInt32
	WorkoutDate time.Time
	Title       sql.NullString
	CreatedAt   sql.NullTime
}

// READ: Get a specific workout by ID
func (q *Queries) GetWorkoutByID(ctx context.Context, workoutID int32) (GetWorkoutByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getWorkoutByID, workoutID)
	var i GetWorkoutByIDRow
	err := row.Scan(
		&i.WorkoutID,
		&i.UserID,
		&i.WorkoutDate,
		&i.Title,
		&i.CreatedAt,
	)
	return i, err
}

const getWorkoutsWithinDateRange = `-- name: GetWorkoutsWithinDateRange :many
SELECT workout_id, workout_date, title, created_at
FROM workouts
WHERE user_id = $1 
AND workout_date BETWEEN $2 AND $3
ORDER BY workout_date DESC
`

type GetWorkoutsWithinDateRangeParams struct {
	UserID        sql.NullInt32
	WorkoutDate   time.Time
	WorkoutDate_2 time.Time
}

type GetWorkoutsWithinDateRangeRow struct {
	WorkoutID   int32
	WorkoutDate time.Time
	Title       sql.NullString
	CreatedAt   sql.NullTime
}

// READ: Get workouts within a date range for a user
func (q *Queries) GetWorkoutsWithinDateRange(ctx context.Context, arg GetWorkoutsWithinDateRangeParams) ([]GetWorkoutsWithinDateRangeRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkoutsWithinDateRange, arg.UserID, arg.WorkoutDate, arg.WorkoutDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkoutsWithinDateRangeRow
	for rows.Next() {
		var i GetWorkoutsWithinDateRangeRow
		if err := rows.Scan(
			&i.WorkoutID,
			&i.WorkoutDate,
			&i.Title,
			&i.CreatedAt,
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

const updateWorkout = `-- name: UpdateWorkout :one
UPDATE workouts
SET workout_date = $1,
    title = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE workout_id = $3 
AND user_id = $4
RETURNING workout_id, workout_date, title, updated_at
`

type UpdateWorkoutParams struct {
	WorkoutDate time.Time
	Title       sql.NullString
	WorkoutID   int32
	UserID      sql.NullInt32
}

type UpdateWorkoutRow struct {
	WorkoutID   int32
	WorkoutDate time.Time
	Title       sql.NullString
	UpdatedAt   sql.NullTime
}

// UPDATE: Modify an existing workout
func (q *Queries) UpdateWorkout(ctx context.Context, arg UpdateWorkoutParams) (UpdateWorkoutRow, error) {
	row := q.db.QueryRowContext(ctx, updateWorkout,
		arg.WorkoutDate,
		arg.Title,
		arg.WorkoutID,
		arg.UserID,
	)
	var i UpdateWorkoutRow
	err := row.Scan(
		&i.WorkoutID,
		&i.WorkoutDate,
		&i.Title,
		&i.UpdatedAt,
	)
	return i, err
}
