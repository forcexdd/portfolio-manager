package repository

import (
	"database/sql"
	"errors"
	dtomodels "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

type IndexRepository interface {
	Create(index *model.Index) error
	GetByName(name string) (*model.Index, error)
	Update(index *model.Index) error
	Delete(index *model.Index) error
	DeleteByName(name string) error
	GetAll() ([]*model.Index, error)
}

type PostgresIndexRepository struct {
	db *sql.DB
}

func NewIndexRepository(db *sql.DB) IndexRepository {
	return &PostgresIndexRepository{db: db}
}

func (p *PostgresIndexRepository) Create(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexID != 0 {
		return ErrIndexAlreadyExists
	}

	var createdIndex *dtomodels.Index
	createdIndex, err = createIndex(p.db, index.Name)
	if err != nil {
		return err
	}
	if index.AssetsFractionMap == nil {
		return nil
	}

	assetsIDQuantityMap := make(map[int]float64)
	assetsIDQuantityMap, err = convertAssetsFractionMapToAssetsIDFractionMap(p.db, index.AssetsFractionMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToIndex(p.db, createdIndex.ID, assetsIDQuantityMap)

	return err
}

func (p *PostgresIndexRepository) GetByName(name string) (*model.Index, error) {
	indexID, err := getIndexIDByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if indexID == 0 {
		return nil, ErrIndexNotFound
	}

	var indexAssets []*dtomodels.IndexAsset
	indexAssets, err = getAllIndexAssetsByIndexID(p.db, indexID)
	if err != nil {
		return nil, err
	}

	assetsFractionMap := make(map[*model.Asset]float64)
	var asset *dtomodels.Asset
	var relationship *dtomodels.IndexAssetRelationship
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

	return &model.Index{
		Name:              name,
		AssetsFractionMap: assetsFractionMap,
	}, nil
}

func (p *PostgresIndexRepository) Update(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexID == 0 {
		return ErrIndexNotFound
	}

	err = deleteAllAssetsFromIndex(p.db, indexID)
	if err != nil {
		return err
	}

	assetsIDFractionMap := make(map[int]float64)
	assetsIDFractionMap, err = convertAssetsFractionMapToAssetsIDFractionMap(p.db, index.AssetsFractionMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToIndex(p.db, indexID, assetsIDFractionMap)

	return err
}

func (p *PostgresIndexRepository) Delete(index *model.Index) error {
	indexID, err := getIndexIDByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexID == 0 {
		return ErrIndexNotFound
	}

	err = deleteAllAssetsFromIndex(p.db, indexID)
	if err != nil {
		return err
	}

	err = deleteIndex(p.db, indexID)

	return err
}

func (p *PostgresIndexRepository) DeleteByName(name string) error {
	return p.Delete(&model.Index{Name: name, AssetsFractionMap: nil})
}

func (p *PostgresIndexRepository) GetAll() ([]*model.Index, error) {
	dtoIndexes, err := getAllIndexes(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoIndexes) == 0 {
		return nil, ErrIndexNotFound
	}

	var indexes []*model.Index
	var index *model.Index
	for _, dtoIndex := range dtoIndexes {
		index, err = p.GetByName(dtoIndex.Name)
		if err != nil {
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
