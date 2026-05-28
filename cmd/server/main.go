package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	httpdelivery "mini-avito/internal/avito_service/delivery/http"
	"mini-avito/internal/avito_service/repository"
	"mini-avito/internal/avito_service/usecase"
	"mini-avito/internal/config"
)

func main() {

	cfg := config.New()
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

	mux := http.NewServeMux()

	mux.HandleFunc("/dbtest", testHandler.DBTest)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	srv := NewServer("8080", mux)

	if err := srv.Run(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
