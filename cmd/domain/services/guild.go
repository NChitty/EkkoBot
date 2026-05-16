package services

import (
	"context"
	"fmt"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
)

type GuildRepository interface {
	GetGuild(ctx context.Context, discordId string) (models.Guild, error)
	CreateGuild(ctx context.Context, discordId string) (models.Guild, error)
}

type GuildService struct {
	guildRepository GuildRepository
}

func NewGuildService(r GuildRepository) *GuildService {
	return &GuildService{r}
}

const GUILD_NOT_FOUND_ERRORF string = "Could not find guild %s"

func (g *GuildService) GetGuild(ctx context.Context, discordId string) (models.Guild, error) {
	if guild, err := g.guildRepository.GetGuild(ctx, discordId); err != nil && err.Error() == fmt.Sprintf(GUILD_NOT_FOUND_ERRORF, discordId) {
		return g.guildRepository.CreateGuild(ctx, discordId)
	} else {
		return guild, err
	}
}
