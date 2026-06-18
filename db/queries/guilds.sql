-- name: GetGuildByDiscordId :one
SELECT * FROM guilds
WHERE discord_id = $1;

-- name: GetGuildById :one
SELECT * FROM guilds
WHERE id = $1;

-- name: CreateGuild :one
INSERT INTO guilds (discord_id, last_updated)
VALUES ($1, now())
ON CONFLICT (discord_id) DO UPDATE
  SET last_updated = NOW()
RETURNING *;

-- name: UpdateLastRanTime :one
UPDATE guilds
SET last_updated = now()
WHERE discord_id = $1
RETURNING *;
