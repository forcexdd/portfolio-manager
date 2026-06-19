package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnStr  string
	Port       string
	MoexAPIURL string
}

func Load() *Config {
	_ = godotenv.Load()

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	moex := os.Getenv("MOEX_API_URL")
	if moex == "" {
		moex = "https://iss.moex.com/iss"
	}

	return &Config{
		DBConnStr:  db,
		Port:       port,
		MoexAPIURL: moex,
	}
}
