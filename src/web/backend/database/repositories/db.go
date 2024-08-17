package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/dto_models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
)

/*
	Portfolio
*/

func createPortfolio(db *sql.DB, name string) (*dto_models.Portfolio, error) {
	query := `INSERT INTO portfolios (name) VALUES ($1) RETURNING id;`

	var portfolioId int
	err := db.QueryRow(query, name).Scan(&portfolioId)
	if err != nil {
		return nil, err
	}

	return &dto_models.Portfolio{
		Id:   portfolioId,
		Name: name,
	}, nil
}

func getPortfolio(db *sql.DB, portfolioId int) (*dto_models.Portfolio, error) {
	query := `SELECT id, name FROM portfolios WHERE id = $1;`

	var portfolio dto_models.Portfolio
	err := db.QueryRow(query, portfolioId).Scan(&portfolio.Id, &portfolio.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolio, nil
}

func updatePortfolio(db *sql.DB, portfolioId int, newName string) error {
	query := `UPDATE portfolios SET name = $1 WHERE id = $2;`
	_, err := db.Exec(query, newName, portfolioId)

	return err
}

func deletePortfolio(db *sql.DB, portfolioId int) error {
	query := `DELETE FROM portfolios WHERE id = $1;`
	_, err := db.Exec(query, portfolioId)

	return err
}

func getPortfolioIdByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM portfolios WHERE name = $1;`

	var portfolioId int
	err := db.QueryRow(query, name).Scan(&portfolioId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return portfolioId, nil
}

func getAllPortfolios(db *sql.DB) ([]*dto_models.Portfolio, error) {
	query := `SELECT id, name FROM portfolios;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolios []*dto_models.Portfolio
	var portfolio dto_models.Portfolio
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

/*
	Stock
*/

func createStock(db *sql.DB, name string, price float64) (*dto_models.Stock, error) {
	query := `INSERT INTO stocks (name, price) VALUES ($1, $2) RETURNING id;`

	var stockId int
	err := db.QueryRow(query, name, price).Scan(&stockId)
	if err != nil {
		return nil, err
	}

	return &dto_models.Stock{
		Id:    stockId,
		Name:  name,
		Price: price,
	}, nil
}

func getStock(db *sql.DB, stockId int) (*dto_models.Stock, error) {
	query := `SELECT id, name, price FROM stocks WHERE id = $1;`

	var stock dto_models.Stock
	err := db.QueryRow(query, stockId).Scan(&stock.Id, &stock.Name, &stock.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &stock, nil
}

func updateStock(db *sql.DB, stockId int, newPrice float64) error {
	query := `UPDATE stocks SET price = $1 WHERE id = $2;`
	_, err := db.Exec(query, newPrice, stockId)

	return err
}

func deleteStock(db *sql.DB, stockId int) error {
	query := `DELETE FROM stocks WHERE id = $1;`
	_, err := db.Exec(query, stockId)

	return err
}

func getStockIdByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM stocks WHERE name = $1;`
	var stockId int
	err := db.QueryRow(query, name).Scan(&stockId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return stockId, nil
}

func getAllStocks(db *sql.DB) ([]*dto_models.Stock, error) {
	query := `SELECT id, name, price FROM stocks;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var stocks []*dto_models.Stock
	var stock dto_models.Stock
	for rows.Next() {
		err = rows.Scan(&stock.Id, &stock.Name, &stock.Price)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, &stock)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

/*
	PortfolioStock
*/

func createPortfolioStock(db *sql.DB, portfolioId int, stockId int) (*dto_models.PortfolioStock, error) {
	query := `INSERT INTO portfolio_stocks (portfolio_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var portfolioStockId int
	err := db.QueryRow(query, portfolioId, stockId).Scan(&portfolioStockId)
	if err != nil {
		return nil, err
	}

	return &dto_models.PortfolioStock{
		Id:          portfolioStockId,
		PortfolioId: portfolioId,
		StockId:     stockId,
	}, nil
}

func getPortfolioStock(db *sql.DB, portfolioStockId int) (*dto_models.PortfolioStock, error) {
	query := `SELECT id, portfolio_id, stock_id FROM portfolio_stocks WHERE id = $1;`

	var portfolioStock dto_models.PortfolioStock
	err := db.QueryRow(query, portfolioStockId).Scan(&portfolioStock.Id, &portfolioStock.PortfolioId, &portfolioStock.StockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolioStock, nil
}

func getAllPortfolioStocksByPortfolioId(db *sql.DB, portfolioId int) ([]*dto_models.PortfolioStock, error) {
	query := `SELECT id, portfolio_id, stock_id FROM portfolio_stocks WHERE portfolio_id = $1;`

	rows, err := db.Query(query, portfolioId)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolioStocks []*dto_models.PortfolioStock
	var portfolioStock dto_models.PortfolioStock
	for rows.Next() {
		err = rows.Scan(&portfolioStock.Id, &portfolioStock.PortfolioId, &portfolioStock.StockId)
		if err != nil {
			return nil, err
		}

		portfolioStocks = append(portfolioStocks, &portfolioStock)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolioStocks, nil
}

func getAllPortfolioStocksByStocksId(db *sql.DB, stockId int) ([]*dto_models.PortfolioStock, error) {
	query := `SELECT id, portfolio_id, stock_id FROM portfolio_stocks WHERE stock_id = $1;`

	rows, err := db.Query(query, stockId)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolioStocks []*dto_models.PortfolioStock
	var portfolioStock dto_models.PortfolioStock
	for rows.Next() {
		err = rows.Scan(&portfolioStock.Id, &portfolioStock.PortfolioId, &portfolioStock.StockId)
		if err != nil {
			return nil, err
		}

		portfolioStocks = append(portfolioStocks, &portfolioStock)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolioStocks, nil
}

func deletePortfolioStock(db *sql.DB, portfolioStockId int) error {
	query := `DELETE FROM portfolio_stocks WHERE id = $1;`
	_, err := db.Exec(query, portfolioStockId)

	return err
}

/*
	PortfolioStockRelationship
*/

func createPortfolioStockRelationship(db *sql.DB, portfolioStockId int, quantity int) (*dto_models.PortfolioStockRelationship, error) {
	query := `INSERT INTO portfolio_stocks_relationship (portfolio_stocks_id, quantity) VALUES ($1, $2) RETURNING id;`

	var portfolioStockRelationshipId int
	err := db.QueryRow(query, portfolioStockId, quantity).Scan(&portfolioStockRelationshipId)
	if err != nil {
		return nil, err
	}

	return &dto_models.PortfolioStockRelationship{
		Id:               portfolioStockRelationshipId,
		PortfolioStockId: portfolioStockId,
		Quantity:         quantity,
	}, nil
}

func getPortfolioStockRelationship(db *sql.DB, relationshipId int) (*dto_models.PortfolioStockRelationship, error) {
	query := `SELECT id, portfolio_stocks_id, quantity FROM portfolio_stocks_relationship WHERE id = $1;`

	var relationship dto_models.PortfolioStockRelationship
	err := db.QueryRow(query, relationshipId).Scan(&relationship.Id, &relationship.PortfolioStockId, &relationship.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updatePortfolioStockRelationship(db *sql.DB, relationshipId int, newQuantity int) error {
	query := `UPDATE portfolio_stocks_relationship SET quantity = $1 WHERE id = $2;`
	_, err := db.Exec(query, newQuantity, relationshipId)

	return err
}

func deletePortfolioStockRelationship(db *sql.DB, relationshipId int) error {
	query := `DELETE FROM portfolio_stocks_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipId)

	return err
}

func getPortfolioStockRelationshipByPortfolioStockId(db *sql.DB, portfolioStockId int) (*dto_models.PortfolioStockRelationship, error) {
	query := `SELECT id, portfolio_stocks_id, quantity FROM portfolio_stocks_relationship WHERE portfolio_stocks_id = $1;`

	var relationship dto_models.PortfolioStockRelationship
	err := db.QueryRow(query, portfolioStockId).Scan(&relationship.Id, &relationship.PortfolioStockId, &relationship.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

/*
	Index
*/

func createIndex(db *sql.DB, name string) (*dto_models.Index, error) {
	query := `INSERT INTO indexes (name) VALUES ($1) RETURNING id;`

	var indexId int
	err := db.QueryRow(query, name).Scan(&indexId)
	if err != nil {
		return nil, err
	}

	return &dto_models.Index{
		Id:   indexId,
		Name: name,
	}, nil
}

func getIndex(db *sql.DB, indexId int) (*dto_models.Index, error) {
	query := `SELECT id, name FROM indexes WHERE id = $1;`

	var index dto_models.Index
	err := db.QueryRow(query, indexId).Scan(&index.Id, &index.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &index, nil
}

func updateIndex(db *sql.DB, indexId int, newName string) error {
	query := `UPDATE indexes SET name = $1 WHERE id = $2;`
	_, err := db.Exec(query, newName, indexId)

	return err
}

func deleteIndex(db *sql.DB, indexId int) error {
	query := `DELETE FROM indexes WHERE id = $1;`
	_, err := db.Exec(query, indexId)

	return err
}

func getIndexIdByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM indexes WHERE name = $1;`

	var indexId int
	err := db.QueryRow(query, name).Scan(&indexId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return indexId, nil
}

func getAllIndexStocksByIndexId(db *sql.DB, indexId int) ([]*dto_models.IndexStock, error) {
	query := `SELECT id, index_id, stock_id FROM index_stocks WHERE portfolio_id = $1;`

	rows, err := db.Query(query, indexId)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var indexStocks []*dto_models.IndexStock
	var indexStock dto_models.IndexStock
	for rows.Next() {
		err = rows.Scan(&indexStock.Id, &indexStock.IndexId, &indexStock.StockId)
		if err != nil {
			return nil, err
		}

		indexStocks = append(indexStocks, &indexStock)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return indexStocks, nil
}

func getAllIndexes(db *sql.DB) ([]*dto_models.Index, error) {
	query := `SELECT id, name FROM indexes;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var indexes []*dto_models.Index
	var index dto_models.Index
	for rows.Next() {
		err = rows.Scan(&index.Id, &index.Name)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, &index)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return indexes, nil
}

/*
	IndexStock
*/

func createIndexStock(db *sql.DB, indexId int, stockId int) (*dto_models.IndexStock, error) {
	query := `INSERT INTO index_stocks (index_id, stock_id) VALUES ($1, $2) RETURNING id;`

	var indexStockId int
	err := db.QueryRow(query, indexId, stockId).Scan(&indexStockId)
	if err != nil {
		return nil, err
	}

	return &dto_models.IndexStock{
		Id:      indexStockId,
		IndexId: indexId,
		StockId: stockId,
	}, nil
}

func getIndexStock(db *sql.DB, indexStockId int) (*dto_models.IndexStock, error) {
	query := `SELECT id, index_id, stock_id FROM index_stocks WHERE id = $1;`

	var indexStock dto_models.IndexStock
	err := db.QueryRow(query, indexStockId).Scan(&indexStock.Id, &indexStock.IndexId, &indexStock.StockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &indexStock, nil
}

func deleteIndexStock(db *sql.DB, indexStockId int) error {
	query := `DELETE FROM index_stocks WHERE id = $1;`
	_, err := db.Exec(query, indexStockId)

	return err
}

/*
	IndexStockRelationship
*/

func createIndexStockRelationship(db *sql.DB, indexStockId int, fraction float64) (*dto_models.IndexStockRelationship, error) {
	query := `INSERT INTO index_stocks_relationship (index_stocks_id, fraction) VALUES ($1, $2) RETURNING id;`

	var relationshipId int
	err := db.QueryRow(query, indexStockId, fraction).Scan(&relationshipId)
	if err != nil {
		return nil, err
	}

	return &dto_models.IndexStockRelationship{
		Id:           relationshipId,
		IndexStockId: indexStockId,
		Fraction:     fraction,
	}, nil
}

func getIndexStockRelationship(db *sql.DB, relationshipId int) (*dto_models.IndexStockRelationship, error) {
	query := `SELECT id, index_stocks_id, fraction FROM index_stocks_relationship WHERE id = $1;`

	var relationship dto_models.IndexStockRelationship
	err := db.QueryRow(query, relationshipId).Scan(&relationship.Id, &relationship.IndexStockId, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updateIndexStockRelationship(db *sql.DB, relationshipId int, newFraction float64) error {
	query := `UPDATE index_stocks_relationship SET fraction = $1 WHERE id = $2;`
	_, err := db.Exec(query, newFraction, relationshipId)

	return err
}

func deleteIndexStockRelationship(db *sql.DB, relationshipId int) error {
	query := `DELETE FROM index_stocks_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipId)

	return err
}

func getIndexStockRelationshipByIndexStockId(db *sql.DB, indexStockId int) (*dto_models.IndexStockRelationship, error) {
	query := `SELECT id, index_stocks_id, fraction FROM index_stocks_relationship WHERE index_stocks_id = $1;`

	var relationship dto_models.IndexStockRelationship
	err := db.QueryRow(query, indexStockId).Scan(&relationship.Id, &relationship.IndexStockId, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

/*
	ADDITIONAL FUNCTIONALITY
*/

func addStockToPortfolio(db *sql.DB, portfolioId, stockId, quantity int) error {
	portfolioStock, err := createPortfolioStock(db, portfolioId, stockId)
	if err != nil {
		return err
	}

	_, err = createPortfolioStockRelationship(db, portfolioStock.Id, quantity)

	return err
}

func addManyStocksToPortfolio(db *sql.DB, portfolioId int, stocksQuantityMap map[int]int) error {
	for stockId, quantity := range stocksQuantityMap {
		err := addStockToPortfolio(db, portfolioId, stockId, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertStocksNameQuantityMapToStocksIdQuantityMap(db *sql.DB, stocksQuantityMap map[string]int) (map[int]int, error) {
	stocksIdQuantityMap := make(map[int]int)
	for stockName, quantity := range stocksQuantityMap {
		stockId, err := getStockIdByName(db, stockName)
		if err != nil {
			return nil, err
		}

		stocksIdQuantityMap[stockId] = quantity
	}

	return stocksIdQuantityMap, nil
}

func addManyStocksToPortfolioByName(db *sql.DB, portfolioName string, stocksQuantityMap map[string]int) error {
	portfolioId, err := getPortfolioIdByName(db, portfolioName)
	if err != nil {
		return err
	}

	stocksIdQuantityMap := make(map[int]int)
	stocksIdQuantityMap, err = convertStocksNameQuantityMapToStocksIdQuantityMap(db, stocksQuantityMap)
	if err != nil {
		return err
	}

	err = addManyStocksToPortfolio(db, portfolioId, stocksIdQuantityMap)

	return err
}

func deleteAllStocksFromPortfolio(db *sql.DB, portfolioId int) error {
	query := `
        DELETE FROM portfolio_stocks_relationship
        WHERE portfolio_stocks_id IN (
            SELECT id FROM portfolio_stocks WHERE portfolio_id = $1
        );
    `
	_, err := db.Exec(query, portfolioId)
	if err != nil {
		return err
	}

	query = `DELETE FROM portfolio_stocks WHERE portfolio_id = $1;`
	_, err = db.Exec(query, portfolioId)

	return err
}

func deleteAllStocksFromPortfolioByName(db *sql.DB, portfolioName string) error {
	portfolioId, err := getPortfolioIdByName(db, portfolioName)
	if err != nil {
		return err
	}

	err = deleteAllStocksFromPortfolio(db, portfolioId)

	return err
}

func convertStocksQuantityMapToStocksIdQuantityMap(db *sql.DB, stocksQuantityMap map[*models.Stock]int) (map[int]int, error) {
	stocksNameQuantityMap := make(map[string]int)
	for stock, quantity := range stocksQuantityMap {
		stocksNameQuantityMap[stock.Name] = quantity
	}

	stocksIdQuantityMap, err := convertStocksNameQuantityMapToStocksIdQuantityMap(db, stocksNameQuantityMap)
	if err != nil {
		return nil, err
	}

	return stocksIdQuantityMap, nil
}

func addStockToIndex(db *sql.DB, indexId int, stockId int, fraction float64) error {
	indexStock, err := createIndexStock(db, indexId, stockId)
	if err != nil {
		return err
	}

	_, err = createIndexStockRelationship(db, indexStock.Id, fraction)

	return err
}

func addManyStocksToIndex(db *sql.DB, indexId int, stocksFractionMap map[int]float64) error {
	for stockId, fraction := range stocksFractionMap {
		err := addStockToIndex(db, indexId, stockId, fraction)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertStocksNameFractionMapToStocksIdFractionMap(db *sql.DB, stocksFractionMap map[string]float64) (map[int]float64, error) {
	stocksIdFractionMap := make(map[int]float64)
	for stockName, fraction := range stocksFractionMap {
		stockId, err := getStockIdByName(db, stockName)
		if err != nil {
			return nil, err
		}

		stocksIdFractionMap[stockId] = fraction
	}

	return stocksIdFractionMap, nil
}

func deleteAllStocksFromIndex(db *sql.DB, indexId int) error {
	query := `
        DELETE FROM index_stocks_relationship
        WHERE index_stocks_id IN (
            SELECT id FROM index_stocks WHERE index_id = $1
        );
    `
	_, err := db.Exec(query, indexId)
	if err != nil {
		return err
	}

	query = `DELETE FROM index_stocks WHERE index_id = $1;`
	_, err = db.Exec(query, indexId)

	return err
}

func convertStocksFractionMapToStocksIdFractionMap(db *sql.DB, stocksFractionMap map[*models.Stock]float64) (map[int]float64, error) {
	stocksNameFractionMap := make(map[string]float64)
	for stock, quantity := range stocksFractionMap {
		stocksNameFractionMap[stock.Name] = quantity
	}

	stocksIDFractionMap, err := convertStocksNameFractionMapToStocksIdFractionMap(db, stocksNameFractionMap)
	if err != nil {
		return nil, err
	}

	return stocksIDFractionMap, nil
}
