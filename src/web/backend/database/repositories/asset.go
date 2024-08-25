package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/dto_models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
)

type AssetRepository interface {
	Create(asset *models.Asset) error
	GetByName(name string) (*models.Asset, error)
	Update(asset *models.Asset) error
	Delete(asset *models.Asset) error
	DeleteByName(name string) error
	GetAll() ([]*models.Asset, error)
}

type PostgresAssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) AssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (p *PostgresAssetRepository) Create(asset *models.Asset) error {
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

func (p *PostgresAssetRepository) GetByName(name string) (*models.Asset, error) {
	assetId, err := getAssetIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if assetId == 0 {
		return nil, ErrAssetNotFound
	}

	var dtoAsset *dto_models.Asset
	dtoAsset, err = getAsset(p.db, assetId)
	if err != nil {
		return nil, err
	}

	return &models.Asset{
		Name:  dtoAsset.Name,
		Price: dtoAsset.Price,
	}, nil
}

func (p *PostgresAssetRepository) Update(asset *models.Asset) error {
	assetId, err := getAssetIdByName(p.db, asset.Name)
	if err != nil {
		return err
	}
	if assetId == 0 {
		return ErrAssetNotFound
	}

	var dtoAsset *dto_models.Asset
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

func (p *PostgresAssetRepository) Delete(asset *models.Asset) error {
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
	return p.Delete(&models.Asset{Name: name}) // Price doesn't matter
}

func (p *PostgresAssetRepository) GetAll() ([]*models.Asset, error) {
	dtoAssets, err := getAllAssets(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoAssets) == 0 {
		return nil, ErrAssetNotFound
	}

	var assets []*models.Asset
	var asset models.Asset
	for _, dtoAsset := range dtoAssets {
		asset.Name = dtoAsset.Name
		asset.Price = dtoAsset.Price

		newAsset := &models.Asset{
			Name:  asset.Name,
			Price: asset.Price,
		}

		assets = append(assets, newAsset)
	}

	return assets, nil
}
