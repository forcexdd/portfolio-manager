package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func CreateNewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	storage := &Storage{db: db}

	// defer closing db connection with error handling
	defer func() {
		closeErr := storage.db.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = storage.db.Ping()
	if err != nil {
		return nil, err
	}

	err = storage.addAllTables()
	if err != nil {
		return nil, err
	}

	return storage, err
}

func DeleteStorage(connString string) error {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	storageToDelete := &Storage{db: db}

	// defer closing db connection with error handling
	defer func() {
		closeErr := storageToDelete.db.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = db.Ping()
	if err != nil {
		return err
	}

	err = storageToDelete.dropAllTables()
	if err != nil {
		return err
	}

	return err
}

func getAllTables() map[string]string {
	return map[string]string{
		"users": `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY
		);`,
		"portfolios": `
		CREATE TABLE IF NOT EXISTS portfolios (
			id SERIAL PRIMARY KEY,
			name VARCHAR(64) NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
		"stocks": `
		CREATE TABLE IF NOT EXISTS stocks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(4) NOT NULL,
			price DECIMAL(20, 10) NOT NULL,
			time TIMESTAMP
		);`,
		"portfolio_stocks": `
		CREATE TABLE IF NOT EXISTS portfolio_stocks (
			id SERIAL PRIMARY KEY,
			portfolio_id INT NOT NULL,
			stock_id INT NOT NULL,
			FOREIGN KEY (portfolio_id) REFERENCES portfolios(id),
			FOREIGN KEY (stock_id) REFERENCES stocks(id)
		);`,
		"portfolio_stocks_relationship": `
		CREATE TABLE IF NOT EXISTS portfolio_stocks_relationship (
			id SERIAL PRIMARY KEY,
			quantity INT NOT NULL,
			portfolio_stocks_id INT NOT NULL,
			FOREIGN KEY (portfolio_stocks_id) REFERENCES portfolio_stocks(id)
		);`,
		"index_stocks": `
		CREATE TABLE IF NOT EXISTS index_stocks (
			id SERIAL PRIMARY KEY,
			name_of_stock VARCHAR(4) NOT NULL,
			fraction DECIMAL(12, 10) NOT NULL,
			time TIMESTAMP
		);`,
	}
}

// Returns an order in which tables must be created
func getTablesOrder() []string {
	return []string{"users", "portfolios", "stocks", "portfolio_stocks", "portfolio_stocks_relationship", "index_stocks"}
}

func (s *Storage) dropAllTables() error {
	for _, table := range getTablesOrder() {
		err := s.DropTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) addAllTables() error {
	for _, table := range getTablesOrder() {
		err := s.CreateTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) DropTable(tableName string) error {
	_, isTable := getAllTables()[tableName]
	if !isTable {
		return errors.New("table not found")
	}

	query := "DROP TABLE IF EXISTS " + tableName + " CASCADE;"
	_, err := s.db.Exec(query)

	return err
}

func (s *Storage) CreateTable(tableName string) error {
	query, isTable := getAllTables()[tableName]
	if !isTable {
		return errors.New("table not found")
	}

	_, err := s.db.Exec(query)

	return err
}
