package tradingplatform

import (
	"errors"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"time"
)

func (m *moexService) parseIndexAssets(index *moexmodels.IndexData, parseTime time.Time) ([]*moexmodels.IndexAssetsData, error) {
	allIndexAssets, err := m.moexApiClient.GetAllIndexAssets(formatTime(parseTime), index)
	if err != nil {
		return nil, err
	}

	return allIndexAssets, nil
}

func (m *moexService) createOrUpdateIndex(index *model.Index) error {
	_, err := m.IndexRepository.GetByName(index.Name)
	if err != nil {
		if errors.Is(err, repository.ErrIndexNotFound) { // Can't find index in database
			err = m.IndexRepository.Create(index)
			if err != nil {
				return err
			}
			m.log.Info("Created index", "name", index.Name)
		}
		return err
	} else { // Found already existing index in database
		err = m.IndexRepository.Update(index)
		if err != nil {
			return err
		}
		m.log.Info("Updated index", "name", index.Name)
	}

	return nil
}

func (m *moexService) createAssetsFractionMapFromIndexAssets(indexAssets []*moexmodels.IndexAssetsData) (map[*model.Asset]float64, error) {
	newAssetsFractionMap := make(map[*model.Asset]float64)
	for _, indexAsset := range indexAssets {
		asset, err := m.AssetRepository.GetByName(indexAsset.SecIDs)
		if err != nil {
			return nil, err
		}

		newAssetsFractionMap[asset] = indexAsset.Weight
	}

	return newAssetsFractionMap, nil
}

func removeIndexByNameFromSlice(indexes []*model.Index, name string) []*model.Index {
	for i, index := range indexes {
		if index.Name == name {
			return removeElementFromSliceByIndex[*model.Index](indexes, i)
		}
	}

	return indexes
}

func (m *moexService) removeOldIndexesFromDB(indexes []*model.Index) error {
	for _, index := range indexes {
		err := m.IndexRepository.Delete(index)
		if err != nil {
			return err
		}
		m.log.Info("Removed index", "name", index.Name)
	}

	return nil
}
