-- name: GetSummoner :one
SELECT * FROM summoners
WHERE id = $1;

-- name: CreateSummoner :one
INSERT INTO summoners (name, tagline, playerUuid)
VALUES ($1, $2, $3)
RETURNING *;
