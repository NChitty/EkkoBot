# Features
- [ ] Add summoner to list to update
- [ ] scheduled messages to discord webhook
    - [ ] Post player W-L
    - [ ] Post player LP
- [ ] batched processing

# Getting Started
1. Clone the repository
2. Build the docker image
3. Grab the following secrets:
    - discord token -> DISCORD_TOKEN_FILE
    - riot token -> RIOT_API_TOKEN_FILE
    - PG secret -> generate using `gpg --gen-random --armor 1 24`
    - PG config secret -> TODO
4. Run the docker compose

# WIP
- [x] Docker + postgres
    - [x] docker secrets (discord token + postgres secrets)
- [x] Migrations
- [x] Capable of registering commands
- [x] Postgres connection (sqlc)
- [ ] Riot API integration
- [ ] Setup command
