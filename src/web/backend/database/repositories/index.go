package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/dto_models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
)

type IndexRepository interface {
	Create(index *models.Index) error
	GetByName(name string) (*models.Index, error)
	Update(index *models.Index) error
	Delete(index *models.Index) error
	DeleteByName(name string) error
	GetAll() ([]*models.Index, error)
}

type PostgresIndexRepository struct {
	db *sql.DB
}

func NewIndexRepository(db *sql.DB) IndexRepository {
	return &PostgresIndexRepository{db: db}
}

func (p *PostgresIndexRepository) Create(index *models.Index) error {
	indexId, err := getIndexIdByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexId != 0 {
		return errors.New("index already exists")
	}

	var createdIndex *dto_models.Index
	createdIndex, err = createIndex(p.db, index.Name)
	if err != nil {
		return err
	}
	if index.AssetsFractionMap == nil {
		return nil
	}

	assetsIdQuantityMap := make(map[int]float64)
	assetsIdQuantityMap, err = convertAssetsFractionMapToAssetsIdFractionMap(p.db, index.AssetsFractionMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToIndex(p.db, createdIndex.Id, assetsIdQuantityMap)

	return err
}

func (p *PostgresIndexRepository) GetByName(name string) (*models.Index, error) {
	indexId, err := getIndexIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if indexId == 0 {
		return nil, nil
	}

	var indexAssets []*dto_models.IndexAsset
	indexAssets, err = getAllIndexAssetsByIndexId(p.db, indexId)
	if err != nil {
		return nil, err
	}

	assetsFractionMap := make(map[*models.Asset]float64)
	var asset *dto_models.Asset
	var relationship *dto_models.IndexAssetRelationship
	for _, indexAsset := range indexAssets {
		asset, err = getAsset(p.db, indexAsset.AssetId)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			return nil, errors.New("asset not found")
		}

		relationship, err = getIndexAssetRelationshipByIndexAssetId(p.db, indexAsset.Id)
		if err != nil {
			return nil, err
		}
		if relationship == nil {
			return nil, errors.New("relationship not found")
		}

		newAsset := models.Asset{
			Name:  asset.Name,
			Price: asset.Price,
		}

		assetsFractionMap[&newAsset] = relationship.Fraction
	}

	return &models.Index{
		Name:              name,
		AssetsFractionMap: assetsFractionMap,
	}, nil
}

func (p *PostgresIndexRepository) Update(index *models.Index) error {
	indexId, err := getIndexIdByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexId == 0 {
		return errors.New("index not found")
	}

	err = deleteAllAssetsFromIndex(p.db, indexId)
	if err != nil {
		return err
	}

	assetsIdFractionMap := make(map[int]float64)
	assetsIdFractionMap, err = convertAssetsFractionMapToAssetsIdFractionMap(p.db, index.AssetsFractionMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToIndex(p.db, indexId, assetsIdFractionMap)

	return err
}

func (p *PostgresIndexRepository) Delete(index *models.Index) error {
	indexId, err := getIndexIdByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexId == 0 {
		return errors.New("index not found")
	}

	err = deleteAllAssetsFromIndex(p.db, indexId)
	if err != nil {
		return err
	}

	err = deleteIndex(p.db, indexId)

	return err
}

func (p *PostgresIndexRepository) DeleteByName(name string) error {
	return p.Delete(&models.Index{Name: name, AssetsFractionMap: nil})
}

func (p *PostgresIndexRepository) GetAll() ([]*models.Index, error) {
	dtoIndexes, err := getAllIndexes(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoIndexes) == 0 {
		return nil, nil
	}

	var indexes []*models.Index
	var index *models.Index
	for _, dtoIndex := range dtoIndexes {
		index, err = p.GetByName(dtoIndex.Name)
		if err != nil {
			return nil, err
		}

		newIndex := &models.Index{
			Name:              index.Name,
			AssetsFractionMap: index.AssetsFractionMap,
		}

		indexes = append(indexes, newIndex)
	}

	return indexes, nil
}
