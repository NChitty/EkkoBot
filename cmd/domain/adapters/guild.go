package adapters

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
	"github.com/NChitty/lol-discord-bot/cmd/domain/services"
	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type GuildRepository struct {
	queries *db.Queries
}

func NewGuildRepository(q *db.Queries) *GuildRepository {
	return &GuildRepository{q}
}

func (r *GuildRepository) CreateGuild(ctx context.Context, discordId string) (models.Guild, error) {
	if guild, err := r.GetGuild(ctx, discordId); err == nil {
		slog.DebugContext(ctx, "Guild exists", "DiscordID", discordId)
		return guild, nil
	} else if err.Error() == fmt.Sprintf(services.GUILD_NOT_FOUND_ERRORF, discordId) {
		if row, err := r.queries.CreateGuild(ctx, pgtype.Text{String: discordId, Valid: true}); err == nil {
			slog.DebugContext(ctx, "Created guild", "DiscordID", discordId)
			return fromRow(row), nil
		} else {
			return models.Guild{}, err
		}
	} else {
		return models.Guild{}, err
	}
}

func (r *GuildRepository) GetGuild(ctx context.Context, discordId string) (models.Guild, error) {
	text := pgtype.Text{String: discordId, Valid: true}
	if row, err := r.queries.GetGuildByDiscordId(ctx, text); err != nil && err.Error() == "no rows in result set" {
		slog.DebugContext(ctx, "Guild does not exist", "DiscordID", discordId)
		return models.Guild{}, fmt.Errorf(services.GUILD_NOT_FOUND_ERRORF, discordId)
	} else if err != nil {
		return models.Guild{}, err
	} else {
		return fromRow(row), nil
	}
}

func fromRow(row db.GuildRow) models.Guild {
	return models.Guild{
		ID:          row.ID,
		DiscordID:   row.DiscordID.String,
		LastUpdated: row.LastUpdated.Time,
	}
}
