package http

import (
	"encoding/json"
	"log"
	"mini-avito/internal/avito_service"
	"mini-avito/internal/middleware"
	"mini-avito/internal/models"
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
		log.Println("[CreateAd] Error: missing user_id in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[CreateAd] Error parsing UUID (%s): %v\n", userIDStr, err)
		http.Error(w, "Unauthorized: invalid user id format", http.StatusUnauthorized)
		return
	}

	var req createAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[CreateAd] Error decoding JSON body: %v\n", err)
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	adID, err := h.usecase.CreateAd(r.Context(), userID, req.Status)
	if err != nil {
		log.Printf("[CreateAd] UseCase error (500): %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[CreateAd] Success! Ad %s created by user %s\n", adID, userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"ad_id": adID.String(),
	})
}

func (h *AdHandler) GetMyAds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	ads, err := h.usecase.GetAdsByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("[GetMyAds] Error fetching ads from DB: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if ads == nil {
		ads = make([]models.Ad, 0)
	}

	log.Printf("[GetMyAds] Success! Returned %d ads for user %s\n", len(ads), userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}
