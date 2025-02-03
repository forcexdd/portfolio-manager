package repository

import (
	"database/sql"
	"errors"

	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	dtomodel "github.com/forcexdd/portfoliomanager/src/database/model"
	"github.com/forcexdd/portfoliomanager/src/model"
)

type IndexRepository interface {
	// Create creates new record of index in DB. If there is another index with the same name returns ErrIndexAlreadyExists
	Create(index *model.Index) error

	// GetByName returns record of index from DB. If there is no index with that name returns ErrIndexNotFound.
	// If there are no assets in index returns ErrAssetNotFound
	GetByName(name string) (*model.Index, error)

	// Update updates record of index in DB. If there is no index with that name returns ErrIndexNotFound
	Update(index *model.Index) error

	// Delete removes index and all possible index records from DB. If there is no index with that name returns ErrIndexNotFound
	Delete(index *model.Index) error

	DeleteByName(name string) error

	// GetAll return all records of indexes from DB. If there are no indexes returns ErrIndexNotFound
	GetAll() ([]*model.Index, error)
}

type postgresIndexRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewIndexRepository(db *sql.DB, log logger.Logger) IndexRepository {
	return &postgresIndexRepository{
		db:  db,
		log: log,
	}
}

func (p *postgresIndexRepository) Create(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		p.log.Error("Failed to get indexID by name", "name", index.Name, "error", err)
		return err
	}
	if indexID != 0 {
		p.log.Warn("Index already exists in DB", "name", index.Name)
		return ErrIndexAlreadyExists
	}

	var createdIndex *dtomodel.Index
	createdIndex, err = createIndex(p.db, index.Name)
	if err != nil {
		p.log.Error("Failed to create index", "name", index.Name, "error", err)
		return err
	}
	if index.AssetsFractionMap == nil {
		p.log.Warn("Index does not have any associated assets", "name", index.Name)
		return nil
	} // If there are no assets then just create name

	err = p.addManyAssetsToIndex(createdIndex.ID, index.AssetsFractionMap)
	if err != nil {
		p.log.Error("Failed to add assets to index", "name", index.Name, "error", err)
		return err
	}

	p.log.Info("Created index", "name", index.Name)

	return nil
}

func (p *postgresIndexRepository) GetByName(name string) (*model.Index, error) {
	indexID, err := getIndexIDByName(p.db, name)
	if err != nil {
		p.log.Error("Failed to get indexID by name", "name", name, "error", err)
		return nil, err
	}
	if indexID == 0 {
		p.log.Warn("Index does not exist in DB", "name", name)
		return nil, ErrIndexNotFound
	}

	assetsFractionMap, err := p.getAssetsFractionMapByIndexID(indexID)
	if err != nil {
		p.log.Error("Failed to get assets fraction map by indexID", "ID", indexID, "error", err)
		return nil, err
	}

	return &model.Index{
		Name:              name,
		AssetsFractionMap: assetsFractionMap,
	}, nil
}

func (p *postgresIndexRepository) Update(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		p.log.Error("Failed to get indexID by name", "name", index.Name, "error", err)
		return err
	}
	if indexID == 0 {
		p.log.Warn("Index does not exist in DB", "name", index.Name)
		return ErrIndexNotFound
	}

	var oldIndex *model.Index
	oldIndex, err = p.GetByName(index.Name)
	if err != nil {
		p.log.Error("Failed to get index by name", "name", index.Name, "error", err)
		return err
	}

	newAssetIDFractionMap, err := p.convertAssetsFractionMapToAssetsIDFractionMap(index.AssetsFractionMap)
	if err != nil {
		return err
	}

	oldAssetIDFractionMap, err := p.convertAssetsFractionMapToAssetsIDFractionMap(oldIndex.AssetsFractionMap)
	if err != nil {
		return err
	}

	err = p.addOrUpdateNewIndexAssets(indexID, oldAssetIDFractionMap, newAssetIDFractionMap)
	if err != nil {
		p.log.Error("Failed to update index assets", "name", index.Name, "error", err)
		return err
	}

	err = p.deleteOldIndexAssets(indexID, oldAssetIDFractionMap, newAssetIDFractionMap)
	if err != nil {
		p.log.Error("Failed to delete old index assets", "name", index.Name, "error", err)
		return err
	}

	p.log.Info("Updated index", "name", index.Name)

	return nil
}

func (p *postgresIndexRepository) Delete(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		p.log.Error("Failed to get indexID by name", "name", index.Name, "error", err)
		return err
	}
	if indexID == 0 {
		p.log.Warn("Index does not exist in DB", "name", index.Name)
		return ErrIndexNotFound
	}

	err = p.deleteAllAssetsFromIndex(indexID)
	if err != nil {
		p.log.Error("Failed to delete all assets from index", "name", index.Name, "error", err)
		return err
	}

	err = deleteIndex(p.db, indexID)
	if err != nil {
		p.log.Error("Failed to delete index", "name", index.Name, "error", err)
		return err
	}

	p.log.Info("Removed index", "name", index.Name)

	return nil
}

func (p *postgresIndexRepository) DeleteByName(name string) error {
	return p.Delete(&model.Index{Name: name, AssetsFractionMap: nil})
}

func (p *postgresIndexRepository) GetAll() ([]*model.Index, error) {
	dtoIndexes, err := getAllIndexes(p.db)
	if err != nil {
		p.log.Error("Failed to get all indexes", "error", err)
		return nil, err
	}
	if len(dtoIndexes) == 0 {
		p.log.Warn("No indexes found in DB")
		return nil, ErrIndexNotFound
	}

	var indexes []*model.Index
	var index *model.Index
	for _, dtoIndex := range dtoIndexes {
		index, err = p.GetByName(dtoIndex.Name)
		if err != nil {
			p.log.Error("Failed to get index from DB", "name", dtoIndex.Name, "error", err)
			return nil, err
		}

		newIndex := &model.Index{
			Name:              index.Name,
			AssetsFractionMap: index.AssetsFractionMap,
		}

		indexes = append(indexes, newIndex)
	}

	return indexes, nil
}

func (p *postgresIndexRepository) getAssetsFractionMapByIndexID(indexID int) (map[*model.Asset]float64, error) {
	indexAssets, err := getAllIndexAssetsByIndexID(p.db, indexID)
	if err != nil {
		return nil, err
	}

	assetsFractionMap := make(map[*model.Asset]float64)
	var asset *dtomodel.Asset
	var relationship *dtomodel.IndexAssetRelationship
	for _, indexAsset := range indexAssets {
		asset, err = getAsset(p.db, indexAsset.AssetID)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			return nil, ErrAssetNotFound
		}

		relationship, err = getIndexAssetRelationshipByIndexAssetID(p.db, indexAsset.ID)
		if err != nil {
			return nil, err
		}
		if relationship == nil {
			return nil, errors.New("relationship not found")
		}

		newAsset := model.Asset{
			Name:  asset.Name,
			Price: asset.Price,
		}

		assetsFractionMap[&newAsset] = relationship.Fraction
	}

	return assetsFractionMap, err
}

func (p *postgresIndexRepository) addManyAssetsToIndex(indexID int, assetsFractionMap map[*model.Asset]float64) error {
	assetsIDFractionMap, err := p.convertAssetsFractionMapToAssetsIDFractionMap(assetsFractionMap)
	if err != nil {
		return err
	}

	for assetID, fraction := range assetsIDFractionMap {
		err = p.addAssetToIndex(indexID, assetID, fraction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresIndexRepository) addOrUpdateNewIndexAssets(indexID int, oldAssetIDFractionMap, newAssetIDFractionMap map[int]float64) error {
	for assetID, fraction := range newAssetIDFractionMap {
		_, assetExists := oldAssetIDFractionMap[assetID]

		if !assetExists {
			err := p.addAssetToIndex(indexID, assetID, fraction) // Since asset already exists in DB we are just adding it to index
			if err != nil {
				return err
			}
		} else {
			err := p.updateIndexAsset(indexID, assetID, fraction) // Updating to new fraction
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresIndexRepository) addAssetToIndex(indexID int, assetID int, fraction float64) error {
	indexAsset, err := createIndexAsset(p.db, indexID, assetID)
	if err != nil {
		return err
	}

	_, err = createIndexAssetRelationship(p.db, indexAsset.ID, fraction)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresIndexRepository) updateIndexAsset(indexID, assetID int, fraction float64) error {
	indexAssetID, err := getIndexAssetIDByIndexIdAndAssetID(p.db, indexID, assetID)
	if err != nil {
		return err
	}

	var indexAssetRelationship *dtomodel.IndexAssetRelationship
	indexAssetRelationship, err = getIndexAssetRelationshipByIndexAssetID(p.db, indexAssetID)
	if err != nil {
		return err
	}

	if indexAssetRelationship.Fraction != fraction {
		err = updateIndexAssetRelationship(p.db, indexAssetRelationship.ID, fraction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresIndexRepository) deleteOldIndexAssets(indexID int, oldAssetIDFractionMap, newAssetIDFractionMap map[int]float64) error {
	for assetID := range oldAssetIDFractionMap {
		_, assetExists := newAssetIDFractionMap[assetID]
		if !assetExists {
			err := p.deleteAssetFromTablesConnectedToIndex(indexID, assetID) // Just delete from connected tables. Shouldn't delete asset itself
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresIndexRepository) deleteAssetFromTablesConnectedToIndex(indexID int, assetID int) error {
	indexAssetID, err := getIndexAssetIDByIndexIdAndAssetID(p.db, indexID, assetID)
	if err != nil {
		return err
	}

	var indexAssetRelationship *dtomodel.IndexAssetRelationship
	indexAssetRelationship, err = getIndexAssetRelationshipByIndexAssetID(p.db, indexAssetID)
	if err != nil {
		return err
	}

	err = deleteIndexAssetRelationship(p.db, indexAssetRelationship.ID)
	if err != nil {
		return err
	}

	err = deleteIndexAsset(p.db, indexAssetID)

	return err
}

func (p *postgresIndexRepository) deleteAllAssetsFromIndex(indexID int) error {
	indexAssets, err := getAllIndexAssetsByIndexID(p.db, indexID)
	if err != nil {
		return err
	}

	for _, indexAsset := range indexAssets {
		var indexAssetRelationship *dtomodel.IndexAssetRelationship
		indexAssetRelationship, err = getIndexAssetRelationshipByIndexAssetID(p.db, indexAsset.ID)
		if err != nil {
			return err
		}

		err = deleteIndexAssetRelationship(p.db, indexAssetRelationship.ID)
		if err != nil {
			return err
		}

		err = deleteIndexAsset(p.db, indexAsset.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresIndexRepository) convertAssetsFractionMapToAssetsIDFractionMap(assetsFractionMap map[*model.Asset]float64) (map[int]float64, error) {
	assetsNameFractionMap := make(map[string]float64)
	for asset, quantity := range assetsFractionMap {
		assetsNameFractionMap[asset.Name] = quantity
	}

	assetsIDFractionMap, err := p.convertAssetsNameFractionMapToAssetsIDFractionMap(assetsNameFractionMap)
	if err != nil {
		return nil, err
	}

	return assetsIDFractionMap, nil
}

func (p *postgresIndexRepository) convertAssetsNameFractionMapToAssetsIDFractionMap(assetsFractionMap map[string]float64) (map[int]float64, error) {
	assetsIDFractionMap := make(map[int]float64)
	for assetName, fraction := range assetsFractionMap {
		assetID, err := getAssetIDByName(p.db, assetName)
		if err != nil {
			return nil, err
		}

		assetsIDFractionMap[assetID] = fraction
	}

	return assetsIDFractionMap, nil
}
