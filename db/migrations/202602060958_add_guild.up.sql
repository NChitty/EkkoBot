CREATE TABLE guilds(
  id bigserial PRIMARY KEY,
  discordId varchar(20),
  lastUpdated TIMESTAMPTZ
)
