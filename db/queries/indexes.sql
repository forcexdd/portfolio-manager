-- name: UpsertIndex :one
INSERT INTO indexes (name) VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING *;

-- name: GetAllIndexes :many
SELECT * FROM indexes;

-- name: UpsertIndexAsset :exec
INSERT INTO index_assets (index_id, asset_id, fraction)
VALUES ($1, $2, $3)
ON CONFLICT (index_id, asset_id) DO UPDATE SET fraction = EXCLUDED.fraction;

-- name: GetIndexAssets :many
SELECT a.id, a.name, a.price, ia.fraction
FROM index_assets ia
JOIN assets a ON ia.asset_id = a.id
WHERE ia.index_id = $1;
