package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/NChitty/lol-discord-bot/cmd/ports/riot"
)

type HttpRiotAdapter struct {
	riotClient riot.RiotClientInterface
}

func NewHttpRiotAdapter(riotClient riot.RiotClientInterface) *HttpRiotAdapter {
	return &HttpRiotAdapter{riotClient}
}

func (a *HttpRiotAdapter) GetSummoner(ctx context.Context, name string, tag string) (models.Summoner, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	params := riot.AccountByRiotIdRequestParams{
		Name:    name,
		Tagline: tag,
	}
	resp, err := a.riotClient.GetAccountByRiotId(reqCtx, params)
	if err != nil {
		slog.ErrorContext(reqCtx, "Failed to get riot account", "name", name, "tag", tag, "error", err)
		return models.Summoner{}, err
	}

	return models.Summoner{
		ID:         -1,
		Name:       resp.Name,
		TagLine:    resp.Tagline,
		PlayerUuid: resp.PlayerUuid,
	}, nil
}

func (a *HttpRiotAdapter) GetRankedStats(ctx context.Context, summoner models.Summoner) (models.SummonerStats, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	params := riot.QueueEntriesByPlayerUuidParams{PlayerUuid: summoner.PlayerUuid}
	resp, err := a.riotClient.GetQueueEntriesByPlayerUuid(reqCtx, params)
	contextValue := reqCtx.Value(commands.CONTEXT_KEY).(commands.CommandContext)
	if err != nil {
		slog.Error("Failed to get ranked queues", "playerUuid", summoner.PlayerUuid, "error", err, "command", contextValue.Command, "request_id", contextValue.RequestId.String())
		return models.SummonerStats{}, err
	}

	respByQueue := make(map[riot.QueueType]*riot.QueueResponse, len(resp))
	for _, res := range resp {
		slog.Debug("Received queue response", "response", fmt.Sprintf("%#v", res), "command", contextValue.Command, "request_id", contextValue.RequestId.String())
		respByQueue[res.QueueType] = res
	}

	return models.SummonerStats{
		Summoner:           summoner,
		FlexGamesPlayed:    respByQueue[riot.RANKED_FLEX].Wins + respByQueue[riot.RANKED_FLEX].Losses,
		FlexTier:           respByQueue[riot.RANKED_FLEX].Tier,
		FlexRank:           respByQueue[riot.RANKED_FLEX].Rank,
		FlexWins:           respByQueue[riot.RANKED_FLEX].Wins,
		FlexLp:             respByQueue[riot.RANKED_FLEX].LeaguePoints,
		SoloDuoGamesPlayed: respByQueue[riot.RANKED_SOLO_DUO].Wins + respByQueue[riot.RANKED_SOLO_DUO].Losses,
		SoloDuoTier:        respByQueue[riot.RANKED_SOLO_DUO].Tier,
		SoloDuoRank:        respByQueue[riot.RANKED_SOLO_DUO].Rank,
		SoloDuoWins:        respByQueue[riot.RANKED_SOLO_DUO].Wins,
		SoloDuoLp:          respByQueue[riot.RANKED_SOLO_DUO].LeaguePoints,
	}, nil
}
