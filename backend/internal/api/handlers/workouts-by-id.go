// GET (ID only), PATCH, DELETE
// In future: add a "CopyWorkoutByDate" so users can easily repeat workouts
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
)

type WorkoutByIDHandler struct {
	queries *sqlc.Queries
}

func NewWorkoutByIDHandler(q *sqlc.Queries) *WorkoutByIDHandler {
	return &WorkoutByIDHandler{
		queries: q,
	}
}

type UpdateWorkoutByIDRequest struct {
	Title string `json:"workouttitle"`
}

func (h *WorkoutByIDHandler) HandleWorkoutsByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetWorkoutByIDForUser(w, r)
	case http.MethodPatch:
		h.UpdateWorkoutByID(w, r)
	case http.MethodDelete:
		h.DeleteWorkoutByID(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *WorkoutByIDHandler) GetWorkoutByIDForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	getWorkoutByIDForUserParams := sqlc.GetWorkoutByIDForUserParams{
		WorkoutID: workoutID,
		UserID:    utils.ToNullInt32(userID),
	}

	workout, err := h.queries.GetWorkoutByIDForUser(r.Context(), getWorkoutByIDForUserParams)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout)
}

func (h *WorkoutByIDHandler) UpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	var request UpdateWorkoutByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateWorkoutParams := sqlc.UpdateWorkoutParams{
		Title:     utils.ToNullString(request.Title),
		WorkoutID: workoutID,
		UserID:    utils.ToNullInt32(userID),
	}

	workout, err := h.queries.UpdateWorkout(r.Context(), updateWorkoutParams)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, workout)
}

func (h *WorkoutByIDHandler) DeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := utils.GetIDFromPath(r.URL.Path)
	if err != nil {
		response.SendError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	deleteWorkoutParams := sqlc.DeleteWorkoutParams{
		UserID:    utils.ToNullInt32(userID),
		WorkoutID: workoutID,
	}

	deletedWorkoutID, err := h.queries.DeleteWorkout(r.Context(), deleteWorkoutParams)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SendError(w, "Workout not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message": "Workout deleted successfully",
		"id":      int(deletedWorkoutID),
	})
}
