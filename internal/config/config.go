package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found")
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getEnvAsInt("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: getEnvAsInt("SERVER_PORT"),
	}
}

func getEnvAsInt(key string) int {
	val := os.Getenv(key)
	if val == "" {
		return 0
	}
	i, _ := strconv.Atoi(val)
	return i
}
