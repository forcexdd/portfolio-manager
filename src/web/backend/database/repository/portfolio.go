package repository

import (
	"database/sql"
	"errors"
	dtomodels "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

type PortfolioRepository interface {
	Create(portfolio *model.Portfolio) error
	GetByName(name string) (*model.Portfolio, error)
	Update(portfolio *model.Portfolio) error
	Delete(portfolio *model.Portfolio) error
	DeleteByName(name string) error
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

	createdPortfolio, err := createPortfolio(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolio.AssetsQuantityMap == nil {
		return nil
	}

	assetsIDQuantityMap := make(map[int]int)
	assetsIDQuantityMap, err = convertAssetsQuantityMapToAssetsIDQuantityMap(p.db, portfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToPortfolio(p.db, createdPortfolio.ID, assetsIDQuantityMap)

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

	var portfolioAssets []*dtomodels.PortfolioAsset
	portfolioAssets, err = getAllPortfolioAssetsByPortfolioID(p.db, portfolioID)
	if err != nil {
		return nil, err
	}

	assetsQuantityMap := make(map[*model.Asset]int)
	var asset *dtomodels.Asset
	var relationship *dtomodels.PortfolioAssetRelationship
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

	err = deleteAllAssetsFromPortfolio(p.db, portfolioID)
	if err != nil {
		return err
	}

	assetsIDQuantityMap := make(map[int]int)
	assetsIDQuantityMap, err = convertAssetsQuantityMapToAssetsIDQuantityMap(p.db, portfolio.AssetsQuantityMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToPortfolio(p.db, portfolioID, assetsIDQuantityMap)

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

	err = deleteAllAssetsFromPortfolio(p.db, portfolioID)
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
