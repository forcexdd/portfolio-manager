package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"log"
)

func main() {
	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	_, err := storage.NewStorage(connString)
	if err != nil {
		log.Fatal(err)
	}

	//err = storage.DeleteStorage()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
