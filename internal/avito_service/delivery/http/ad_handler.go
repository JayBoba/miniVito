package http

import (
	"encoding/json"
	"mini-avito/internal/avito_service"
	"mini-avito/internal/middleware"
	"net/http"

	"github.com/google/uuid"
)

type AdHandler struct {
	usecase avito_service.AdUseCase
}

func NewAdHandler(usecase avito_service.AdUseCase) *AdHandler {
	return &AdHandler{usecase: usecase}
}

type createAdRequest struct {
	Status string `json:"status"`
}

func (h *AdHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: missing user id in context", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Unauthorized: invalid user id format", http.StatusUnauthorized)
		return
	}

	var req createAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	adID, err := h.usecase.CreateAd(r.Context(), userID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"ad_id": adID.String(),
	})
}
