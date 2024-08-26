package tradingplatform

import (
	"errors"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/client"
	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"time"
)

type AssetExchangeService interface {
	ParseAllAssetsIntoDb() error
	ParseAllIndexesIntoDb() error
}

type MoexService struct {
	AssetRepository repository.AssetRepository
	IndexRepository repository.IndexRepository
	moexApiClient   *client.MoexApiClient
	time            time.Time
}

func NewTradingPlatformService(assetRepository repository.AssetRepository, indexRepository repository.IndexRepository) AssetExchangeService {
	newService := &MoexService{
		AssetRepository: assetRepository,
		IndexRepository: indexRepository,
		time:            time.Time{},
	}
	newService.setApiClient()

	return newService
}

func (m *MoexService) ParseAllAssetsIntoDb() error {
	var allAssets []*moexmodels.AssetData
	var err error
	allAssets, m.time, err = m.parseLatestAssets(getMaxDaysBeforeLatestDate())
	if err != nil {
		return err
	}

	var allAssetsInDb []*model.Asset
	allAssetsInDb, err = m.AssetRepository.GetAll()
	if err != nil {
		return err
	}

	for _, asset := range allAssets {
		newAsset := &model.Asset{
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

	var allIndexesInDb []*model.Index
	allIndexesInDb, err = m.IndexRepository.GetAll()
	if err != nil {
		return err
	}

	for _, index := range allIndexes {
		var indexAssets []*moexmodels.IndexAssetsData
		indexAssets, err = m.parseIndexAssets(index, m.time)
		if err != nil {
			return err
		}
		if len(indexAssets) == 0 { // No assets in index (something wrong with index API response)
			continue
		}

		newAssetsFractionMap := make(map[*model.Asset]float64)
		newAssetsFractionMap, err = m.createAssetsFractionMapFromIndexAssets(indexAssets)
		if err != nil {
			if errors.Is(err, repository.ErrAssetNotFound) { // There is no such asset in database (something wrong with index API response)
				continue
			}
			return err
		}

		newIndex := &model.Index{
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
	m.moexApiClient = client.NewMoexApiClient()
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
