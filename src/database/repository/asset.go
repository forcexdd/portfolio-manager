package repository

import (
	"database/sql"

	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	dtomodel "github.com/forcexdd/portfoliomanager/src/database/model"
	"github.com/forcexdd/portfoliomanager/src/model"
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

	// GetAll return all records of asset from DB. If there are no assets returns ErrAssetNotFound
	GetAll() ([]*model.Asset, error)
}

type postgresAssetRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewAssetRepository(db *sql.DB, log logger.Logger) AssetRepository {
	return &postgresAssetRepository{
		db:  db,
		log: log,
	}
}

func (p *postgresAssetRepository) Create(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		p.log.Error("Failed to get assetID by name", "name", asset.Name, "error", err)
		return err
	}
	if assetID != 0 {
		p.log.Warn("Asset already exists in DB", "name", asset.Name, "error", ErrAssetAlreadyExists)
		return ErrAssetAlreadyExists
	}

	_, err = createAsset(p.db, asset.Name, asset.Price)
	if err != nil {
		p.log.Error("Failed to create asset", "name", asset.Name, "error", err)
		return err
	}

	p.log.Info("Created asset", "name", asset.Name)

	return nil
}

func (p *postgresAssetRepository) GetByName(name string) (*model.Asset, error) {
	assetID, err := getAssetIDByName(p.db, name)
	if err != nil {
		p.log.Error("Failed to get assetID by name", "name", name, "error", err)
		return nil, err
	}
	if assetID == 0 {
		p.log.Warn("Asset does not exists in DB", "name", name, "error", ErrAssetNotFound)
		return nil, ErrAssetNotFound
	}

	var dtoAsset *dtomodel.Asset
	dtoAsset, err = getAsset(p.db, assetID)
	if err != nil {
		p.log.Error("Failed to get asset by assetID", "ID", assetID, "error", err)
		return nil, err
	}

	return &model.Asset{
		Name:  dtoAsset.Name,
		Price: dtoAsset.Price,
	}, nil
}

func (p *postgresAssetRepository) Update(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		p.log.Error("Failed to get assetID by name", "name", asset.Name, "error", err)
		return err
	}
	if assetID == 0 {
		p.log.Warn("Asset does not exists in DB", "name", asset.Name, "error", ErrAssetNotFound)
		return ErrAssetNotFound
	}

	var oldAsset *model.Asset
	oldAsset, err = p.GetByName(asset.Name)
	if err != nil {
		p.log.Error("Failed to get asset by name", "name", asset.Name, "error", err)
		return err
	}

	if oldAsset.Price != asset.Price {
		err = updateAsset(p.db, assetID, asset.Price)
		if err != nil {
			p.log.Error("Failed to update asset", "name", asset.Name, "error", err)
			return err
		}
		p.log.Info("Updated asset", "name", asset.Name)
	} else {
		p.log.Info("Asset is up to date", "name", asset.Name)
	}

	return nil
}

func (p *postgresAssetRepository) Delete(asset *model.Asset) error {
	assetID, err := getAssetIDByName(p.db, asset.Name)
	if err != nil {
		p.log.Error("Failed to get assetID by name", "name", asset.Name, "error", err)
		return err
	}
	if assetID == 0 {
		p.log.Warn("Asset does not exists in DB", "name", asset.Name, "error", ErrAssetNotFound)
		return ErrAssetNotFound
	}

	err = deleteAssetFromConnectedTables(p.db, assetID)
	if err != nil {
		p.log.Error("Failed to delete asset from connected tables", "name", asset.Name, "error", err)
		return err
	}

	err = deleteAsset(p.db, assetID)
	if err != nil {
		p.log.Error("Failed to delete asset", "name", asset.Name, "error", err)
		return err
	}

	p.log.Info("Removed asset", "name", asset.Name)

	return nil
}

func (p *postgresAssetRepository) DeleteByName(name string) error {
	return p.Delete(&model.Asset{Name: name}) // Price doesn't matter
}

func (p *postgresAssetRepository) GetAll() ([]*model.Asset, error) {
	dtoAssets, err := getAllAssets(p.db)
	if err != nil {
		p.log.Error("Failed to get all assets", "error", err)
		return nil, err
	}
	if len(dtoAssets) == 0 {
		p.log.Warn("No assets found in DB", "error", ErrAssetNotFound)
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
