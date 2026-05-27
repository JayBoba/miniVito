package http

import (
	"encoding/json"
	"mini-avito/internal/avito_service"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"
)

type AuthHandler struct {
	usecase avito_service.UserUseCase
}

func NewAuthHandler(uc avito_service.UserUseCase) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
	}
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type registerResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

var loginRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	if !loginRegex.MatchString(req.Login) {
		http.Error(w, "Bad request: login must be 3-30 characters long and contain only english letters, numbers, hyphens or underscores", http.StatusBadRequest)
		return
	}

	passLen := utf8.RuneCountInString(req.Password)
	if passLen < 8 {
		http.Error(w, "Bad request: password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	if len(req.Password) > 72 {
		http.Error(w, "Bad request: password is too long (max 72 bytes)", http.StatusBadRequest)
		return
	}

	userID, err := h.usecase.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	json.NewEncoder(w).Encode(registerResponse{
		UserID: userID,
	})

}
