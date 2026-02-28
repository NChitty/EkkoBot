-- name: GetSummoner :one
SELECT * FROM summoners
WHERE id = $1;

-- name: CreateSummoner :one
INSERT INTO summoners (name, tag_line, player_uuid)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSummonerByNameAndTag :one
SELECT * FROM summoners
WHERE name = $1 AND tag_line = $2;
