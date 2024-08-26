package repository

import (
	"database/sql"
	dtomodel "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

type AssetRepository interface {
	// Create creates new record of asset in DB. If there is another asset with the same name returns ErrAssetAlreadyExists
	Create(asset *model.Asset) error

	// GetByName returns record of asset from DB. If there is no asset with that name returns ErrAssetNotFound
	GetByName(name string) (*model.Asset, error)

	// Update updates record of asset in DB. If there is no asset with that name returns ErrAssetNotFound
	Update(asset *model.Asset) error

	// Delete removes all possible asset records from DB. If there is no asset with that name returns ErrAssetNotFound
	Delete(asset *model.Asset) error

	DeleteByName(name string) error

	// GetAll return all records of asset from DB. If there is no assets returns ErrAssetNotFound
	GetAll() ([]*model.Asset, error)
}

type PostgresAssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) AssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (p *PostgresAssetRepository) Create(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetID != 0 {
		return ErrAssetAlreadyExists
	}

	_, err = createAsset(p.db, asset.Name, asset.Price)

	return err
}

func (p *PostgresAssetRepository) GetByName(name string) (*model.Asset, error) {
	assetID, err := getAssetIDByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if assetID == 0 {
		return nil, ErrAssetNotFound
	}

	var dtoAsset *dtomodel.Asset
	dtoAsset, err = getAsset(p.db, assetID)
	if err != nil {
		return nil, err
	}

	return &model.Asset{
		Name:  dtoAsset.Name,
		Price: dtoAsset.Price,
	}, nil
}

func (p *PostgresAssetRepository) Update(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetID == 0 {
		return ErrAssetNotFound
	}

	err = updateAsset(p.db, assetID, asset.Price)

	return err
}

func (p *PostgresAssetRepository) Delete(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetID == 0 {
		return ErrAssetNotFound
	}

	err = deleteAssetFromConnectedTables(p.db, assetID)
	if err != nil {
		return err
	}

	err = deleteAsset(p.db, assetID)

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
