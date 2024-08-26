package repository

import (
	"database/sql"
	"errors"
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

type PostgresPortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) PortfolioRepository {
	return &PostgresPortfolioRepository{db: db}
}

func (p *PostgresPortfolioRepository) Create(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioID != 0 {
		return ErrPortfolioAlreadyExists
	}

	var createdPortfolio *dtomodel.Portfolio
	createdPortfolio, err = createPortfolio(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolio.AssetsQuantityMap == nil {
		return nil
	} // If there are no assets just create name

	err = p.addManyAssetsToPortfolio(createdPortfolio.ID, portfolio.AssetsQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) GetByName(name string) (*model.Portfolio, error) {
	portfolioID, err := getPortfolioIDByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if portfolioID == 0 {
		return nil, ErrPortfolioNotFound
	}

	var portfolioAssets []*dtomodel.PortfolioAsset
	portfolioAssets, err = getAllPortfolioAssetsByPortfolioID(p.db, portfolioID)
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

	return &model.Portfolio{
		Name:              name,
		AssetsQuantityMap: assetsQuantityMap,
	}, nil
}

func (p *PostgresPortfolioRepository) Update(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioID == 0 {
		return ErrPortfolioNotFound
	}

	err = p.deleteAllAssetsFromPortfolio(portfolioID)
	if err != nil {
		return err
	}

	err = p.addManyAssetsToPortfolio(portfolioID, portfolio.AssetsQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) Delete(portfolio *model.Portfolio) error {
	portfolioID, err := getPortfolioIDByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioID == 0 {
		return ErrPortfolioNotFound
	}

	err = p.deleteAllAssetsFromPortfolio(portfolioID)
	if err != nil {
		return err
	}

	err = deletePortfolio(p.db, portfolioID)

	return err
}

func (p *PostgresPortfolioRepository) DeleteByName(name string) error {
	return p.Delete(&model.Portfolio{Name: name, AssetsQuantityMap: nil})
}

func (p *PostgresPortfolioRepository) GetAll() ([]*model.Portfolio, error) {
	dtoPortfolios, err := getAllPortfolios(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoPortfolios) == 0 {
		return nil, ErrPortfolioNotFound
	}

	var portfolios []*model.Portfolio
	var portfolio *model.Portfolio
	for _, dtoPortfolio := range dtoPortfolios {
		portfolio, err = p.GetByName(dtoPortfolio.Name)
		if err != nil {
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

func (p *PostgresPortfolioRepository) deleteAllAssetsFromPortfolio(portfolioID int) error {
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

func (p *PostgresPortfolioRepository) addManyAssetsToPortfolio(portfolioID int, assetsQuantityMap map[*model.Asset]int) error {
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

func (p *PostgresPortfolioRepository) addAssetToPortfolio(portfolioID, assetID, quantity int) error {
	portfolioAsset, err := createPortfolioAsset(p.db, portfolioID, assetID)
	if err != nil {
		return err
	}

	_, err = createPortfolioAssetRelationship(p.db, portfolioAsset.ID, quantity)

	return err
}

func (p *PostgresPortfolioRepository) convertAssetsQuantityMapToAssetsIDQuantityMap(assetsQuantityMap map[*model.Asset]int) (map[int]int, error) {
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

func (p *PostgresPortfolioRepository) convertAssetsNameQuantityMapToAssetsIDQuantityMap(assetsQuantityMap map[string]int) (map[int]int, error) {
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
