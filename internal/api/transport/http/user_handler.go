package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-arch-template/internal/api/usecase"
)

// HandleUsers handles user list and creation
// @Summary List or create users
// @Description Get list of all users or create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body usecase.CreateUserRequest false "User data"
// @Success 200 {array} usecase.UserResponse "List of users"
// @Success 201 {object} usecase.UserResponse "Created user"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [get]
// @Router /api/users [post]

type UserHandler struct {
	useCase *usecase.UserUseCase
}

func NewUserHandler(useCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		users, err := h.useCase.ListUsers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, users)

	case http.MethodPost:
		var req usecase.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := h.useCase.CreateUser(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusCreated, user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUser handles single user operations
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} usecase.UserResponse "User information"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [get]
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")

	switch r.Method {
	case http.MethodGet:
		user, err := h.useCase.GetUser(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		respondJSON(w, http.StatusOK, user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

