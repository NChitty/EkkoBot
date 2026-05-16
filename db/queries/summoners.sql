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

-- name: UpdateSummoner :one
WITH summoner_upsert AS (
  INSERT INTO summoners (name, tag_line, player_uuid)
  VALUES ($1, $2, $3)
  ON CONFLICT (player_uuid)
  DO NOTHING
  RETURNING id
),
summoner_id AS (
  SELECT id from summoner_upsert
  UNION
  SELECT id from summoners WHERE player_uuid = $3
),
inserted AS (
  INSERT INTO guild_summoners (
    summoner_id,
    guild_id,
    flex_games_played,
    flex_tier,
    flex_rank,
    flex_wins,
    flex_lp,
    solo_duo_games_played,
    solo_duo_tier,
    solo_duo_rank,
    solo_duo_wins,
    solo_duo_lp,
    last_updated
  )
  VALUES (
    (SELECT id FROM summoner_id),
    (SELECT id FROM guilds WHERE guilds.discord_id = $4),
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    NOW()
  )
  RETURNING *
)
SELECT * FROM inserted
JOIN summoners ON summoners.id = inserted.summoner_id
JOIN guilds ON guilds.id = inserted.guild_id;
