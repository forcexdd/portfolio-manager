package repositories

import (
	"database/sql"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/dto_models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
)

type IndexRepository interface {
	Create(index *models.Index) error
	GetByName(name string) (*models.Index, error)
	Update(index *models.Index) error
	Delete(index *models.Index) error
	DeleteByName(name string) error
	GetAll() ([]*models.Index, error)
}

type PostgresIndexRepository struct {
	db *sql.DB
}

func NewIndexRepository(db *sql.DB) IndexRepository {
	return &PostgresIndexRepository{db: db}
}

func (p *PostgresIndexRepository) Create(index *models.Index) error {
	createdPortfolio, err := createIndex(p.db, index.Name)
	if err != nil {
		return err
	}
	if index.StocksFractionMap == nil {
		return nil
	}

	stocksIdQuantityMap := make(map[int]float64)
	stocksIdQuantityMap, err = convertStocksFractionMapToStocksIdFractionMap(p.db, index.StocksFractionMap)

	err = addManyStocksToIndex(p.db, createdPortfolio.Id, stocksIdQuantityMap)

	return err
}

func (p *PostgresIndexRepository) GetByName(name string) (*models.Index, error) {
	indexId, err := getIndexIdByName(p.db, name)
	if err != nil {
		return nil, err
	}
	if indexId == 0 {
		return nil, nil
	}

	var indexStocks []*dto_models.IndexStock
	indexStocks, err = getAllIndexStocksByIndexId(p.db, indexId)
	if err != nil {
		return nil, err
	}

	stocksFractionMap := make(map[*models.Stock]float64)
	var stock *dto_models.Stock
	var relationship *dto_models.IndexStockRelationship
	for _, indexStock := range indexStocks {
		stock, err = getStock(p.db, indexStock.StockId)
		if err != nil {
			return nil, err
		}
		if stock == nil {
			return nil, errors.New("stock not found")
		}

		relationship, err = getIndexStockRelationshipByIndexStockId(p.db, indexStock.Id)
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

		stocksFractionMap[&newStock] = relationship.Fraction
	}

	return &models.Index{
		Name:              name,
		StocksFractionMap: stocksFractionMap,
	}, nil
}

func (p *PostgresIndexRepository) Update(index *models.Index) error {
	indexId, err := getIndexIdByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexId == 0 {
		return errors.New("index not found")
	}

	err = deleteAllStocksFromIndex(p.db, indexId)
	if err != nil {
		return err
	}

	stocksIdFractionMap := make(map[int]float64)
	stocksIdFractionMap, err = convertStocksFractionMapToStocksIdFractionMap(p.db, index.StocksFractionMap)
	if err != nil {
		return err
	}

	err = addManyStocksToIndex(p.db, indexId, stocksIdFractionMap)

	return err
}

func (p *PostgresIndexRepository) Delete(index *models.Index) error {
	indexId, err := getIndexIdByName(p.db, index.Name)
	if err != nil {
		return err
	}
	if indexId == 0 {
		return errors.New("index not found")
	}

	err = deleteAllStocksFromIndex(p.db, indexId)
	if err != nil {
		return err
	}

	err = deleteIndex(p.db, indexId)

	return err
}

func (p *PostgresIndexRepository) DeleteByName(name string) error {
	return p.Delete(&models.Index{Name: name, StocksFractionMap: nil})
}

func (p *PostgresIndexRepository) GetAll() ([]*models.Index, error) {
	dtoIndexes, err := getAllIndexes(p.db)
	if err != nil {
		return nil, err
	}

	var indexes []*models.Index
	var index *models.Index
	for _, dtoIndex := range dtoIndexes {
		index, err = p.GetByName(dtoIndex.Name)
		if err != nil {
			return nil, err
		}

		newIndex := &models.Index{Name: index.Name,
			StocksFractionMap: index.StocksFractionMap,
		}

		indexes = append(indexes, newIndex)
	}

	return indexes, nil
}
