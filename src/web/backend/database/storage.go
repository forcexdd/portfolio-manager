package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sql.DB
}

func CreateNewStorage(connString string) (*Storage, error) {
	var newDb Storage
	var err error

	newDb.db, err = sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Checking connection
	err = newDb.db.Ping()
	if err != nil {
		return nil, err
	}

	err = newDb.addAllTables()
	if err != nil {
		return nil, err
	}

	err = newDb.db.Close()
	if err != nil {
		return nil, err
	}

	return &newDb, nil
}

func DeleteStorage(connString string) error {
	var dbToDelete Storage
	var err error

	dbToDelete.db, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	// Checking connection
	err = dbToDelete.db.Ping()
	if err != nil {
		return err
	}

	err = dbToDelete.dropAllTables()
	if err != nil {
		return err
	}

	err = dbToDelete.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) dropAllTables() error {
	err := s.deleteUsersTable()
	if err != nil {
		return err
	}

	err = s.deletePortfoliosTable()
	if err != nil {
		return err
	}

	err = s.deleteStocksTable()
	if err != nil {
		return err
	}

	err = s.deletePortfolioStocksTable()
	if err != nil {
		return err
	}

	err = s.deletePortfolioStocksRelationshipTable()
	if err != nil {
		return err
	}

	err = s.deleteIndexStockTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) addAllTables() error {
	err := s.createUsersTable()
	if err != nil {
		return err
	}

	err = s.createPortfoliosTable()
	if err != nil {
		return err
	}

	err = s.createStocksTable()
	if err != nil {
		return err
	}

	err = s.createPortfolioStocksTable()
	if err != nil {
		return err
	}

	err = s.createPortfolioStocksRelationshipTable()
	if err != nil {
		return err
	}

	err = s.createIndexStockTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createUsersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteUsersTable() error {
	query := "DROP TABLE IF EXISTS users CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createPortfoliosTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS portfolios (
		id SERIAL PRIMARY KEY,
	    name VARCHAR(64) NOT NULL,
		user_id INT NOT NULL,
	    FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deletePortfoliosTable() error {
	query := "DROP TABLE IF EXISTS portfolios CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createStocksTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS stocks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(4) NOT NULL,
		price DECIMAL(20, 10) NOT NULL,
		time TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteStocksTable() error {
	query := "DROP TABLE IF EXISTS stocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createPortfolioStocksTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS portfolio_stocks (
		id SERIAL PRIMARY KEY,
		portfolio_id INT NOT NULL,
		stock_id INT NOT NULL,
		FOREIGN KEY (portfolio_id) REFERENCES portfolios(id),
		FOREIGN KEY (stock_id) REFERENCES stocks(id)
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deletePortfolioStocksTable() error {
	query := "DROP TABLE IF EXISTS portfolio_stocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createPortfolioStocksRelationshipTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS portfolio_stocks_relationship (
		id SERIAL PRIMARY KEY,
		quantity INT NOT NULL,
		portfolio_stocks_id INT NOT NULL,
		FOREIGN KEY (portfolio_stocks_id) REFERENCES portfolio_stocks(id)
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deletePortfolioStocksRelationshipTable() error {
	query := "DROP TABLE IF EXISTS portfolio_stocks_relationship CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createIndexStockTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS index_stocks (
		id SERIAL PRIMARY KEY,
		name_of_stock VARCHAR(4) NOT NULL,
		fraction DECIMAL(12, 10) NOT NULL,
		time TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteIndexStockTable() error {
	query := "DROP TABLE IF EXISTS index_stocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
