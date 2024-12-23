// For /users/{id} endpoint - specific users only
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserByIDHandler struct {
	queries *sqlc.Queries
}

func NewUserByIDHandler(q *sqlc.Queries) *UserByIDHandler {
	return &UserByIDHandler{
		queries: q,
	}
}

func (h *UserByIDHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/")) // "/users/3"
	parts := strings.Split(cleanPath, "/")

	if len(parts) != 3 { // /users/{id} should have exactly 3 parts
		response.SendError(w, "Invalid URL format - must be '/users/{user_id}'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, parts)
	case http.MethodPatch:
		h.UpdateUser(w, r, parts)
	case http.MethodDelete:
		h.DeleteUser(w, r, parts)
	default:
		response.SendError(w, "Method not allowed - only GetUser, UpdateUser, and DeleteUser at path '/users/{user_id}'", http.StatusMethodNotAllowed)
	}
}

func (h *UserByIDHandler) GetUser(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUser(r.Context(), int32(id))
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user)
}

func (h *UserByIDHandler) UpdateUser(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
		Username string `json:"username,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := sqlc.UpdateUserParams{
		UserID:   int32(id),
		Email:    request.Email,
		Username: request.Username,
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			response.SendError(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		params.PasswordHash = string(hashedPassword)
	}

	user, err := h.queries.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		// handle dupes
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				response.SendError(w, "Email already in use", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				response.SendError(w, "Username already taken", http.StatusConflict)
			} else {
				response.SendError(w, "Duplicate value", http.StatusConflict)
			}
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user)
}

func (h *UserByIDHandler) DeleteUser(w http.ResponseWriter, r *http.Request, parts []string) {
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.DeleteUser(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, user, http.StatusOK) // Not StatusNoContent bc this is a soft delete
}
