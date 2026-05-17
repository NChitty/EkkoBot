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
    (SELECT id FROM summoners WHERE summoners.player_uuid = $1),
    (SELECT id FROM guilds WHERE guilds.id = $2),
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
JOIN summoners ON summoners.id = inserted.summoner_id;
