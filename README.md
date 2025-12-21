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
    - PG config secret -> the file has the format `db:5432:ekkobot:[password]`
4. Run the docker compose

# WIP
- [ ] Clean up messaging
    - [ ] Add embed to "track" command for W-L and current LP
- [ ] Schedule for checking LP gains/losses
    - [ ] Add guild column for message update
    - [ ] Add column for last update
    - [ ] Add table for most recent queue stats
    - [ ] Display latest game stats
- [ ] Deploy to docker swarm
