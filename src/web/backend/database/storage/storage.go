package storage

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db          *sql.DB
	connStr     string
	allTables   map[string]string
	tablesOrder []string
}

func NewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		Db:          db,
		connStr:     connString,
		allTables:   getAllTables(),
		tablesOrder: getTablesOrder(),
	}

	err = storage.Db.Ping()
	if err != nil {
		return nil, err
	}

	err = storage.addAllTables()
	if err != nil {
		return nil, err
	}

	return storage, err
}

func (s *Storage) DeleteStorage() error {
	var err error

	// defer closing db connection with error handling
	defer func() {
		closeErr := s.Db.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = s.Db.Ping()
	if err != nil {
		return err
	}

	err = s.dropAllTables()
	if err != nil {
		return err
	}

	return err
}

func getAllTables() map[string]string {
	return map[string]string{
		"portfolios": `
		CREATE TABLE IF NOT EXISTS portfolios (
			id SERIAL PRIMARY KEY,
			name VARCHAR(64) NOT NULL UNIQUE
		);`,
		"stocks": `
		CREATE TABLE IF NOT EXISTS stocks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(4) NOT NULL UNIQUE,
			price DECIMAL(32, 15) NOT NULL
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
		"indexes": `
		CREATE TABLE IF NOT EXISTS indexes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(64) NOT NULL UNIQUE
		);`,
		"index_stocks": `
		CREATE TABLE IF NOT EXISTS index_stocks (
			id SERIAL PRIMARY KEY,
			index_id INT NOT NULL,
			stock_id INT NOT NULL,
			FOREIGN KEY (index_id) REFERENCES indexes(id),
			FOREIGN KEY (stock_id) REFERENCES stocks(id)
		);`,
		"index_stocks_relationship": `
		CREATE TABLE IF NOT EXISTS index_stocks_relationship (
			id SERIAL PRIMARY KEY,
			fraction DECIMAL(17, 15) NOT NULL,
			index_stocks_id INT NOT NULL,
			FOREIGN KEY (index_stocks_id) REFERENCES index_stocks(id)
		);`,
	}
}

// Returns an order in which tables must be created
func getTablesOrder() []string {
	return []string{"portfolios", "stocks", "portfolio_stocks", "portfolio_stocks_relationship", "indexes", "index_stocks", "index_stocks_relationship"}
}

func (s *Storage) dropAllTables() error {
	for _, table := range s.tablesOrder {
		err := s.DropTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) addAllTables() error {
	for _, table := range s.tablesOrder {
		err := s.CreateTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) DropTable(tableName string) error {
	_, isTable := s.allTables[tableName]
	if !isTable {
		return errors.New("table not found")
	}

	query := "DROP TABLE IF EXISTS " + tableName + " CASCADE;"
	_, err := s.Db.Exec(query)

	return err
}

func (s *Storage) CreateTable(tableName string) error {
	query, isTable := s.allTables[tableName]
	if !isTable {
		return errors.New("table not found")
	}

	_, err := s.Db.Exec(query)

	return err
}
