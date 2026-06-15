package adapters

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
	"github.com/NChitty/lol-discord-bot/cmd/domain/services"
	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SummonerRepository struct {
	queries db.QuerierTx
	conn    *pgx.Conn
}

func NewSummonerRepository(q db.QuerierTx, conn *pgx.Conn) *SummonerRepository {
	return &SummonerRepository{q, conn}
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
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return models.SummonerStats{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTransaction(tx);
	guild, err := qtx.CreateGuild(ctx, pgtype.Text{String: ctx.Value(commands.CONTEXT_KEY).(commands.CommandContext).DiscordId, Valid: true})
	if err != nil {
		return models.SummonerStats{}, err
	}
	createParams := db.CreateSummonerParams{
		Name:       pgtype.Text{String: stats.Summoner.Name, Valid: true},
		TagLine:    pgtype.Text{String: stats.Summoner.TagLine, Valid: true},
		PlayerUuid: pgtype.Text{String: stats.Summoner.PlayerUuid, Valid: true},
	}
	summoner, err := qtx.CreateSummoner(ctx, createParams)
	if err != nil {
	  return models.SummonerStats{}, err
	}

	params := db.UpdateSummonerParams{
		SummonerID:         pgtype.Int8{Int64: summoner.ID, Valid: true},
		GuildID:            pgtype.Int8{Int64: guild.ID, Valid: true},
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
	err = tx.Commit(ctx)
	if err != nil {
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
