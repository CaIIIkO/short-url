package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type ServiceInterface interface {
	Register(ctx context.Context, input *RegisterRequest) (*User, error)
	Authenticate(ctx context.Context, input *LoginRequest) (string, error)
}

type Handler struct {
	service ServiceInterface
}

func NewUserHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if input.Email == "" || input.Password == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	//Приведение емайла к нижнему регистру
	input.Email = strings.ToLower(input.Email)

	//Вызов сервиса
	token, err := h.service.Authenticate(r.Context(), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if input.Email == "" || input.Password == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	//Приведение емайла к нижнему регистру
	input.Email = strings.ToLower(input.Email)

	//Вызов сервиса
	user, err := h.service.Register(r.Context(), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
