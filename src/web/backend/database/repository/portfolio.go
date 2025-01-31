package repository

import (
	"database/sql"
	"errors"

	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	dtomodel "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

type PortfolioRepository interface {
	// Create creates new record of portfolio in DB. If there is another portfolio with the same name returns ErrPortfolioAlreadyExists
	Create(portfolio *model.Portfolio) error

	// GetByName returns record of portfolio from DB. If there is no portfolio with that name returns ErrPortfolioNotFound.
	// If there are no assets in portfolio returns ErrAssetNotFound
	GetByName(name string) (*model.Portfolio, error)

	// Update updates record of portfolio in DB. If there is no portfolio with that name returns ErrPortfolioNotFound
	Update(portfolio *model.Portfolio) error

	// Delete removes portfolio and all possible portfolio records from DB. If there is no portfolio with that name returns ErrPortfolioNotFound
	Delete(portfolio *model.Portfolio) error

	DeleteByName(name string) error

	// GetAll return all records of portfolio from DB. If there are no portfolios returns ErrPortfolioNotFound
	GetAll() ([]*model.Portfolio, error)
}

type postgresPortfolioRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPortfolioRepository(db *sql.DB, log logger.Logger) PortfolioRepository {
	return &postgresPortfolioRepository{
		db:  db,
		log: log,
	}
}

func (p *postgresPortfolioRepository) Create(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		p.log.Error("Failed to get portfolioID by name", "name", portfolio.Name, "error", err)
		return err
	}
	if portfolioID != 0 {
		p.log.Error("Portfolio already exists in DB", "name", portfolio.Name)
		return ErrPortfolioAlreadyExists
	}

	var createdPortfolio *dtomodel.Portfolio
	createdPortfolio, err = createPortfolio(p.db, portfolio.Name)
	if err != nil {
		p.log.Error("Failed to create portfolio", "name", portfolio.Name, "error", err)
		return err
	}
	if portfolio.AssetsQuantityMap == nil {
		p.log.Warn("Portfolio does not gave associated assets", "name", portfolio.Name)
		return nil
	} // If there are no assets just create name

	err = p.addManyAssetsToPortfolio(createdPortfolio.ID, portfolio.AssetsQuantityMap)
	if err != nil {
		p.log.Error("Failed to add assets to portfolio", "name", portfolio.Name, "error", err)
		return err
	}

	p.log.Info("Created portfolio", "name", portfolio.Name)

	return nil
}

func (p *postgresPortfolioRepository) GetByName(name string) (*model.Portfolio, error) {
	portfolioID, err := getPortfolioIDByName(p.db, name)
	if err != nil {
		p.log.Error("Failed to get portfolioID by name", "name", name, "error", err)
		return nil, err
	}
	if portfolioID == 0 {
		p.log.Warn("Portfolio does not exist in DB", "name", name)
		return nil, ErrPortfolioNotFound
	}

	assetsQuantityMap, err := p.getAssetsQuantityMapByPortfolioID(portfolioID)
	if err != nil {
		p.log.Error("Failed to get assets associated to portfolioID", "ID", portfolioID, "error", err)
		return nil, err
	}

	return &model.Portfolio{
		Name:              name,
		AssetsQuantityMap: assetsQuantityMap,
	}, nil
}

func (p *postgresPortfolioRepository) Update(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		p.log.Error("Failed to get portfolioID by name", "name", portfolio.Name, "error", err)
		return err
	}
	if portfolioID == 0 {
		p.log.Warn("Portfolio does not exist in DB", "name", portfolio.Name)
		return ErrPortfolioNotFound
	}

	var oldPortfolio *model.Portfolio
	oldPortfolio, err = p.GetByName(portfolio.Name)
	if err != nil {
		p.log.Error("Failed to get portfolio by name", "name", portfolio.Name, "error", err)
		return err
	}

	newAssetIDQuantityMap, err := p.convertAssetsQuantityMapToAssetsIDQuantityMap(portfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	oldAssetIDQuantityMap, err := p.convertAssetsQuantityMapToAssetsIDQuantityMap(oldPortfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	err = p.addOrUpdateNewPortfolioAssets(portfolioID, oldAssetIDQuantityMap, newAssetIDQuantityMap)
	if err != nil {
		p.log.Error("Failed to update portfolio assets", "name", portfolio.Name, "error", err)
		return err
	}

	err = p.deleteOldPortfolioAssets(portfolioID, oldAssetIDQuantityMap, newAssetIDQuantityMap)
	if err != nil {
		p.log.Error("Failed to delete portfolio assets", "name", portfolio.Name, "error", err)
		return err
	}

	p.log.Info("Updated portfolio", "name", portfolio.Name)

	return nil
}

func (p *postgresPortfolioRepository) Delete(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		p.log.Error("Failed to get portfolioID by name", "name", portfolio.Name, "error", err)
		return err
	}
	if portfolioID == 0 {
		p.log.Warn("Portfolio does not exist in DB", "name", portfolio.Name)
		return ErrPortfolioNotFound
	}

	err = p.deleteAllAssetsFromPortfolio(portfolioID)
	if err != nil {
		p.log.Error("Failed to delete all assets from portfolio", "name", portfolio.Name, "error", err)
		return err
	}

	err = deletePortfolio(p.db, portfolioID)
	if err != nil {
		p.log.Error("Failed to delete portfolio", "name", portfolio.Name, "error", err)
		return err
	}

	p.log.Info("Removed portfolio", "name", portfolio.Name)

	return nil
}

func (p *postgresPortfolioRepository) DeleteByName(name string) error {
	return p.Delete(&model.Portfolio{Name: name, AssetsQuantityMap: nil})
}

func (p *postgresPortfolioRepository) GetAll() ([]*model.Portfolio, error) {
	dtoPortfolios, err := getAllPortfolios(p.db)
	if err != nil {
		p.log.Error("Failed to get all portfolios", "error", err)
		return nil, err
	}
	if len(dtoPortfolios) == 0 {
		p.log.Warn("No portfolios found in DB")
		return nil, ErrPortfolioNotFound
	}

	var portfolios []*model.Portfolio
	var portfolio *model.Portfolio
	for _, dtoPortfolio := range dtoPortfolios {
		portfolio, err = p.GetByName(dtoPortfolio.Name)
		if err != nil {
			p.log.Error("Failed to get portfolio by name", "name", dtoPortfolio.Name)
			return nil, err
		}

		newPortfolio := &model.Portfolio{
			Name:              portfolio.Name,
			AssetsQuantityMap: portfolio.AssetsQuantityMap,
		}

		portfolios = append(portfolios, newPortfolio)
	}

	return portfolios, nil
}

func (p *postgresPortfolioRepository) getAssetsQuantityMapByPortfolioID(portfolioID int) (map[*model.Asset]int, error) {
	portfolioAssets, err := getAllPortfolioAssetsByPortfolioID(p.db, portfolioID)
	if err != nil {
		return nil, err
	}

	assetsQuantityMap := make(map[*model.Asset]int)
	var asset *dtomodel.Asset
	var relationship *dtomodel.PortfolioAssetRelationship
	for _, portfolioAsset := range portfolioAssets {
		asset, err = getAsset(p.db, portfolioAsset.AssetID)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			return nil, ErrAssetNotFound
		}

		relationship, err = getPortfolioAssetRelationshipByPortfolioAssetID(p.db, portfolioAsset.ID)
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

		assetsQuantityMap[&newAsset] = relationship.Quantity
	}

	return assetsQuantityMap, nil
}

func (p *postgresPortfolioRepository) addManyAssetsToPortfolio(portfolioID int, assetsQuantityMap map[*model.Asset]int) error {
	assetsIDQuantityMap, err := p.convertAssetsQuantityMapToAssetsIDQuantityMap(assetsQuantityMap)
	if err != nil {
		return err
	}

	for assetID, quantity := range assetsIDQuantityMap {
		err = p.addAssetToPortfolio(portfolioID, assetID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresPortfolioRepository) addOrUpdateNewPortfolioAssets(portfolioID int, oldAssetIDQuantityMap, newAssetIDQuantityMap map[int]int) error {
	for assetID, quantity := range newAssetIDQuantityMap {
		_, assetExists := oldAssetIDQuantityMap[assetID]

		if !assetExists {
			err := p.addAssetToPortfolio(portfolioID, assetID, quantity) // Since asset already exists in DB we are just adding it to portfolio
			if err != nil {
				return err
			}
		} else {
			err := p.updatePortfolioAsset(portfolioID, assetID, quantity) // Updating to new quantity
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresPortfolioRepository) addAssetToPortfolio(portfolioID, assetID, quantity int) error {
	portfolioAsset, err := createPortfolioAsset(p.db, portfolioID, assetID)
	if err != nil {
		return err
	}

	_, err = createPortfolioAssetRelationship(p.db, portfolioAsset.ID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresPortfolioRepository) updatePortfolioAsset(portfolioID, assetID, quantity int) error {
	portfolioAssetID, err := getPortfolioAssetIDByPortfolioIdAndAssetID(p.db, portfolioID, assetID)
	if err != nil {
		return err
	}

	var portfolioAssetRelationship *dtomodel.PortfolioAssetRelationship
	portfolioAssetRelationship, err = getPortfolioAssetRelationshipByPortfolioAssetID(p.db, portfolioAssetID)
	if err != nil {
		return err
	}

	if portfolioAssetRelationship.Quantity != quantity {
		err = updatePortfolioAssetRelationship(p.db, portfolioAssetRelationship.ID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresPortfolioRepository) deleteOldPortfolioAssets(portfolioID int, oldAssetIDQuantityMap, newAssetIDQuantityMap map[int]int) error {
	for assetID := range oldAssetIDQuantityMap {
		_, assetExists := newAssetIDQuantityMap[assetID]
		if !assetExists {
			err := p.deleteAssetFromTablesConnectedToPortfolio(portfolioID, assetID) // Just delete from connected tables. Shouldn't delete asset itself
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresPortfolioRepository) deleteAssetFromTablesConnectedToPortfolio(portfolioID int, assetID int) error {
	portfolioAssetID, err := getPortfolioAssetIDByPortfolioIdAndAssetID(p.db, portfolioID, assetID)
	if err != nil {
		return err
	}

	var portfolioAssetRelationship *dtomodel.PortfolioAssetRelationship
	portfolioAssetRelationship, err = getPortfolioAssetRelationshipByPortfolioAssetID(p.db, portfolioAssetID)
	if err != nil {
		return err
	}

	err = deletePortfolioAssetRelationship(p.db, portfolioAssetRelationship.ID)
	if err != nil {
		return err
	}

	err = deletePortfolioAsset(p.db, portfolioAssetID)

	return err
}

func (p *postgresPortfolioRepository) deleteAllAssetsFromPortfolio(portfolioID int) error {
	portfolioAssets, err := getAllPortfolioAssetsByPortfolioID(p.db, portfolioID)
	if err != nil {
		return err
	}

	for _, portfolioAsset := range portfolioAssets {
		var portfolioAssetRelationship *dtomodel.PortfolioAssetRelationship
		portfolioAssetRelationship, err = getPortfolioAssetRelationshipByPortfolioAssetID(p.db, portfolioAsset.ID)
		if err != nil {
			return err
		}

		err = deletePortfolioAssetRelationship(p.db, portfolioAssetRelationship.ID)
		if err != nil {
			return err
		}

		err = deletePortfolioAsset(p.db, portfolioAsset.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresPortfolioRepository) convertAssetsQuantityMapToAssetsIDQuantityMap(assetsQuantityMap map[*model.Asset]int) (map[int]int, error) {
	assetsNameQuantityMap := make(map[string]int)
	for asset, quantity := range assetsQuantityMap {
		assetsNameQuantityMap[asset.Name] = quantity
	}

	assetsIDQuantityMap, err := p.convertAssetsNameQuantityMapToAssetsIDQuantityMap(assetsNameQuantityMap)
	if err != nil {
		return nil, err
	}

	return assetsIDQuantityMap, nil
}

func (p *postgresPortfolioRepository) convertAssetsNameQuantityMapToAssetsIDQuantityMap(assetsQuantityMap map[string]int) (map[int]int, error) {
	assetsIDQuantityMap := make(map[int]int)
	for assetName, quantity := range assetsQuantityMap {
		assetID, err := getAssetIDByName(p.db, assetName)
		if err != nil {
			return nil, err
		}

		assetsIDQuantityMap[assetID] = quantity
	}

	return assetsIDQuantityMap, nil
}
