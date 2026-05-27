package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	httpdelivery "mini-avito/internal/avito_service/delivery/http"
	"mini-avito/internal/avito_service/repository"
	"mini-avito/internal/avito_service/usecase"
	"mini-avito/internal/config"
)

func main() {
	cfg := config.New()
	//fmt.Printf("DEBUG CONFIG: Host='%s', Port='%d'\n", cfg.DBHost, cfg.DBPort)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userUC := usecase.NewUserUseCase(userRepo)
	authHandler := httpdelivery.NewAuthHandler(userUC)

	mux := http.NewServeMux()
	testRepo := repository.NewPostgresRepository(db)
	testHandler := httpdelivery.NewHandler(testRepo)
	mux.HandleFunc("/dbtest", testHandler.DBTest)
	mux.HandleFunc("/register", authHandler.Register)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
