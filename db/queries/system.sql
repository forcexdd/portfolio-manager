-- name: HasParsedDate :one
SELECT EXISTS(SELECT 1 FROM parse_logs WHERE date = $1);

-- name: MarkDateParsed :exec
INSERT INTO parse_logs (date) VALUES ($1) ON CONFLICT DO NOTHING;
