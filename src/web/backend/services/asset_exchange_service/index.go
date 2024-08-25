package asset_exchange_service

import (
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/repositories"
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_models"
	"log"
	"time"
)

func (m *MoexService) parseIndexAssets(index *moex_models.IndexData, parseTime time.Time) ([]*moex_models.IndexAssetsData, error) {
	allIndexAssets, err := m.moexApiClient.GetAllIndexAssets(formatTime(parseTime), index)
	if err != nil {
		return nil, err
	}

	return allIndexAssets, nil
}

func (m *MoexService) createOrUpdateIndex(index *models.Index) error {
	_, err := m.IndexRepository.GetByName(index.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrIndexNotFound) { // Can't find index in database
			err = m.IndexRepository.Create(index)
			if err != nil {
				return err
			}
		}
		return err
	} else { // Found already existing index in database
		err = m.IndexRepository.Update(index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) createAssetsFractionMapFromIndexAssets(indexAssets []*moex_models.IndexAssetsData) (map[*models.Asset]float64, error) {
	newAssetsFractionMap := make(map[*models.Asset]float64)
	for _, indexAsset := range indexAssets {
		asset, err := m.AssetRepository.GetByName(indexAsset.SecIds)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			log.Printf("Asset %s does not exist", indexAsset.SecIds)
		}

		newAssetsFractionMap[asset] = indexAsset.Weight
	}

	return newAssetsFractionMap, nil
}

func removeIndexByNameFromSlice(indexes []*models.Index, name string) []*models.Index {
	for i, index := range indexes {
		if index.Name == name {
			return removeElementFromSliceByIndex[*models.Index](indexes, i)
		}
	}

	return indexes
}

func (m *MoexService) removeOldIndexesFromDb(indexes []*models.Index) error {
	for _, index := range indexes {
		err := m.IndexRepository.Delete(index)
		if err != nil {
			return err
		}
	}

	return nil
}
