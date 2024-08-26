package repository

import (
	"database/sql"
	"errors"
	dtomodels "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
)

/*
	Portfolio
*/

func createPortfolio(db *sql.DB, name string) (*dtomodels.Portfolio, error) {
	query := `INSERT INTO portfolios (name) VALUES ($1) RETURNING id;`

	var portfolioId int
	err := db.QueryRow(query, name).Scan(&portfolioId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.Portfolio{
		Id:   portfolioId,
		Name: name,
	}, nil
}

func getPortfolio(db *sql.DB, portfolioId int) (*dtomodels.Portfolio, error) {
	query := `SELECT id, name FROM portfolios WHERE id = $1;`

	var portfolio dtomodels.Portfolio
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

func getAllPortfolios(db *sql.DB) ([]*dtomodels.Portfolio, error) {
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

	var portfolios []*dtomodels.Portfolio
	var portfolio dtomodels.Portfolio
	for rows.Next() {
		err = rows.Scan(&portfolio.Id, &portfolio.Name)
		if err != nil {
			return nil, err
		}

		newPortfolio := &dtomodels.Portfolio{
			Id:   portfolio.Id,
			Name: portfolio.Name,
		}

		portfolios = append(portfolios, newPortfolio)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolios, nil
}

/*
	Asset
*/

func createAsset(db *sql.DB, name string, price float64) (*dtomodels.Asset, error) {
	query := `INSERT INTO assets (name, price) VALUES ($1, $2) RETURNING id;`

	var assetId int
	err := db.QueryRow(query, name, price).Scan(&assetId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.Asset{
		Id:    assetId,
		Name:  name,
		Price: price,
	}, nil
}

func getAsset(db *sql.DB, assetId int) (*dtomodels.Asset, error) {
	query := `SELECT id, name, price FROM assets WHERE id = $1;`

	var asset dtomodels.Asset
	err := db.QueryRow(query, assetId).Scan(&asset.Id, &asset.Name, &asset.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

func updateAsset(db *sql.DB, assetId int, newPrice float64) error {
	query := `UPDATE assets SET price = $1 WHERE id = $2;`
	_, err := db.Exec(query, newPrice, assetId)

	return err
}

func deleteAsset(db *sql.DB, assetId int) error {
	query := `DELETE FROM assets WHERE id = $1;`
	_, err := db.Exec(query, assetId)

	return err
}

func getAssetIdByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM assets WHERE name = $1;`

	var assetId int
	err := db.QueryRow(query, name).Scan(&assetId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return assetId, nil
}

func getAllAssets(db *sql.DB) ([]*dtomodels.Asset, error) {
	query := `SELECT id, name, price FROM assets;`
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

	var assets []*dtomodels.Asset
	var asset dtomodels.Asset
	for rows.Next() {
		err = rows.Scan(&asset.Id, &asset.Name, &asset.Price)
		if err != nil {
			return nil, err
		}

		newAsset := &dtomodels.Asset{
			Id:    asset.Id,
			Name:  asset.Name,
			Price: asset.Price,
		}

		assets = append(assets, newAsset)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func deleteAssetFromConnectedTables(db *sql.DB, assetId int) error {
	query := `
        DELETE FROM portfolio_assets_relationship
        WHERE portfolio_assets_id IN (
            SELECT id FROM portfolio_assets WHERE asset_id = $1
        );
    `
	_, err := db.Exec(query, assetId)
	if err != nil {
		return err
	}

	query = `DELETE FROM portfolio_assets WHERE asset_id = $1;`
	_, err = db.Exec(query, assetId)
	if err != nil {
		return err
	}

	query = `
        DELETE FROM index_assets_relationship
        WHERE index_assets_id IN (
            SELECT id FROM index_assets WHERE asset_id = $1
        );
    `
	_, err = db.Exec(query, assetId)
	if err != nil {
		return err
	}

	query = `DELETE FROM index_assets WHERE asset_id = $1;`
	_, err = db.Exec(query, assetId)

	return err
}

/*
	PortfolioAsset
*/

func createPortfolioAsset(db *sql.DB, portfolioId int, assetId int) (*dtomodels.PortfolioAsset, error) {
	query := `INSERT INTO portfolio_assets (portfolio_id, asset_id) VALUES ($1, $2) RETURNING id;`

	var portfolioAssetId int
	err := db.QueryRow(query, portfolioId, assetId).Scan(&portfolioAssetId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.PortfolioAsset{
		Id:          portfolioAssetId,
		PortfolioId: portfolioId,
		AssetId:     assetId,
	}, nil
}

func getPortfolioAsset(db *sql.DB, portfolioAssetId int) (*dtomodels.PortfolioAsset, error) {
	query := `SELECT id, portfolio_id, asset_id FROM portfolio_assets WHERE id = $1;`

	var portfolioAsset dtomodels.PortfolioAsset
	err := db.QueryRow(query, portfolioAssetId).Scan(&portfolioAsset.Id, &portfolioAsset.PortfolioId, &portfolioAsset.AssetId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolioAsset, nil
}

func getAllPortfolioAssetsByPortfolioId(db *sql.DB, portfolioId int) ([]*dtomodels.PortfolioAsset, error) {
	query := `SELECT id, portfolio_id, asset_id FROM portfolio_assets WHERE portfolio_id = $1;`
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

	var portfolioAssets []*dtomodels.PortfolioAsset
	var portfolioAsset dtomodels.PortfolioAsset
	for rows.Next() {
		err = rows.Scan(&portfolioAsset.Id, &portfolioAsset.PortfolioId, &portfolioAsset.AssetId)
		if err != nil {
			return nil, err
		}

		newPortfolioAssets := &dtomodels.PortfolioAsset{
			Id:          portfolioAsset.Id,
			PortfolioId: portfolioAsset.PortfolioId,
			AssetId:     portfolioAsset.AssetId,
		}

		portfolioAssets = append(portfolioAssets, newPortfolioAssets)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolioAssets, nil
}

func getAllPortfolioAssetsByAssetId(db *sql.DB, assetId int) ([]*dtomodels.PortfolioAsset, error) {
	query := `SELECT id, portfolio_id, asset_id FROM portfolio_assets WHERE asset_id = $1;`
	rows, err := db.Query(query, assetId)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolioAssets []*dtomodels.PortfolioAsset
	var portfolioAsset dtomodels.PortfolioAsset
	for rows.Next() {
		err = rows.Scan(&portfolioAsset.Id, &portfolioAsset.PortfolioId, &portfolioAsset.AssetId)
		if err != nil {
			return nil, err
		}

		newPortfolioAssets := &dtomodels.PortfolioAsset{
			Id:          portfolioAsset.Id,
			PortfolioId: portfolioAsset.PortfolioId,
			AssetId:     portfolioAsset.AssetId,
		}

		portfolioAssets = append(portfolioAssets, newPortfolioAssets)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolioAssets, nil
}

func deletePortfolioAsset(db *sql.DB, portfolioAssetId int) error {
	query := `DELETE FROM portfolio_assets WHERE id = $1;`
	_, err := db.Exec(query, portfolioAssetId)

	return err
}

/*
	PortfolioAssetRelationship
*/

func createPortfolioAssetRelationship(db *sql.DB, portfolioAssetId int, quantity int) (*dtomodels.PortfolioAssetRelationship, error) {
	query := `INSERT INTO portfolio_assets_relationship (portfolio_assets_id, quantity) VALUES ($1, $2) RETURNING id;`

	var portfolioAssetRelationshipId int
	err := db.QueryRow(query, portfolioAssetId, quantity).Scan(&portfolioAssetRelationshipId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.PortfolioAssetRelationship{
		Id:               portfolioAssetRelationshipId,
		PortfolioAssetId: portfolioAssetId,
		Quantity:         quantity,
	}, nil
}

func getPortfolioAssetRelationship(db *sql.DB, relationshipId int) (*dtomodels.PortfolioAssetRelationship, error) {
	query := `SELECT id, portfolio_assets_id, quantity FROM portfolio_assets_relationship WHERE id = $1;`

	var relationship dtomodels.PortfolioAssetRelationship
	err := db.QueryRow(query, relationshipId).Scan(&relationship.Id, &relationship.PortfolioAssetId, &relationship.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updatePortfolioAssetRelationship(db *sql.DB, relationshipId int, newQuantity int) error {
	query := `UPDATE portfolio_assets_relationship SET quantity = $1 WHERE id = $2;`
	_, err := db.Exec(query, newQuantity, relationshipId)

	return err
}

func deletePortfolioAssetRelationship(db *sql.DB, relationshipId int) error {
	query := `DELETE FROM portfolio_assets_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipId)

	return err
}

func getPortfolioAssetRelationshipByPortfolioAssetId(db *sql.DB, portfolioAssetId int) (*dtomodels.PortfolioAssetRelationship, error) {
	query := `SELECT id, portfolio_assets_id, quantity FROM portfolio_assets_relationship WHERE portfolio_assets_id = $1;`

	var relationship dtomodels.PortfolioAssetRelationship
	err := db.QueryRow(query, portfolioAssetId).Scan(&relationship.Id, &relationship.PortfolioAssetId, &relationship.Quantity)
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

func createIndex(db *sql.DB, name string) (*dtomodels.Index, error) {
	query := `INSERT INTO indexes (name) VALUES ($1) RETURNING id;`

	var indexId int
	err := db.QueryRow(query, name).Scan(&indexId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.Index{
		Id:   indexId,
		Name: name,
	}, nil
}

func getIndex(db *sql.DB, indexId int) (*dtomodels.Index, error) {
	query := `SELECT id, name FROM indexes WHERE id = $1;`

	var index dtomodels.Index
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

func getAllIndexes(db *sql.DB) ([]*dtomodels.Index, error) {
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

	var indexes []*dtomodels.Index
	var index dtomodels.Index
	for rows.Next() {
		err = rows.Scan(&index.Id, &index.Name)
		if err != nil {
			return nil, err
		}

		newIndex := &dtomodels.Index{
			Id:   index.Id,
			Name: index.Name,
		}

		indexes = append(indexes, newIndex)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return indexes, nil
}

/*
	IndexAsset
*/

func createIndexAsset(db *sql.DB, indexId int, assetId int) (*dtomodels.IndexAsset, error) {
	query := `INSERT INTO index_assets (index_id, asset_id) VALUES ($1, $2) RETURNING id;`

	var indexAssetId int
	err := db.QueryRow(query, indexId, assetId).Scan(&indexAssetId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.IndexAsset{
		Id:      indexAssetId,
		IndexId: indexId,
		AssetId: assetId,
	}, nil
}

func getIndexAsset(db *sql.DB, indexAssetId int) (*dtomodels.IndexAsset, error) {
	query := `SELECT id, index_id, asset_id FROM index_assets WHERE id = $1;`

	var indexAsset dtomodels.IndexAsset
	err := db.QueryRow(query, indexAssetId).Scan(&indexAsset.Id, &indexAsset.IndexId, &indexAsset.AssetId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &indexAsset, nil
}

func deleteIndexAsset(db *sql.DB, indexAssetId int) error {
	query := `DELETE FROM index_assets WHERE id = $1;`
	_, err := db.Exec(query, indexAssetId)

	return err
}

func getAllIndexAssetsByIndexId(db *sql.DB, indexId int) ([]*dtomodels.IndexAsset, error) {
	query := `SELECT id, index_id, asset_id FROM index_assets WHERE index_id = $1;`
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

	var indexAssets []*dtomodels.IndexAsset
	var indexAsset dtomodels.IndexAsset
	for rows.Next() {
		err = rows.Scan(&indexAsset.Id, &indexAsset.IndexId, &indexAsset.AssetId)
		if err != nil {
			return nil, err
		}

		newIndexAssets := &dtomodels.IndexAsset{
			Id:      indexAsset.Id,
			IndexId: indexAsset.IndexId,
			AssetId: indexAsset.AssetId,
		}

		indexAssets = append(indexAssets, newIndexAssets)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return indexAssets, nil
}

/*
	IndexAssetRelationship
*/

func createIndexAssetRelationship(db *sql.DB, indexAssetId int, fraction float64) (*dtomodels.IndexAssetRelationship, error) {
	query := `INSERT INTO index_assets_relationship (index_assets_id, fraction) VALUES ($1, $2) RETURNING id;`

	var relationshipId int
	err := db.QueryRow(query, indexAssetId, fraction).Scan(&relationshipId)
	if err != nil {
		return nil, err
	}

	return &dtomodels.IndexAssetRelationship{
		Id:           relationshipId,
		IndexAssetId: indexAssetId,
		Fraction:     fraction,
	}, nil
}

func getIndexAssetRelationship(db *sql.DB, relationshipId int) (*dtomodels.IndexAssetRelationship, error) {
	query := `SELECT id, index_assets_id, fraction FROM index_assets_relationship WHERE id = $1;`

	var relationship dtomodels.IndexAssetRelationship
	err := db.QueryRow(query, relationshipId).Scan(&relationship.Id, &relationship.IndexAssetId, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updateIndexAssetRelationship(db *sql.DB, relationshipId int, newFraction float64) error {
	query := `UPDATE index_assets_relationship SET fraction = $1 WHERE id = $2;`
	_, err := db.Exec(query, newFraction, relationshipId)

	return err
}

func deleteIndexAssetRelationship(db *sql.DB, relationshipId int) error {
	query := `DELETE FROM index_assets_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipId)

	return err
}

func getIndexAssetRelationshipByIndexAssetId(db *sql.DB, indexAssetId int) (*dtomodels.IndexAssetRelationship, error) {
	query := `SELECT id, index_assets_id, fraction FROM index_assets_relationship WHERE index_assets_id = $1;`

	var relationship dtomodels.IndexAssetRelationship
	err := db.QueryRow(query, indexAssetId).Scan(&relationship.Id, &relationship.IndexAssetId, &relationship.Fraction)
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

func addAssetToPortfolio(db *sql.DB, portfolioId, assetId, quantity int) error {
	portfolioAsset, err := createPortfolioAsset(db, portfolioId, assetId)
	if err != nil {
		return err
	}

	_, err = createPortfolioAssetRelationship(db, portfolioAsset.Id, quantity)

	return err
}

func addManyAssetsToPortfolio(db *sql.DB, portfolioId int, assetsQuantityMap map[int]int) error {
	for assetId, quantity := range assetsQuantityMap {
		err := addAssetToPortfolio(db, portfolioId, assetId, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertAssetsNameQuantityMapToAssetsIdQuantityMap(db *sql.DB, assetsQuantityMap map[string]int) (map[int]int, error) {
	assetsIdQuantityMap := make(map[int]int)
	for assetName, quantity := range assetsQuantityMap {
		assetId, err := getAssetIdByName(db, assetName)
		if err != nil {
			return nil, err
		}

		assetsIdQuantityMap[assetId] = quantity
	}

	return assetsIdQuantityMap, nil
}

func addManyAssetsToPortfolioByName(db *sql.DB, portfolioName string, assetsQuantityMap map[string]int) error {
	portfolioId, err := getPortfolioIdByName(db, portfolioName)
	if err != nil {
		return err
	}

	assetsIdQuantityMap := make(map[int]int)
	assetsIdQuantityMap, err = convertAssetsNameQuantityMapToAssetsIdQuantityMap(db, assetsQuantityMap)
	if err != nil {
		return err
	}

	err = addManyAssetsToPortfolio(db, portfolioId, assetsIdQuantityMap)

	return err
}

func deleteAllAssetsFromPortfolio(db *sql.DB, portfolioId int) error {
	query := `
        DELETE FROM portfolio_assets_relationship
        WHERE portfolio_assets_id IN (
            SELECT id FROM portfolio_assets WHERE portfolio_id = $1
        );
    `
	_, err := db.Exec(query, portfolioId)
	if err != nil {
		return err
	}

	query = `DELETE FROM portfolio_assets WHERE portfolio_id = $1;`
	_, err = db.Exec(query, portfolioId)

	return err
}

func deleteAllAssetsFromPortfolioByName(db *sql.DB, portfolioName string) error {
	portfolioId, err := getPortfolioIdByName(db, portfolioName)
	if err != nil {
		return err
	}

	err = deleteAllAssetsFromPortfolio(db, portfolioId)

	return err
}

func convertAssetsQuantityMapToAssetsIdQuantityMap(db *sql.DB, assetsQuantityMap map[*model.Asset]int) (map[int]int, error) {
	assetsNameQuantityMap := make(map[string]int)
	for asset, quantity := range assetsQuantityMap {
		assetsNameQuantityMap[asset.Name] = quantity
	}

	assetsIdQuantityMap, err := convertAssetsNameQuantityMapToAssetsIdQuantityMap(db, assetsNameQuantityMap)
	if err != nil {
		return nil, err
	}

	return assetsIdQuantityMap, nil
}

func addAssetToIndex(db *sql.DB, indexId int, assetId int, fraction float64) error {
	indexAsset, err := createIndexAsset(db, indexId, assetId)
	if err != nil {
		return err
	}

	_, err = createIndexAssetRelationship(db, indexAsset.Id, fraction)

	return err
}

func addManyAssetsToIndex(db *sql.DB, indexId int, assetsFractionMap map[int]float64) error {
	for assetId, fraction := range assetsFractionMap {
		err := addAssetToIndex(db, indexId, assetId, fraction)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertAssetsNameFractionMapToAssetsIdFractionMap(db *sql.DB, assetsFractionMap map[string]float64) (map[int]float64, error) {
	assetsIdFractionMap := make(map[int]float64)
	for assetName, fraction := range assetsFractionMap {
		assetId, err := getAssetIdByName(db, assetName)
		if err != nil {
			return nil, err
		}

		assetsIdFractionMap[assetId] = fraction
	}

	return assetsIdFractionMap, nil
}

func deleteAllAssetsFromIndex(db *sql.DB, indexId int) error {
	query := `
        DELETE FROM index_assets_relationship
        WHERE index_assets_id IN (
            SELECT id FROM index_assets WHERE index_id = $1
        );
    `
	_, err := db.Exec(query, indexId)
	if err != nil {
		return err
	}

	query = `DELETE FROM index_assets WHERE index_id = $1;`
	_, err = db.Exec(query, indexId)

	return err
}

func convertAssetsFractionMapToAssetsIdFractionMap(db *sql.DB, assetsFractionMap map[*model.Asset]float64) (map[int]float64, error) {
	assetsNameFractionMap := make(map[string]float64)
	for asset, quantity := range assetsFractionMap {
		assetsNameFractionMap[asset.Name] = quantity
	}

	assetsIDFractionMap, err := convertAssetsNameFractionMapToAssetsIdFractionMap(db, assetsNameFractionMap)
	if err != nil {
		return nil, err
	}

	return assetsIDFractionMap, nil
}
