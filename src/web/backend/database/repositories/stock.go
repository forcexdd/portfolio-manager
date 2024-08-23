package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/dto_models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
)

type StockRepository interface {
	Create(stock *models.Stock) error
	GetByName(name string) (*models.Stock, error)
	Update(stock *models.Stock) error
	Delete(stock *models.Stock) error
	DeleteByName(name string) error
	GetAll() ([]*models.Stock, error)
}

type PostgresStockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &PostgresStockRepository{db: db}
}

func (p *PostgresStockRepository) Create(stock *models.Stock) error {
	stockId, err := getStockIdByName(p.db, stock.Name)
	if err != nil {
		return err
	}
	if stockId != 0 {
		return errors.New("stock already exists")
	}

	_, err = createStock(p.db, stock.Name, stock.Price)

	return err
}

func (p *PostgresStockRepository) GetByName(name string) (*models.Stock, error) {
	stockId, err := getStockIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if stockId == 0 {
		return nil, errors.New("stock not found")
	}

	var dtoStock *dto_models.Stock
	dtoStock, err = getStock(p.db, stockId)
	if err != nil {
		return nil, err
	}

	return &models.Stock{
		Name:  dtoStock.Name,
		Price: dtoStock.Price,
	}, nil
}

func (p *PostgresStockRepository) Update(stock *models.Stock) error {
	stockId, err := getStockIdByName(p.db, stock.Name)
	if err != nil {
		return err
	}
	if stockId == 0 {
		return errors.New("stock not found")
	}

	var dtoStock *dto_models.Stock
	dtoStock, err = getStock(p.db, stockId)
	if err != nil {
		return err
	}
	if dtoStock.Name != stock.Name {
		return errors.New("stock name does not match")
	}

	err = updateStock(p.db, stockId, stock.Price)

	return err
}

func (p *PostgresStockRepository) Delete(stock *models.Stock) error {
	stockId, err := getStockIdByName(p.db, stock.Name)
	if err != nil {
		return err
	}
	if stockId == 0 {
		return errors.New("stock not found")
	}

	err = deleteStockFromConnectedTables(p.db, stockId)
	if err != nil {
		return err
	}

	err = deleteStock(p.db, stockId)

	return err
}

func (p *PostgresStockRepository) DeleteByName(name string) error {
	return p.Delete(&models.Stock{Name: name}) // Price doesn't matter
}

func (p *PostgresStockRepository) GetAll() ([]*models.Stock, error) {
	dtoStocks, err := getAllStocks(p.db)
	if err != nil {
		return nil, err
	}

	var stocks []*models.Stock
	var stock models.Stock
	for _, dtoStock := range dtoStocks {
		stock.Name = dtoStock.Name
		stock.Price = dtoStock.Price

		newStock := &models.Stock{
			Name:  stock.Name,
			Price: stock.Price,
		}

		stocks = append(stocks, newStock)
	}

	return stocks, nil
}
