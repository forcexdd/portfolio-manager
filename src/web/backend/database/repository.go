package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

/*
BASIC CRUD
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

// Returns stock id
func (s *Storage) CreateStock(name string, price float64, time string) (int, error) {
	query := `INSERT INTO stocks (name, price, time) VALUES ($1, $2, $3) RETURNING id;`

	var stockID int
	err := s.db.QueryRow(query, name, price, time).Scan(&stockID)

	if err != nil {
		return 0, err
	}

	return stockID, nil
}

// Returns stock name, price and time it's created
func (s *Storage) GetStock(stockID int) (string, float64, string, error) {
	query := `SELECT name, price, time FROM stocks WHERE id = $1;`

	var name string
	var price float64
	var time string
	err := s.db.QueryRow(query, stockID).Scan(&name, &price, &time)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, "", nil // No rows found
		}
		return "", 0, "", err
	}

	return name, price, time, nil
}

func (s *Storage) UpdateStock(stockID int, newPrice float64, newTime string) error {
	query := `UPDATE stocks SET price = $1, time = $2 WHERE id = $3;`
	_, err := s.db.Exec(query, newPrice, newTime, stockID)

	return err
}

func (s *Storage) DeleteStock(stockID int) error {
	query := `DELETE FROM stocks WHERE id = $1;`
	_, err := s.db.Exec(query, stockID)

	return err
}

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

// Returns index stock ID
func (s *Storage) CreateIndexStock(name string, fraction float64, time string) (int, error) {
	query := `INSERT INTO index_stocks (name_of_stock, fraction, time) VALUES ($1, $2, $3) RETURNING id;`

	var indexStockID int
	err := s.db.QueryRow(query, name, fraction, time).Scan(&indexStockID)

	if err != nil {
		return 0, err
	}

	return indexStockID, nil
}

// Returns stock name, fraction and time stock has been created
func (s *Storage) GetIndexStock(indexStockID int) (string, float64, string, error) {
	query := `SELECT name_of_stock, fraction, time FROM index_stocks WHERE id = $1;`

	var name string
	var fraction float64
	var time string
	err := s.db.QueryRow(query, indexStockID).Scan(&name, &fraction, &time)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, "", nil // No rows found
		}
		return "", 0, "", err
	}

	return name, fraction, time, nil
}

func (s *Storage) UpdateIndexStock(indexStockID int, newFraction float64, newTime string) error {
	query := `UPDATE index_stocks SET fraction = $1, time = $2 WHERE id = $3;`
	_, err := s.db.Exec(query, newFraction, newTime, indexStockID)

	return err
}

func (s *Storage) DeleteIndexStock(indexStockID int) error {
	query := `DELETE FROM index_stocks WHERE id = $1;`
	_, err := s.db.Exec(query, indexStockID)

	return err
}
