package repository

import (
	"database/sql"
	"errors"
	dtomodels "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

type AssetRepository interface {
	Create(asset *model.Asset) error
	GetByName(name string) (*model.Asset, error)
	Update(asset *model.Asset) error
	Delete(asset *model.Asset) error
	DeleteByName(name string) error
	GetAll() ([]*model.Asset, error)
}

type PostgresAssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) AssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (p *PostgresAssetRepository) Create(asset *model.Asset) error {
	assetId, err := getAssetIdByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetId != 0 {
		return ErrAssetAlreadyExists
	}

	_, err = createAsset(p.db, asset.Name, asset.Price)

	return err
}

func (p *PostgresAssetRepository) GetByName(name string) (*model.Asset, error) {
	assetId, err := getAssetIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if assetId == 0 {
		return nil, ErrAssetNotFound
	}

	var dtoAsset *dtomodels.Asset
	dtoAsset, err = getAsset(p.db, assetId)
	if err != nil {
		return nil, err
	}

	return &model.Asset{
		Name:  dtoAsset.Name,
		Price: dtoAsset.Price,
	}, nil
}

func (p *PostgresAssetRepository) Update(asset *model.Asset) error {
	assetId, err := getAssetIdByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetId == 0 {
		return ErrAssetNotFound
	}

	var dtoAsset *dtomodels.Asset
	dtoAsset, err = getAsset(p.db, assetId)
	if err != nil {
		return err
	}
	if dtoAsset.Name != asset.Name {
		return errors.New("asset name does not match")
	}

	err = updateAsset(p.db, assetId, asset.Price)

	return err
}

func (p *PostgresAssetRepository) Delete(asset *model.Asset) error {
	assetId, err := getAssetIdByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetId == 0 {
		return ErrAssetNotFound
	}

	err = deleteAssetFromConnectedTables(p.db, assetId)
	if err != nil {
		return err
	}

	err = deleteAsset(p.db, assetId)

	return err
}

func (p *PostgresAssetRepository) DeleteByName(name string) error {
	return p.Delete(&model.Asset{Name: name}) // Price doesn't matter
}

func (p *PostgresAssetRepository) GetAll() ([]*model.Asset, error) {
	dtoAssets, err := getAllAssets(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoAssets) == 0 {
		return nil, ErrAssetNotFound
	}

	var assets []*model.Asset
	var asset model.Asset
	for _, dtoAsset := range dtoAssets {
		asset.Name = dtoAsset.Name
		asset.Price = dtoAsset.Price

		newAsset := &model.Asset{
			Name:  asset.Name,
			Price: asset.Price,
		}

		assets = append(assets, newAsset)
	}

	return assets, nil
}
