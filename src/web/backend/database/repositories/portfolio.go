package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/dto_models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
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
		return errors.New("portfolio already exists")
	}

	createdPortfolio, err := createPortfolio(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolio.StocksQuantityMap == nil {
		return nil
	}

	stocksIdQuantityMap := make(map[int]int)
	stocksIdQuantityMap, err = convertStocksQuantityMapToStocksIdQuantityMap(p.db, portfolio.StocksQuantityMap)
	if err != nil {
		return err
	}

	err = addManyStocksToPortfolio(p.db, createdPortfolio.Id, stocksIdQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) GetByName(name string) (*models.Portfolio, error) {
	portfolioId, err := getPortfolioIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if portfolioId == 0 {
		return nil, errors.New("portfolio not found")
	}

	var portfolioStocks []*dto_models.PortfolioStock
	portfolioStocks, err = getAllPortfolioStocksByPortfolioId(p.db, portfolioId)
	if err != nil {
		return nil, err
	}

	stocksQuantityMap := make(map[*models.Stock]int)
	var stock *dto_models.Stock
	var relationship *dto_models.PortfolioStockRelationship
	for _, portfolioStock := range portfolioStocks {
		stock, err = getStock(p.db, portfolioStock.StockId)
		if err != nil {
			return nil, err
		}
		if stock == nil {
			return nil, errors.New("stock not found")
		}

		relationship, err = getPortfolioStockRelationshipByPortfolioStockId(p.db, portfolioStock.Id)
		if err != nil {
			return nil, err
		}
		if relationship == nil {
			return nil, errors.New("relationship not found")
		}

		newStock := models.Stock{
			Name:  stock.Name,
			Price: stock.Price,
		}

		stocksQuantityMap[&newStock] = relationship.Quantity
	}

	return &models.Portfolio{
		Name:              name,
		StocksQuantityMap: stocksQuantityMap,
	}, nil
}

func (p *PostgresPortfolioRepository) Update(portfolio *models.Portfolio) error {
	portfolioId, err := getPortfolioIdByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioId == 0 {
		return errors.New("portfolio not found")
	}

	err = deleteAllStocksFromPortfolio(p.db, portfolioId)
	if err != nil {
		return err
	}

	stocksIdQuantityMap := make(map[int]int)
	stocksIdQuantityMap, err = convertStocksQuantityMapToStocksIdQuantityMap(p.db, portfolio.StocksQuantityMap)
	if err != nil {
		return err
	}

	err = addManyStocksToPortfolio(p.db, portfolioId, stocksIdQuantityMap)

	return err
}

func (p *PostgresPortfolioRepository) Delete(portfolio *models.Portfolio) error {
	portfolioId, err := getPortfolioIdByName(p.db, portfolio.Name)
	if err != nil {
		return err
	}
	if portfolioId == 0 {
		return errors.New("portfolio not found")
	}

	err = deleteAllStocksFromPortfolio(p.db, portfolioId)
	if err != nil {
		return err
	}

	err = deletePortfolio(p.db, portfolioId)

	return err
}

func (p *PostgresPortfolioRepository) DeleteByName(name string) error {
	return p.Delete(&models.Portfolio{Name: name, StocksQuantityMap: nil})
}

func (p *PostgresPortfolioRepository) GetAll() ([]*models.Portfolio, error) {
	dtoPortfolios, err := getAllPortfolios(p.db)
	if err != nil {
		return nil, err
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
			StocksQuantityMap: portfolio.StocksQuantityMap,
		}

		portfolios = append(portfolios, newPortfolio)
	}

	return portfolios, nil
}
