-- name: CreatePortfolio :one
INSERT INTO portfolios (name) VALUES ($1) RETURNING *;

-- name: GetPortfolio :one
SELECT * FROM portfolios WHERE id = $1;

-- name: GetPortfolioByName :one
SELECT * FROM portfolios WHERE name = $1;

-- name: DeletePortfolio :exec
DELETE FROM portfolios WHERE id = $1;

-- name: GetAllPortfolios :many
SELECT * FROM portfolios;

-- name: UpsertPortfolioAsset :exec
INSERT INTO portfolio_assets (portfolio_id, asset_id, quantity) 
VALUES ($1, $2, $3)
ON CONFLICT (portfolio_id, asset_id) DO UPDATE SET quantity = EXCLUDED.quantity;

-- name: UpsertPortfolioAssetLot :exec
INSERT INTO portfolio_assets (portfolio_id, asset_id, lot_size) 
VALUES ($1, $2, $3)
ON CONFLICT (portfolio_id, asset_id) DO UPDATE SET lot_size = EXCLUDED.lot_size;

-- name: DeletePortfolioAsset :exec
DELETE FROM portfolio_assets WHERE portfolio_id = $1 AND asset_id = $2;

-- name: GetPortfolioAssets :many
SELECT a.id, a.name, a.price, pa.quantity, pa.lot_size
FROM portfolio_assets pa
JOIN assets a ON pa.asset_id = a.id
WHERE pa.portfolio_id = $1;
