package main

import (
	"net/http"
)
import "github.com/forcexdd/StockPortfolioManager/src/web/backend/handlers"

func main() {
	http.HandleFunc("/static/", handlers.HandleStaticFiles)
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/following_index", handlers.HandleHome)
	http.HandleFunc("/manager", handlers.HandleManager)
	http.HandleFunc("/add_portfolio", handlers.HandleAddPortfolio)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
