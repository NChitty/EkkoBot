package adapters

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
	"github.com/NChitty/lol-discord-bot/cmd/domain/services"
	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/jackc/pgx/v5/pgtype"
)

type SummonerRepository struct {
	queries *db.Queries
}

func NewSummonerRepository(q *db.Queries) *SummonerRepository {
	return &SummonerRepository{q}
}

func (s *SummonerRepository) GetSummoner(ctx context.Context, name string, tag string) (models.Summoner, error) {
	params := db.GetSummonerByNameAndTagParams{
		Name:    pgtype.Text{String: name, Valid: true},
		TagLine: pgtype.Text{String: tag, Valid: true},
	}
	if row, err := s.queries.GetSummonerByNameAndTag(ctx, params); err == nil {
		return models.Summoner{
			ID:         row.ID,
			Name:       row.Name.String,
			TagLine:    row.TagLine.String,
			PlayerUuid: row.PlayerUuid.String,
		}, nil
	} else if err.Error() == "no rows in result set" {
		return models.Summoner{}, fmt.Errorf(services.SUMMONER_NOT_FOUND_ERRORF, name, tag)
	} else {
		return models.Summoner{}, err
	}
}

func (s *SummonerRepository) SaveSummoner(ctx context.Context, stats models.SummonerStats) (models.SummonerStats, error) {
	params := db.UpdateSummonerParams{
		Name:       pgtype.Text{String: stats.Summoner.Name, Valid: true},
		TagLine:    pgtype.Text{String: stats.Summoner.TagLine, Valid: true},
		PlayerUuid: pgtype.Text{String: stats.Summoner.PlayerUuid, Valid: true},
		// Not a fan, assumes I know what the context this is being called in
		DiscordID:          pgtype.Text{String: ctx.Value(commands.CONTEXT_KEY).(commands.CommandContext).DiscordId},
		FlexGamesPlayed:    pgtype.Int4{Int32: int32(stats.FlexGamesPlayed), Valid: true},
		FlexTier:           pgtype.Text{String: stats.FlexTier, Valid: true},
		FlexRank:           pgtype.Text{String: stats.FlexRank, Valid: true},
		FlexWins:           pgtype.Int4{Int32: int32(stats.FlexWins), Valid: true},
		FlexLp:             pgtype.Int4{Int32: int32(stats.FlexLp), Valid: true},
		SoloDuoGamesPlayed: pgtype.Int4{Int32: int32(stats.SoloDuoGamesPlayed), Valid: true},
		SoloDuoTier:        pgtype.Text{String: stats.SoloDuoTier, Valid: true},
		SoloDuoRank:        pgtype.Text{String: stats.SoloDuoRank, Valid: true},
		SoloDuoWins:        pgtype.Int4{Int32: int32(stats.SoloDuoWins), Valid: true},
		SoloDuoLp:          pgtype.Int4{Int32: int32(stats.SoloDuoLp), Valid: true},
	}
	row, err := s.queries.UpdateSummoner(ctx, params)
	if err != nil {
		slog.Error("Failed to update summoner stats", "name", stats.Summoner.Name, "tagline", stats.Summoner.TagLine)
		return models.SummonerStats{}, err
	}

	return models.SummonerStats{
		Summoner: models.Summoner{
			ID:         row.ID_2,
			Name:       row.Name.String,
			TagLine:    row.TagLine.String,
			PlayerUuid: row.PlayerUuid.String,
		},
		FlexGamesPlayed:    int(row.FlexGamesPlayed.Int32),
		FlexTier:           row.FlexTier.String,
		FlexRank:           row.FlexRank.String,
		FlexWins:           int(row.FlexWins.Int32),
		FlexLp:             int(row.FlexLp.Int32),
		SoloDuoGamesPlayed: int(row.SoloDuoGamesPlayed.Int32),
		SoloDuoTier:        row.SoloDuoTier.String,
		SoloDuoRank:        row.SoloDuoRank.String,
		SoloDuoWins:        int(row.SoloDuoWins.Int32),
		SoloDuoLp:          int(row.SoloDuoLp.Int32),
	}, nil
}
