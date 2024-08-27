package repository

import (
	"database/sql"
	"errors"
	dtomodel "github.com/forcexdd/portfoliomanager/src/web/backend/database/model"
)

/*
	Portfolio
*/

func createPortfolio(db *sql.DB, name string) (*dtomodel.Portfolio, error) {
	query := `INSERT INTO portfolios (name) VALUES ($1) RETURNING id;`

	var portfolioID int
	err := db.QueryRow(query, name).Scan(&portfolioID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.Portfolio{
		ID:   portfolioID,
		Name: name,
	}, nil
}

func getPortfolio(db *sql.DB, portfolioID int) (*dtomodel.Portfolio, error) {
	query := `SELECT id, name FROM portfolios WHERE id = $1;`

	var portfolio dtomodel.Portfolio
	err := db.QueryRow(query, portfolioID).Scan(&portfolio.ID, &portfolio.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolio, nil
}

func updatePortfolio(db *sql.DB, portfolioID int, newName string) error {
	query := `UPDATE portfolios SET name = $1 WHERE id = $2;`
	_, err := db.Exec(query, newName, portfolioID)

	return err
}

func deletePortfolio(db *sql.DB, portfolioID int) error {
	query := `DELETE FROM portfolios WHERE id = $1;`
	_, err := db.Exec(query, portfolioID)

	return err
}

func getPortfolioIDByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM portfolios WHERE name = $1;`

	var portfolioID int
	err := db.QueryRow(query, name).Scan(&portfolioID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return portfolioID, nil
}

func getAllPortfolios(db *sql.DB) ([]*dtomodel.Portfolio, error) {
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

	var portfolios []*dtomodel.Portfolio
	var portfolio dtomodel.Portfolio
	for rows.Next() {
		err = rows.Scan(&portfolio.ID, &portfolio.Name)
		if err != nil {
			return nil, err
		}

		newPortfolio := &dtomodel.Portfolio{
			ID:   portfolio.ID,
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

func createAsset(db *sql.DB, name string, price float64) (*dtomodel.Asset, error) {
	query := `INSERT INTO assets (name, price) VALUES ($1, $2) RETURNING id;`

	var assetID int
	err := db.QueryRow(query, name, price).Scan(&assetID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.Asset{
		ID:    assetID,
		Name:  name,
		Price: price,
	}, nil
}

func getAsset(db *sql.DB, assetID int) (*dtomodel.Asset, error) {
	query := `SELECT id, name, price FROM assets WHERE id = $1;`

	var asset dtomodel.Asset
	err := db.QueryRow(query, assetID).Scan(&asset.ID, &asset.Name, &asset.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

func updateAsset(db *sql.DB, assetID int, newPrice float64) error {
	query := `UPDATE assets SET price = $1 WHERE id = $2;`
	_, err := db.Exec(query, newPrice, assetID)

	return err
}

func deleteAsset(db *sql.DB, assetID int) error {
	query := `DELETE FROM assets WHERE id = $1;`
	_, err := db.Exec(query, assetID)

	return err
}

func getAssetIDByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM assets WHERE name = $1;`

	var assetID int
	err := db.QueryRow(query, name).Scan(&assetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return assetID, nil
}

func getAllAssets(db *sql.DB) ([]*dtomodel.Asset, error) {
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

	var assets []*dtomodel.Asset
	var asset dtomodel.Asset
	for rows.Next() {
		err = rows.Scan(&asset.ID, &asset.Name, &asset.Price)
		if err != nil {
			return nil, err
		}

		newAsset := &dtomodel.Asset{
			ID:    asset.ID,
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

func deleteAssetFromConnectedTables(db *sql.DB, assetID int) error {
	query := `
        DELETE FROM portfolio_assets_relationship
        WHERE portfolio_assets_id IN (
            SELECT id FROM portfolio_assets WHERE asset_id = $1
        );
    `
	_, err := db.Exec(query, assetID)
	if err != nil {
		return err
	}

	query = `DELETE FROM portfolio_assets WHERE asset_id = $1;`
	_, err = db.Exec(query, assetID)
	if err != nil {
		return err
	}

	query = `
        DELETE FROM index_assets_relationship
        WHERE index_assets_id IN (
            SELECT id FROM index_assets WHERE asset_id = $1
        );
    `
	_, err = db.Exec(query, assetID)
	if err != nil {
		return err
	}

	query = `DELETE FROM index_assets WHERE asset_id = $1;`
	_, err = db.Exec(query, assetID)

	return err
}

/*
	PortfolioAsset
*/

func createPortfolioAsset(db *sql.DB, portfolioID int, assetID int) (*dtomodel.PortfolioAsset, error) {
	query := `INSERT INTO portfolio_assets (portfolio_id, asset_id) VALUES ($1, $2) RETURNING id;`

	var portfolioAssetID int
	err := db.QueryRow(query, portfolioID, assetID).Scan(&portfolioAssetID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.PortfolioAsset{
		ID:          portfolioAssetID,
		PortfolioID: portfolioID,
		AssetID:     assetID,
	}, nil
}

func getPortfolioAsset(db *sql.DB, portfolioAssetID int) (*dtomodel.PortfolioAsset, error) {
	query := `SELECT id, portfolio_id, asset_id FROM portfolio_assets WHERE id = $1;`

	var portfolioAsset dtomodel.PortfolioAsset
	err := db.QueryRow(query, portfolioAssetID).Scan(&portfolioAsset.ID, &portfolioAsset.PortfolioID, &portfolioAsset.AssetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &portfolioAsset, nil
}

func deletePortfolioAsset(db *sql.DB, portfolioAssetID int) error {
	query := `DELETE FROM portfolio_assets WHERE id = $1;`
	_, err := db.Exec(query, portfolioAssetID)

	return err
}

func getPortfolioAssetIDByPortfolioIdAndAssetID(db *sql.DB, portfolioID int, assetID int) (int, error) {
	query := `SELECT id FROM portfolio_assets WHERE portfolio_id = $1 AND asset_id = $2;`

	var portfolioAssetID int
	err := db.QueryRow(query, portfolioID, assetID).Scan(&portfolioAssetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return portfolioAssetID, nil
}

func getAllPortfolioAssetsByPortfolioID(db *sql.DB, portfolioID int) ([]*dtomodel.PortfolioAsset, error) {
	query := `SELECT id, portfolio_id, asset_id FROM portfolio_assets WHERE portfolio_id = $1;`
	rows, err := db.Query(query, portfolioID)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var portfolioAssets []*dtomodel.PortfolioAsset
	var portfolioAsset dtomodel.PortfolioAsset
	for rows.Next() {
		err = rows.Scan(&portfolioAsset.ID, &portfolioAsset.PortfolioID, &portfolioAsset.AssetID)
		if err != nil {
			return nil, err
		}

		newPortfolioAssets := &dtomodel.PortfolioAsset{
			ID:          portfolioAsset.ID,
			PortfolioID: portfolioAsset.PortfolioID,
			AssetID:     portfolioAsset.AssetID,
		}

		portfolioAssets = append(portfolioAssets, newPortfolioAssets)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return portfolioAssets, nil
}

/*
	PortfolioAssetRelationship
*/

func createPortfolioAssetRelationship(db *sql.DB, portfolioAssetID int, quantity int) (*dtomodel.PortfolioAssetRelationship, error) {
	query := `INSERT INTO portfolio_assets_relationship (portfolio_assets_id, quantity) VALUES ($1, $2) RETURNING id;`

	var portfolioAssetRelationshipID int
	err := db.QueryRow(query, portfolioAssetID, quantity).Scan(&portfolioAssetRelationshipID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.PortfolioAssetRelationship{
		ID:               portfolioAssetRelationshipID,
		PortfolioAssetID: portfolioAssetID,
		Quantity:         quantity,
	}, nil
}

func getPortfolioAssetRelationship(db *sql.DB, relationshipID int) (*dtomodel.PortfolioAssetRelationship, error) {
	query := `SELECT id, portfolio_assets_id, quantity FROM portfolio_assets_relationship WHERE id = $1;`

	var relationship dtomodel.PortfolioAssetRelationship
	err := db.QueryRow(query, relationshipID).Scan(&relationship.ID, &relationship.PortfolioAssetID, &relationship.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updatePortfolioAssetRelationship(db *sql.DB, relationshipID int, newQuantity int) error {
	query := `UPDATE portfolio_assets_relationship SET quantity = $1 WHERE id = $2;`
	_, err := db.Exec(query, newQuantity, relationshipID)

	return err
}

func deletePortfolioAssetRelationship(db *sql.DB, relationshipID int) error {
	query := `DELETE FROM portfolio_assets_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipID)

	return err
}

func getPortfolioAssetRelationshipByPortfolioAssetID(db *sql.DB, portfolioAssetID int) (*dtomodel.PortfolioAssetRelationship, error) {
	query := `SELECT id, portfolio_assets_id, quantity FROM portfolio_assets_relationship WHERE portfolio_assets_id = $1;`

	var relationship dtomodel.PortfolioAssetRelationship
	err := db.QueryRow(query, portfolioAssetID).Scan(&relationship.ID, &relationship.PortfolioAssetID, &relationship.Quantity)
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

func createIndex(db *sql.DB, name string) (*dtomodel.Index, error) {
	query := `INSERT INTO indexes (name) VALUES ($1) RETURNING id;`

	var indexID int
	err := db.QueryRow(query, name).Scan(&indexID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.Index{
		ID:   indexID,
		Name: name,
	}, nil
}

func getIndex(db *sql.DB, indexID int) (*dtomodel.Index, error) {
	query := `SELECT id, name FROM indexes WHERE id = $1;`

	var index dtomodel.Index
	err := db.QueryRow(query, indexID).Scan(&index.ID, &index.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &index, nil
}

func updateIndex(db *sql.DB, indexID int, newName string) error {
	query := `UPDATE indexes SET name = $1 WHERE id = $2;`
	_, err := db.Exec(query, newName, indexID)

	return err
}

func deleteIndex(db *sql.DB, indexID int) error {
	query := `DELETE FROM indexes WHERE id = $1;`
	_, err := db.Exec(query, indexID)

	return err
}

func getIndexIDByName(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM indexes WHERE name = $1;`

	var indexID int
	err := db.QueryRow(query, name).Scan(&indexID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return indexID, nil
}

func getAllIndexes(db *sql.DB) ([]*dtomodel.Index, error) {
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

	var indexes []*dtomodel.Index
	var index dtomodel.Index
	for rows.Next() {
		err = rows.Scan(&index.ID, &index.Name)
		if err != nil {
			return nil, err
		}

		newIndex := &dtomodel.Index{
			ID:   index.ID,
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

func createIndexAsset(db *sql.DB, indexID int, assetID int) (*dtomodel.IndexAsset, error) {
	query := `INSERT INTO index_assets (index_id, asset_id) VALUES ($1, $2) RETURNING id;`

	var indexAssetID int
	err := db.QueryRow(query, indexID, assetID).Scan(&indexAssetID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.IndexAsset{
		ID:      indexAssetID,
		IndexID: indexID,
		AssetID: assetID,
	}, nil
}

func getIndexAsset(db *sql.DB, indexAssetID int) (*dtomodel.IndexAsset, error) {
	query := `SELECT id, index_id, asset_id FROM index_assets WHERE id = $1;`

	var indexAsset dtomodel.IndexAsset
	err := db.QueryRow(query, indexAssetID).Scan(&indexAsset.ID, &indexAsset.IndexID, &indexAsset.AssetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &indexAsset, nil
}

func deleteIndexAsset(db *sql.DB, indexAssetID int) error {
	query := `DELETE FROM index_assets WHERE id = $1;`
	_, err := db.Exec(query, indexAssetID)

	return err
}

func getIndexAssetIDByIndexIdAndAssetID(db *sql.DB, indexID int, assetID int) (int, error) {
	query := `SELECT id FROM index_assets WHERE index_id = $1 AND asset_id = $2;`

	var indexAssetID int
	err := db.QueryRow(query, indexID, assetID).Scan(&indexAssetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return indexAssetID, nil
}

func getAllIndexAssetsByIndexID(db *sql.DB, indexID int) ([]*dtomodel.IndexAsset, error) {
	query := `SELECT id, index_id, asset_id FROM index_assets WHERE index_id = $1;`
	rows, err := db.Query(query, indexID)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var indexAssets []*dtomodel.IndexAsset
	var indexAsset dtomodel.IndexAsset
	for rows.Next() {
		err = rows.Scan(&indexAsset.ID, &indexAsset.IndexID, &indexAsset.AssetID)
		if err != nil {
			return nil, err
		}

		newIndexAssets := &dtomodel.IndexAsset{
			ID:      indexAsset.ID,
			IndexID: indexAsset.IndexID,
			AssetID: indexAsset.AssetID,
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

func createIndexAssetRelationship(db *sql.DB, indexAssetID int, fraction float64) (*dtomodel.IndexAssetRelationship, error) {
	query := `INSERT INTO index_assets_relationship (index_assets_id, fraction) VALUES ($1, $2) RETURNING id;`

	var relationshipID int
	err := db.QueryRow(query, indexAssetID, fraction).Scan(&relationshipID)
	if err != nil {
		return nil, err
	}

	return &dtomodel.IndexAssetRelationship{
		ID:           relationshipID,
		IndexAssetID: indexAssetID,
		Fraction:     fraction,
	}, nil
}

func getIndexAssetRelationship(db *sql.DB, relationshipID int) (*dtomodel.IndexAssetRelationship, error) {
	query := `SELECT id, index_assets_id, fraction FROM index_assets_relationship WHERE id = $1;`

	var relationship dtomodel.IndexAssetRelationship
	err := db.QueryRow(query, relationshipID).Scan(&relationship.ID, &relationship.IndexAssetID, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}

func updateIndexAssetRelationship(db *sql.DB, relationshipID int, newFraction float64) error {
	query := `UPDATE index_assets_relationship SET fraction = $1 WHERE id = $2;`
	_, err := db.Exec(query, newFraction, relationshipID)

	return err
}

func deleteIndexAssetRelationship(db *sql.DB, relationshipID int) error {
	query := `DELETE FROM index_assets_relationship WHERE id = $1;`
	_, err := db.Exec(query, relationshipID)

	return err
}

func getIndexAssetRelationshipByIndexAssetID(db *sql.DB, indexAssetID int) (*dtomodel.IndexAssetRelationship, error) {
	query := `SELECT id, index_assets_id, fraction FROM index_assets_relationship WHERE index_assets_id = $1;`

	var relationship dtomodel.IndexAssetRelationship
	err := db.QueryRow(query, indexAssetID).Scan(&relationship.ID, &relationship.IndexAssetID, &relationship.Fraction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &relationship, nil
}
