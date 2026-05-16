package commands

import (
	"context"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
)

type SummonerServicer interface {
	GetSummonerStats(ctx context.Context, name string, tag string) (models.SummonerStats, error)
}

type GuildServicer interface {
	GetGuild(ctx context.Context, discordId string) (models.Guild, error)
}
