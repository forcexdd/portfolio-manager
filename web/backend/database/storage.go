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
	err := s.deleteStocksTable()
	if err != nil {
		return err
	}

	err = s.deleteUserPortfolioTable()
	if err != nil {
		return err
	}

	err = s.deleteIndexStockTable()
	if err != nil {
		return err
	}

	err = s.deleteManyUserPortfolioWithManyStocksTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) addAllTables() error {
	err := s.createStocksTable()
	if err != nil {
		return err
	}

	err = s.createUserPortfolioTable()
	if err != nil {
		return err
	}

	err = s.createIndexStockTable()
	if err != nil {
		return err
	}

	err = s.createManyUserPortfolioWithManyStocksTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createStocksTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS Stocks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(5) NOT NULL,
		price DECIMAL(20, 10) NOT NULL,
		time TIMESTAMP NOT NULL
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteStocksTable() error {
	query := "DROP TABLE IF EXISTS Stocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createUserPortfolioTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS UserPortfolios (
		id SERIAL PRIMARY KEY,
	    name VARCHAR(256) NOT NULL
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteUserPortfolioTable() error {
	query := "DROP TABLE IF EXISTS UserPortfolios CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createIndexStockTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS IndexStocks (
		id SERIAL PRIMARY KEY,
		name_of_stock VARCHAR(5) NOT NULL,
		fraction DECIMAL(12, 10) NOT NULL,
		time TIMESTAMP NOT NULL
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteIndexStockTable() error {
	query := "DROP TABLE IF EXISTS IndexStocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) createManyUserPortfolioWithManyStocksTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS ManyUserPortfolioWithManyStocks (
		id SERIAL PRIMARY KEY,
		portfolio_id INT NOT NULL,
		stock_id INT NOT NULL,
		quantity INT NOT NULL,
		FOREIGN KEY (portfolio_id) REFERENCES UserPortfolios(id),
		FOREIGN KEY (stock_id) REFERENCES Stocks(id)
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Storage) deleteManyUserPortfolioWithManyStocksTable() error {
	query := "DROP TABLE IF EXISTS ManyUserPortfolioWithManyStocks CASCADE;"

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
