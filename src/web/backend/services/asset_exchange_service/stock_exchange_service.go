package asset_exchange_service

import (
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/repositories"
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_api_client"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_models"
	"time"
)

type AssetExchangeService interface {
	ParseAllAssetsIntoDb() error
	ParseAllIndexesIntoDb() error
}

type MoexService struct {
	AssetRepository repositories.AssetRepository
	IndexRepository repositories.IndexRepository
	moexApiClient   *moex_api_client.MoexApiClient
	time            time.Time
}

func NewAssetExchangeService(assetRepository repositories.AssetRepository, indexRepository repositories.IndexRepository) AssetExchangeService {
	newService := &MoexService{
		AssetRepository: assetRepository,
		IndexRepository: indexRepository,
		time:            time.Time{},
	}
	newService.setApiClient()

	return newService
}

func (m *MoexService) ParseAllAssetsIntoDb() error {
	var allAssets []*moex_models.AssetData
	var err error
	allAssets, m.time, err = m.parseLatestAssets(getMaxDaysBeforeLatestDate())
	if err != nil {
		return err
	}

	var allAssetsInDb []*models.Asset
	allAssetsInDb, err = m.AssetRepository.GetAll()
	if err != nil {
		return err
	}

	for _, asset := range allAssets {
		newAsset := &models.Asset{
			Name:  asset.SecId,
			Price: asset.CurPrice,
		}

		err = m.createOrUpdateAsset(newAsset)
		if err != nil {
			return err
		}

		allAssetsInDb = removeAssetByNameFromSlice(allAssetsInDb, newAsset.Name)
	}

	if len(allAssetsInDb) > 0 {
		err = m.removeOldAssetsFromDb(allAssetsInDb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) ParseAllIndexesIntoDb() error {
	if m.time.IsZero() {
		return errors.New("assets should be parsed first")
	}

	allIndexes, err := m.moexApiClient.GetAllIndexes()
	if err != nil {
		return err
	}

	var allIndexesInDb []*models.Index
	allIndexesInDb, err = m.IndexRepository.GetAll()
	if err != nil {
		return err
	}

	for _, index := range allIndexes {
		var indexAssets []*moex_models.IndexAssetsData
		indexAssets, err = m.parseIndexAssets(index, m.time)
		if err != nil {
			return err
		}
		if len(indexAssets) == 0 { // No assets in index (something wrong with index API response)
			continue
		}

		newAssetsFractionMap := make(map[*models.Asset]float64)
		newAssetsFractionMap, err = m.createAssetsFractionMapFromIndexAssets(indexAssets)
		if err != nil {
			if errors.Is(err, repositories.ErrAssetNotFound) { // There is no such asset in database (something wrong with index API response)
				continue
			}
			return err
		}

		newIndex := &models.Index{
			Name:              index.IndexId,
			AssetsFractionMap: newAssetsFractionMap,
		}

		err = m.createOrUpdateIndex(newIndex)
		if err != nil {
			return err
		}

		allIndexesInDb = removeIndexByNameFromSlice(allIndexesInDb, newIndex.Name)
	}

	if len(allIndexesInDb) > 0 {
		err = m.removeOldIndexesFromDb(allIndexesInDb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) setApiClient() {
	m.moexApiClient = moex_api_client.NewMoexApiClient()
}

func getMaxDaysBeforeLatestDate() int {
	return 15
}

func removeElementFromSliceByIndex[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func getCurrentTime() time.Time {
	currTime := time.Now()
	return currTime
}

func formatTime(time time.Time) string {
	return time.Format("2006-01-02")
}
