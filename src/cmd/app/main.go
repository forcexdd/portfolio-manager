package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database"
	"log"
)

func main() {
	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	err := database.DeleteStorage(connString)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.CreateNewStorage(connString)
	if err != nil {
		log.Fatal(err)
	}
}
