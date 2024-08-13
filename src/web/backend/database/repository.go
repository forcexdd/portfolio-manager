package database

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	_ "github.com/lib/pq"
)

/*
	BASIC CRUD
*/

/*
	Portfolio
*/

func (s *Storage) CreatePortfolio(name string) (*models.Portfolio, error) {
	query := `INSERT INTO portfolios (name) VALUES ($1) RETURNING id;`

	var portfolioID int
	err := s.db.QueryRow(query, name).Scan(&portfolioID)
	if err != nil {
		return nil, err
	}

	return &models.Portfolio{
		Id:   portfolioID,
		Name: name,
	}, nil
}

func (s *Storage) GetPortfolio(portfolioID int) (*models.Portfolio, error) {
	query := `SELECT id, name FROM portfolios WHERE id = $1;`

	var portfolio models.Portfolio
	err := s.db.QueryRow(query, portfolioID).Scan(&portfolio.Id, &portfolio.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolio, nil
}

func (s *Storage) UpdatePortfolio(portfolioID int, newName string) error {
	query := `UPDATE portfolios SET name = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, newName, portfolioID)

	return err
}

func (s *Storage) DeletePortfolio(portfolioID int) error {
	query := `DELETE FROM portfolios WHERE id = $1;`
	_, err := s.db.Exec(query, portfolioID)

	return err
}

func (s *Storage) GetPortfolioIDByName(name string) (int, error) {
	query := `SELECT id FROM portfolios WHERE name = $1;`

	var portfolioID int
	err := s.db.QueryRow(query, name).Scan(&portfolioID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return portfolioID, nil
}

/*
	Stock
*/

func (s *Storage) CreateStock(name string, price float64) (*models.Stock, error) {
	query := `INSERT INTO stocks (name, price) VALUES ($1, $2) RETURNING id;`

	var stockID int
	err := s.db.QueryRow(query, name, price).Scan(&stockID)
	if err != nil {
		return nil, err
	}

	return &models.Stock{
		Id:    stockID,
		Name:  name,
		Price: price,
	}, nil
}

func (s *Storage) GetStock(stockID int) (*models.Stock, error) {
	query := `SELECT id, name, price FROM stocks WHERE id = $1;`

	var stock models.Stock
	err := s.db.QueryRow(query, stockID).Scan(&stock.Id, &stock.Name, &stock.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &stock, nil
}

func (s *Storage) UpdateStock(stockID int, newPrice float64) error {
	query := `UPDATE stocks SET price = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, newPrice, stockID)

	return err
}

func (s *Storage) DeleteStock(stockID int) error {
	query := `DELETE FROM stocks WHERE id = $1;`
	_, err := s.db.Exec(query, stockID)

	return err
}

func (s *Storage) GetStockIDByName(name string) (int, error) {
	query := `SELECT id FROM stocks WHERE name = $1;`
	var stockID int
	err := s.db.QueryRow(query, name).Scan(&stockID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return stockID, nil
}

/*
	PortfolioStock
*/

// Returns PortfolioStockID
func (s *Storage) CreatePortfolioStock(portfolioID int, stockID int) (*models.PortfolioStock, error) {
	query := `INSERT INTO portfolio_stocks (portfolio_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var portfolioStockID int
	err := s.db.QueryRow(query, portfolioID, stockID).Scan(&portfolioStockID)
	if err != nil {
		return nil, err
	}

	return &models.PortfolioStock{
		Id:          portfolioStockID,
		PortfolioId: portfolioID,
		StockId:     stockID,
	}, nil
}

func (s *Storage) GetPortfolioStock(portfolioStockID int) (*models.PortfolioStock, error) {
	query := `SELECT id, portfolio_id, stock_id FROM portfolio_stocks WHERE id = $1;`

	var portfolioStock models.PortfolioStock
	err := s.db.QueryRow(query, portfolioStockID).Scan(&portfolioStock.Id, &portfolioStock.PortfolioId, &portfolioStock.StockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolioStock, nil
}

func (s *Storage) DeletePortfolioStock(portfolioStockID int) error {
	query := `DELETE FROM portfolio_stocks WHERE id = $1;`
	_, err := s.db.Exec(query, portfolioStockID)

	return err
}

/*
	PortfolioStockRelationship
*/

func (s *Storage) CreatePortfolioStockRelationship(portfolioStockID int, quantity int) (*models.PortfolioStockRelationship, error) {
	query := `INSERT INTO portfolio_stocks_relationship (portfolio_stocks_id, quantity) VALUES ($1, $2) RETURNING id;`

	var portfolioStockRelationshipID int
	err := s.db.QueryRow(query, portfolioStockID, quantity).Scan(&portfolioStockRelationshipID)
	if err != nil {
		return nil, err
	}

	return &models.PortfolioStockRelationship{
		Id:               portfolioStockRelationshipID,
		PortfolioStockId: portfolioStockID,
		Quantity:         quantity,
	}, nil
}

func (s *Storage) GetPortfolioStockRelationship(relationshipID int) (*models.PortfolioStockRelationship, error) {
	query := `SELECT id, portfolio_stocks_id, quantity FROM portfolio_stocks_relationship WHERE id = $1;`

	var relationship models.PortfolioStockRelationship
	err := s.db.QueryRow(query, relationshipID).Scan(&relationship.Id, &relationship.PortfolioStockId, &relationship.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func (s *Storage) UpdatePortfolioStockRelationship(relationshipID int, newQuantity int) error {
	query := `UPDATE portfolio_stocks_relationship SET quantity = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, newQuantity, relationshipID)

	return err
}

func (s *Storage) DeletePortfolioStockRelationship(relationshipID int) error {
	query := `DELETE FROM portfolio_stocks_relationship WHERE id = $1;`
	_, err := s.db.Exec(query, relationshipID)

	return err
}

/*
	Index
*/

func (s *Storage) CreateIndex(name string) (*models.Index, error) {
	query := `INSERT INTO indexes (name) VALUES ($1) RETURNING id;`

	var indexID int
	err := s.db.QueryRow(query, name).Scan(&indexID)
	if err != nil {
		return nil, err
	}

	return &models.Index{
		Id:   indexID,
		Name: name,
	}, nil
}

func (s *Storage) GetIndex(indexID int) (*models.Index, error) {
	query := `SELECT id, name FROM indexes WHERE id = $1;`

	var index models.Index
	err := s.db.QueryRow(query, indexID).Scan(&index.Id, &index.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &index, nil
}

func (s *Storage) UpdateIndex(indexID int, newName string) error {
	query := `UPDATE indexes SET name = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, newName, indexID)

	return err
}

func (s *Storage) DeleteIndex(indexID int) error {
	query := `DELETE FROM indexes WHERE id = $1;`
	_, err := s.db.Exec(query, indexID)

	return err
}

/*
	IndexStock
*/

func (s *Storage) CreateIndexStock(indexID int, stockID int) (*models.IndexStock, error) {
	query := `INSERT INTO index_stocks (index_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var indexStockID int
	err := s.db.QueryRow(query, indexID, stockID).Scan(&indexStockID)
	if err != nil {
		return nil, err
	}

	return &models.IndexStock{
		Id:      indexStockID,
		IndexId: indexID,
		StockId: stockID,
	}, nil
}

func (s *Storage) GetIndexStock(indexStockID int) (*models.IndexStock, error) {
	query := `SELECT id, index_id, stock_id FROM index_stocks WHERE id = $1;`

	var indexStock models.IndexStock
	err := s.db.QueryRow(query, indexStockID).Scan(&indexStock.Id, &indexStock.IndexId, &indexStock.StockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &indexStock, nil
}

func (s *Storage) DeleteIndexStock(indexStockID int) error {
	query := `DELETE FROM index_stocks WHERE id = $1;`
	_, err := s.db.Exec(query, indexStockID)

	return err
}

/*
	IndexStockRelationship
*/

func (s *Storage) CreateIndexStockRelationship(indexStockID int, fraction float64) (*models.IndexStockRelationship, error) {
	query := `INSERT INTO index_stocks_relationship (index_stocks_id, fraction) VALUES ($1, $2) RETURNING id;`

	var relationshipID int
	err := s.db.QueryRow(query, indexStockID, fraction).Scan(&relationshipID)
	if err != nil {
		return nil, err
	}

	return &models.IndexStockRelationship{
		Id:           relationshipID,
		IndexStockId: indexStockID,
		Fraction:     fraction,
	}, nil
}

func (s *Storage) GetIndexStockRelationship(relationshipID int) (*models.IndexStockRelationship, error) {
	query := `SELECT id, index_stocks_id, fraction FROM index_stocks_relationship WHERE id = $1;`

	var relationship models.IndexStockRelationship
	err := s.db.QueryRow(query, relationshipID).Scan(&relationship.Id, &relationship.IndexStockId, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func (s *Storage) UpdateIndexStockRelationship(relationshipID int, newFraction float64) error {
	query := `UPDATE index_stocks_relationship SET fraction = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, newFraction, relationshipID)

	return err
}

func (s *Storage) DeleteIndexStockRelationship(relationshipID int) error {
	query := `DELETE FROM index_stocks_relationship WHERE id = $1;`
	_, err := s.db.Exec(query, relationshipID)

	return err
}

/*
	ADDITIONAL FUNCTIONALITY
*/

func (s *Storage) AddStockToPortfolio(portfolioID, stockID, quantity int) error {
	portfolioStock, err := s.CreatePortfolioStock(portfolioID, stockID)
	if err != nil {
		return err
	}

	_, err = s.CreatePortfolioStockRelationship(portfolioStock.Id, quantity)

	return err
}

func (s *Storage) AddManyStocksToPortfolio(portfolioID int, stocksQuantityMap map[int]int) error {
	for stockID, quantity := range stocksQuantityMap {
		err := s.AddStockToPortfolio(portfolioID, stockID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) convertStocksQuantityMapToStocksIDQuantityMap(stocksQuantityMap map[string]int) (map[int]int, error) {
	stocksIDQuantityMap := make(map[int]int)
	for stockName, quantity := range stocksQuantityMap {
		stockID, err := s.GetStockIDByName(stockName)
		if err != nil {
			return nil, err
		}

		stocksIDQuantityMap[stockID] = quantity
	}

	return stocksIDQuantityMap, nil
}

func (s *Storage) AddManyStocksToPortfolioByName(portfolioName string, stocksQuantityMap map[string]int) error {
	portfolioID, err := s.GetPortfolioIDByName(portfolioName)
	if err != nil {
		return err
	}

	stocksIDQuantityMap := make(map[int]int)
	stocksIDQuantityMap, err = s.convertStocksQuantityMapToStocksIDQuantityMap(stocksQuantityMap)
	if err != nil {
		return err
	}

	err = s.AddManyStocksToPortfolio(portfolioID, stocksIDQuantityMap)

	return err
}

func (s *Storage) GetAllPortfolios() ([]*models.Portfolio, error) {
	query := `SELECT id, name FROM portfolios;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolios []*models.Portfolio
	var portfolio models.Portfolio
	for rows.Next() {
		err = rows.Scan(&portfolio.Id, &portfolio.Name)
		if err != nil {
			return nil, err
		}

		portfolios = append(portfolios, &portfolio)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (s *Storage) DeleteAllStocksFromPortfolio(portfolioID int) error {
	query := `
        DELETE FROM portfolio_stocks_relationship
        WHERE portfolio_stocks_id IN (
            SELECT id FROM portfolio_stocks WHERE portfolio_id = $1
        );
    `
	_, err := s.db.Exec(query, portfolioID)
	if err != nil {
		return err
	}

	query = `DELETE FROM portfolio_stocks WHERE portfolio_id = $1;`
	_, err = s.db.Exec(query, portfolioID)

	return err
}

func (s *Storage) DeleteAllStocksFromPortfolioByName(portfolioName string) error {
	portfolioID, err := s.GetPortfolioIDByName(portfolioName)
	if err != nil {
		return err
	}

	err = s.DeleteAllStocksFromPortfolio(portfolioID)

	return err
}
