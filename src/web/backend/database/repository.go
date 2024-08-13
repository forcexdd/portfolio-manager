package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

/*
	BASIC CRUD
*/

/*
	Portfolio
*/

// Returns portfolio ID
func (s *Storage) CreatePortfolio(name string) (int, error) {
	query := `INSERT INTO portfolios (name) VALUES ($1) RETURNING id;`

	var portfolioID int
	err := s.db.QueryRow(query, name).Scan(&portfolioID)

	if err != nil {
		return 0, err
	}

	return portfolioID, nil
}

// Returns portfolio name
func (s *Storage) GetPortfolio(portfolioID int) (string, error) {
	query := `SELECT name FROM portfolios WHERE id = $1;`

	var name string
	err := s.db.QueryRow(query, portfolioID).Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No rows found
		}
		return "", err
	}

	return name, nil
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
		if err == sql.ErrNoRows {
			return 0, nil // No rows found
		}
		return 0, err
	}

	return portfolioID, nil
}

/*
	Stock
*/

// Returns stock id
func (s *Storage) CreateStock(name string, price float64) (int, error) {
	query := `INSERT INTO stocks (name, price) VALUES ($1, $2) RETURNING id;`

	var stockID int
	err := s.db.QueryRow(query, name, price).Scan(&stockID)

	if err != nil {
		return 0, err
	}

	return stockID, nil
}

// Returns stock name, price
func (s *Storage) GetStock(stockID int) (string, float64, error) {
	query := `SELECT name, price FROM stocks WHERE id = $1;`

	var name string
	var price float64
	err := s.db.QueryRow(query, stockID).Scan(&name, &price)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, nil // No rows found
		}
		return "", 0, err
	}

	return name, price, nil
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
		if err == sql.ErrNoRows {
			return 0, nil // No rows found
		}
		return 0, err
	}

	return stockID, nil
}

/*
	PortfolioStock
*/

// Returns PortfolioStockID
func (s *Storage) CreatePortfolioStock(portfolioID int, stockID int) (int, error) {
	query := `INSERT INTO portfolio_stocks (portfolio_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var portfolioStockID int
	err := s.db.QueryRow(query, portfolioID, stockID).Scan(&portfolioStockID)

	if err != nil {
		return 0, err
	}

	return portfolioStockID, nil
}

// Returns portfolio ID & stock ID
func (s *Storage) GetPortfolioStock(portfolioStockID int) (int, int, error) {
	query := `SELECT portfolio_id, stock_id FROM portfolio_stocks WHERE id = $1;`

	var portfolioID, stockID int
	err := s.db.QueryRow(query, portfolioStockID).Scan(&portfolioID, &stockID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil // No rows found
		}
		return 0, 0, err
	}

	return portfolioID, stockID, nil
}

func (s *Storage) DeletePortfolioStock(portfolioStockID int) error {
	query := `DELETE FROM portfolio_stocks WHERE id = $1;`
	_, err := s.db.Exec(query, portfolioStockID)

	return err
}

/*
	PortfolioStockRelationship
*/

// Returns relationship ID
func (s *Storage) CreatePortfolioStockRelationship(portfolioStockID int, quantity int) (int, error) {
	query := `INSERT INTO portfolio_stocks_relationship (portfolio_stocks_id, quantity) VALUES ($1, $2) RETURNING id;`

	var portfolioStockRelationshipID int
	err := s.db.QueryRow(query, portfolioStockID, quantity).Scan(&portfolioStockRelationshipID)

	if err != nil {
		return 0, err
	}

	return portfolioStockRelationshipID, nil
}

// Returns PortfolioStockID & quantity
func (s *Storage) GetPortfolioStockRelationship(relationshipID int) (int, int, error) {
	query := `SELECT portfolio_stocks_id, quantity FROM portfolio_stocks_relationship WHERE id = $1;`

	var portfolioStockID, quantity int
	err := s.db.QueryRow(query, relationshipID).Scan(&portfolioStockID, &quantity)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil // No rows found
		}
		return 0, 0, err
	}

	return portfolioStockID, quantity, nil
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

// Returns indexID
func (s *Storage) CreateIndex(name string) (int, error) {
	query := `INSERT INTO indexes (name) VALUES ($1) RETURNING id;`

	var indexID int
	err := s.db.QueryRow(query, name).Scan(&indexID)

	if err != nil {
		return 0, err
	}

	return indexID, nil
}

// Returns name of the index
func (s *Storage) GetIndex(indexID int) (string, error) {
	query := `SELECT name FROM indexes WHERE id = $1;`

	var name string
	err := s.db.QueryRow(query, indexID).Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No rows found
		}
		return "", err
	}

	return name, nil
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

// Returns index stock ID
func (s *Storage) CreateIndexStock(indexID int, stockID int) (int, error) {
	query := `INSERT INTO index_stocks (index_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var indexStockID int
	err := s.db.QueryRow(query, indexID, stockID).Scan(&indexStockID)

	if err != nil {
		return 0, err
	}

	return indexStockID, nil
}

// Returns indexID, stockID
func (s *Storage) GetIndexStock(indexStockID int) (int, int, error) {

	query := `SELECT index_id, stock_id FROM index_stocks WHERE id = $1;`

	var indexID, stockID int
	err := s.db.QueryRow(query, indexStockID).Scan(&indexID, &stockID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil // No rows found
		}
		return 0, 0, err
	}

	return indexID, stockID, nil
}

func (s *Storage) DeleteIndexStock(indexStockID int) error {
	query := `DELETE FROM index_stocks WHERE id = $1;`
	_, err := s.db.Exec(query, indexStockID)

	return err
}

/*
	IndexStockRelationship
*/

// Returns relationshipID
func (s *Storage) CreateIndexStockRelationship(indexStockID int, fraction float64) (int, error) {
	query := `INSERT INTO index_stocks_relationship (index_stocks_id, fraction) VALUES ($1, $2) RETURNING id;`

	var relationshipID int
	err := s.db.QueryRow(query, indexStockID, fraction).Scan(&relationshipID)

	if err != nil {
		return 0, err
	}

	return relationshipID, nil
}

// Returns indexStockID, fraction
func (s *Storage) GetIndexStockRelationship(relationshipID int) (int, float64, error) {
	query := `SELECT index_stocks_id, fraction FROM index_stocks_relationship WHERE id = $1;`

	var indexStockID int
	var fraction float64
	err := s.db.QueryRow(query, relationshipID).Scan(&indexStockID, &fraction)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil // No rows found
		}
		return 0, 0, err
	}

	return indexStockID, fraction, nil
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
	portfolioStockID, err := s.CreatePortfolioStock(portfolioID, stockID)
	if err != nil {
		return err
	}

	_, err = s.CreatePortfolioStockRelationship(portfolioStockID, quantity)

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
