package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"mini-avito/internal/models"
)

type mockAdRepo struct {
	CreateAdFunc       func(ctx context.Context, ad models.Ad) (uuid.UUID, error)
	GetAdsByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]models.Ad, error)
}

func (m *mockAdRepo) CreateAd(ctx context.Context, ad models.Ad) (uuid.UUID, error) {
	return m.CreateAdFunc(ctx, ad)
}

func (m *mockAdRepo) GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error) {
	return m.GetAdsByUserIDFunc(ctx, userID)
}

type mockPublisher struct {
	PublishAdCreatedFunc func(ctx context.Context, adID string) error
}

func (m *mockPublisher) PublishAdCreated(ctx context.Context, adID string) error {
	return m.PublishAdCreatedFunc(ctx, adID)
}

func TestCreateAd_Success(t *testing.T) {
	expectedID := uuid.New()
	userID := uuid.New()

	repo := &mockAdRepo{
		CreateAdFunc: func(ctx context.Context, ad models.Ad) (uuid.UUID, error) {
			if ad.UserID != userID {
				t.Errorf("expected UserID %v, got %v", userID, ad.UserID)
			}
			return expectedID, nil
		},
	}

	published := false
	pub := &mockPublisher{
		PublishAdCreatedFunc: func(ctx context.Context, adID string) error {
			if adID != expectedID.String() {
				t.Errorf("expected adID %s, got %s", expectedID.String(), adID)
			}
			published = true
			return nil
		},
	}

	uc := NewAdUseCase(repo, pub)

	id, err := uc.CreateAd(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != expectedID {
		t.Errorf("expected ID %v, got %v", expectedID, id)
	}

	if !published {
		t.Error("expected PublishAdCreated to be called")
	}
}
