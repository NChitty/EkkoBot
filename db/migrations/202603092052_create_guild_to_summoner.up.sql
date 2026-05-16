CREATE TABLE guild_summoners(
  id bigserial PRIMARY KEY,
  summoner_id bigint REFERENCES summoners(id),
  guild_id bigint REFERENCES guilds(id),
  flex_games_played int,
  flex_tier text,
  flex_rank text,
  flex_wins int,
  flex_lp int,
  solo_duo_games_played int,
  solo_duo_tier text,
  solo_duo_rank text,
  solo_duo_wins int,
  solo_duo_lp int,
  last_updated TIMESTAMPTZ,
  UNIQUE (summoner_id, guild_id)
);
