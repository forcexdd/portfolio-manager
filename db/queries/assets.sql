-- name: UpsertAsset :one
INSERT INTO assets (name, price) VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE SET price = EXCLUDED.price RETURNING *;

-- name: GetAssetByName :one
SELECT * FROM assets WHERE name = $1;

-- name: GetAllAssets :many
SELECT * FROM assets;
