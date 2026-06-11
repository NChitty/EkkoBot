ALTER TABLE guild_summoners
ADD CONSTRAINT guild_summoners_summoner_id_guild_id_key UNIQUE (summoner_id, guild_id);
