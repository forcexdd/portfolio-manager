package main

import (
	"fmt"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_api_client"
)

func main() {
	apiClient := moex_api_client.NewMoexApiClient()
	a, _ := apiClient.GetAllStocks("2024-07-25")
	for _, b := range a {
		fmt.Println(b)
	}
}
