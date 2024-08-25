package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/database/dto_models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
)

type PortfolioRepository interface {
	Create(portfolio *models.Portfolio) error
	GetByName(name string) (*models.Portfolio, error)
	Update(portfolio *models.Portfolio) error
	Delete(portfolio *models.Portfolio) error
	DeleteByName(name string) error
	GetAll() ([]*models.Portfolio, error)
}

type PostgresPortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) PortfolioRepository {
	return &PostgresPortfolioRepository{db: db}
}

func (p *PostgresPortfolioRepository) Create(portfolio *models.Portfolio) error {
	portfolioId, err := getPortfolioIdByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioId != 0 {
		return ErrPortfolioAlreadyExists
	}

	createdPortfolio, err := createPortfolio(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolio.AssetsQuantityMap == nil {
		return nil
	}

	assetsIdQuantityMap := make(map[int]int)
	assetsIdQuantityMap, err = convertAssetsQuantityMapToAssetsIdQuantityMap(p.db, portfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToPortfolio(p.db, createdPortfolio.Id, assetsIdQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) GetByName(name string) (*models.Portfolio, error) {
	portfolioId, err := getPortfolioIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if portfolioId == 0 {
		return nil, ErrPortfolioNotFound
	}

	var portfolioAssets []*dto_models.PortfolioAsset
	portfolioAssets, err = getAllPortfolioAssetsByPortfolioId(p.db, portfolioId)
	if err != nil {
		return nil, err
	}

	assetsQuantityMap := make(map[*models.Asset]int)
	var asset *dto_models.Asset
	var relationship *dto_models.PortfolioAssetRelationship
	for _, portfolioAsset := range portfolioAssets {
		asset, err = getAsset(p.db, portfolioAsset.AssetId)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			return nil, ErrAssetNotFound
		}

		relationship, err = getPortfolioAssetRelationshipByPortfolioAssetId(p.db, portfolioAsset.Id)
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

		assetsQuantityMap[&newAsset] = relationship.Quantity
	}

	return &models.Portfolio{
		Name:              name,
		AssetsQuantityMap: assetsQuantityMap,
	}, nil
}

func (p *PostgresPortfolioRepository) Update(portfolio *models.Portfolio) error {
	portfolioId, err := getPortfolioIdByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioId == 0 {
		return ErrPortfolioNotFound
	}

	err = deleteAllAssetsFromPortfolio(p.db, portfolioId)
	if err != nil {
		return err
	}

	assetsIdQuantityMap := make(map[int]int)
	assetsIdQuantityMap, err = convertAssetsQuantityMapToAssetsIdQuantityMap(p.db, portfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToPortfolio(p.db, portfolioId, assetsIdQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) Delete(portfolio *models.Portfolio) error {
	portfolioId, err := getPortfolioIdByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioId == 0 {
		return ErrPortfolioNotFound
	}

	err = deleteAllAssetsFromPortfolio(p.db, portfolioId)
	if err != nil {
		return err
	}

	err = deletePortfolio(p.db, portfolioId)

	return err
}

func (p *PostgresPortfolioRepository) DeleteByName(name string) error {
	return p.Delete(&models.Portfolio{Name: name, AssetsQuantityMap: nil})
}

func (p *PostgresPortfolioRepository) GetAll() ([]*models.Portfolio, error) {
	dtoPortfolios, err := getAllPortfolios(p.db)
	if err != nil {
		return nil, err
	}
	if len(dtoPortfolios) == 0 {
		return nil, ErrPortfolioNotFound
	}

	var portfolios []*models.Portfolio
	var portfolio *models.Portfolio
	for _, dtoPortfolio := range dtoPortfolios {
		portfolio, err = p.GetByName(dtoPortfolio.Name)
		if err != nil {
			return nil, err
		}

		newPortfolio := &models.Portfolio{
			Name:              portfolio.Name,
			AssetsQuantityMap: portfolio.AssetsQuantityMap,
		}

		portfolios = append(portfolios, newPortfolio)
	}

	return portfolios, nil
}
