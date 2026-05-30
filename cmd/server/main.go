package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	httpdelivery "mini-avito/internal/avito_service/delivery/http"
	"mini-avito/internal/avito_service/repository"
	"mini-avito/internal/avito_service/usecase"
	"mini-avito/internal/config"
	"mini-avito/internal/jwt"
	"mini-avito/internal/middleware"
)

func main() {

	cfg := config.New()
	jwt.Init(jwt.Config{
		SecretKey:     cfg.JWTSecret,
		TokenDuration: 24 * time.Hour, // токен живет сутки!!!
	})
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	testRepo := repository.NewPostgresRepository(db)
	testHandler := httpdelivery.NewHandler(testRepo)

	userRepo := repository.NewUserRepo(db)
	userUC := usecase.NewUserUseCase(userRepo)
	authHandler := httpdelivery.NewAuthHandler(userUC)

	adRepo := repository.NewAdRepo(db)
	adUC := usecase.NewAdUseCase(adRepo)
	adHandler := httpdelivery.NewAdHandler(adUC)

	mux := http.NewServeMux()

	mux.HandleFunc("/dbtest", testHandler.DBTest)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(string)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access granted! Your User ID is: " + userID))
	})))

	mux.Handle("/ads/create", middleware.AuthMiddleware(http.HandlerFunc(adHandler.CreateAd)))
	mux.Handle("/ads/my", middleware.AuthMiddleware(http.HandlerFunc(adHandler.GetMyAds)))
	srv := NewServer("8080", mux)

	if err := srv.Run(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
