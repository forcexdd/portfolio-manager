package main

import (
	"net/http"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/handlers"
	"log"
)

func main() {
	http.HandleFunc("/static/", handlers.HandleStaticFiles)
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/following_index", handlers.HandleHome)
	http.HandleFunc("/manager", handlers.HandleManager)
	http.HandleFunc("/add_portfolio", handlers.HandleAddPortfolio)
	
	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	_, err := database.CreateNewStorage(connString)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	//err = storage.DeleteStorage()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
