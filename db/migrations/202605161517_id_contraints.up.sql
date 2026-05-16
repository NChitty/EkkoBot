ALTER TABLE summoners
ADD CONSTRAINT summoners_player_uuid_key UNIQUE (player_uuid);
ALTER TABLE guilds
ADD CONSTRAINT guild_discord_id_key UNIQUE (discord_id);
