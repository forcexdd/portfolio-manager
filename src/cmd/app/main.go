package main

import (
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/repositories"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/storage"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service"
	"log"
	"time"
)

func main() {
	start := time.Now()

	const connString = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		panic(err)
	}
	//db.DeleteStorage()

	assetRepository := repositories.NewAssetRepository(db.GetDb())
	indexRepository := repositories.NewIndexRepository(db.GetDb())

	assetExchangeService := asset_exchange_service.NewAssetExchangeService(assetRepository, indexRepository)

	err = assetExchangeService.ParseAllAssetsIntoDb()
	if err != nil {
		panic(err)
	}

	err = assetExchangeService.ParseAllIndexesIntoDb()
	if err != nil {
		panic(err)
	}

	err = db.CloseConnection()
	if err != nil {
		panic(err)
	}

	log.Println("Time passed: ", time.Since(start))
}
