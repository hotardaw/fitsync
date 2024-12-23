// For /users endpoint - no specific users
package handlers

import (
	"encoding/json"
	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Holds dependencies for user handlers
type UserHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(q *sqlc.Queries) *UserHandler {
	return &UserHandler{
		queries: q,
	}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Ensure only /users endpoint is handled
	if len(parts) != 2 {
		response.SendError(w, "Invalid URL - must be '/users'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r, parts)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response.SendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Username:     request.Username,
	})

	if err != nil {
		// handle dupe emails/usernames
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

	response.SendSuccess(w, user, http.StatusCreated)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request, parts []string) {
	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		response.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, users)
}

/*
1. Define handler funcs that process incoming HTTP reqs
2. Map URLs/routes to their corresponding handlers
3. Handle req validation, auth, and authorization
4. Format & send HTTP responses

Import the SQLc queries
Define the HTTP hanndler functions that use these queries
Handle the HTTP request/response lifecycle

Then register these handlers to make them accessible
*/
