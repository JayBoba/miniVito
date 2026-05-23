package http

import (
	"encoding/json"
	"net/http"

	"mini-avito/internal/avito_service/repository"
)

type Handler struct {
	repo *repository.PostgresRepository
}

func NewHandler(repo *repository.PostgresRepository) *Handler {
	return &Handler{repo: repo}
}

type TestRequest struct {
	Text string `json:"text"`
}

func (h *Handler) DBTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.repo.SaveTestMessage(r.Context(), req.Text); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message saved successfully"))
}