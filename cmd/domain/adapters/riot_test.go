package adapters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/NChitty/lol-discord-bot/cmd/ports/riot"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRiotClient implements RiotClientInterface for testing
type mockRiotClient struct {
	t              *testing.T
	accountResp    *riot.AccountByRiotIdResponse
	accountErr     error
	queueResp      []*riot.QueueResponse
	queueErr       error
	expectedParams riot.AccountByRiotIdRequestParams
}

func (m *mockRiotClient) GetAccountByRiotId(ctx context.Context, params riot.AccountByRiotIdRequestParams) (*riot.AccountByRiotIdResponse, error) {
	m.expectedParams = params
	return m.accountResp, m.accountErr
}

func (m *mockRiotClient) GetQueueEntriesByPlayerUuid(ctx context.Context, params riot.QueueEntriesByPlayerUuidParams) ([]*riot.QueueResponse, error) {
	return m.queueResp, m.queueErr
}

func TestGetSummonerSuccess(t *testing.T) {
	ctx := context.Background()

	expectedSummoner := &riot.AccountByRiotIdResponse{
		PlayerUuid: "test-puuid-123",
		Name:       "TestSummoner",
		Tagline:    "NA1",
	}

	mockClient := &mockRiotClient{
		t:           t,
		accountResp: expectedSummoner,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	summoner, err := adapter.GetSummoner(ctx, "TestSummoner", "NA1")

	require.NoError(t, err)
	assert.Equal(t, "TestSummoner", summoner.Name)
	assert.Equal(t, "NA1", summoner.TagLine)
	assert.Equal(t, "test-puuid-123", summoner.PlayerUuid)
	assert.Equal(t, "TestSummoner", mockClient.expectedParams.Name)
	assert.Equal(t, "NA1", mockClient.expectedParams.Tagline)
}

func TestGetSummonerTimeout(t *testing.T) {
	// Create a context that's already cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	cancel() // Cancel immediately

	mockClient := &mockRiotClient{
		t:          t,
		accountErr: context.DeadlineExceeded,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	_, err := adapter.GetSummoner(ctx, "TestSummoner", "NA1")

	assert.Error(t, err)
	// Check for deadline exceeded error (message may vary by Go version)
	isDeadlineErr := err.Error() == "context deadline exceeded" || err.Error() == "context: deadline exceeded" || err == context.DeadlineExceeded
	assert.True(t, isDeadlineErr, "error should be deadline exceeded, got: %v", err)
}

func TestGetSummonerRiotError(t *testing.T) {
	ctx := context.Background()

	mockClient := &mockRiotClient{
		t:          t,
		accountErr: assert.AnError,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	_, err := adapter.GetSummoner(ctx, "TestSummoner", "NA1")

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestGetRankedStatsSuccess(t *testing.T) {
	ctx := context.WithValue(context.Background(), commands.CONTEXT_KEY, commands.CommandContext{
		Command:   "test-command",
		RequestId: uuid.New(),
		DiscordId: "123456",
	})

	summoner := models.Summoner{
		ID:         -1,
		Name:       "TestSummoner",
		TagLine:    "NA1",
		PlayerUuid: "test-puuid-123",
	}

	queueResponses := []*riot.QueueResponse{
		{
			QueueType:    riot.RANKED_FLEX,
			LeagueId:     "test-league-flex",
			Tier:         "GOLD",
			Rank:         "I",
			PlayerUuid:   "test-puuid-123",
			LeaguePoints: 50,
			Wins:         15,
			Losses:       10,
			Veteran:      false,
			Inactive:     false,
			FreshBlood:   true,
			HotStreak:    false,
		},
		{
			QueueType:    riot.RANKED_SOLO_DUO,
			LeagueId:     "test-league-solo",
			Tier:         "PLATINUM",
			Rank:         "II",
			PlayerUuid:   "test-puuid-123",
			LeaguePoints: 75,
			Wins:         25,
			Losses:       20,
			Veteran:      true,
			Inactive:     false,
			FreshBlood:   false,
			HotStreak:    true,
		},
	}

	mockClient := &mockRiotClient{
		t:          t,
		accountResp: &riot.AccountByRiotIdResponse{PlayerUuid: summoner.PlayerUuid, Name: summoner.Name, Tagline: summoner.TagLine},
		queueResp:  queueResponses,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	stats, err := adapter.GetRankedStats(ctx, summoner)

	require.NoError(t, err)
	assert.Equal(t, summoner, stats.Summoner)
	assert.Equal(t, 25, stats.FlexGamesPlayed)
	assert.Equal(t, "GOLD", stats.FlexTier)
	assert.Equal(t, "I", stats.FlexRank)
	assert.Equal(t, 15, stats.FlexWins)
	assert.Equal(t, 50, stats.FlexLp)
	assert.Equal(t, 45, stats.SoloDuoGamesPlayed)
	assert.Equal(t, "PLATINUM", stats.SoloDuoTier)
	assert.Equal(t, "II", stats.SoloDuoRank)
	assert.Equal(t, 25, stats.SoloDuoWins)
	assert.Equal(t, 75, stats.SoloDuoLp)
}

func TestGetRankedStatsEmptyQueues(t *testing.T) {
	ctx := context.WithValue(context.Background(), commands.CONTEXT_KEY, commands.CommandContext{
		Command:   "test-command",
		RequestId: uuid.New(),
		DiscordId: "123456",
	})

	summoner := models.Summoner{
		ID:         -1,
		Name:       "TestSummoner",
		TagLine:    "NA1",
		PlayerUuid: "test-puuid-123",
	}

	mockClient := &mockRiotClient{
		t:         t,
		queueResp: []*riot.QueueResponse{},
	}

	adapter := NewHttpRiotAdapter(mockClient)

	_, err := adapter.GetRankedStats(ctx, summoner)

	// This should return an error because we try to access map keys on empty responses
	assert.Error(t, err)
}

func TestGetRankedStatsMissingQueue(t *testing.T) {
	ctx := context.WithValue(context.Background(), commands.CONTEXT_KEY, commands.CommandContext{
		Command:   "test-command",
		RequestId: uuid.New(),
		DiscordId: "123456",
	})

	summoner := models.Summoner{
		ID:         -1,
		Name:       "TestSummoner",
		TagLine:    "NA1",
		PlayerUuid: "test-puuid-123",
	}

	// Only provide one queue type
	queueResponses := []*riot.QueueResponse{
		{
			QueueType:    riot.RANKED_SOLO_DUO,
			LeagueId:     "test-league-solo",
			Tier:         "PLATINUM",
			Rank:         "II",
			PlayerUuid:   "test-puuid-123",
			LeaguePoints: 75,
			Wins:         25,
			Losses:       20,
		},
	}

	mockClient := &mockRiotClient{
		t:         t,
		queueResp: queueResponses,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	stats, err := adapter.GetRankedStats(ctx, summoner)

	require.NoError(t, err)
	// Solo/Duo should have data
	assert.Equal(t, 45, stats.SoloDuoGamesPlayed)
	assert.Equal(t, "PLATINUM", stats.SoloDuoTier)
	// Flex should be default/zero values
	assert.Equal(t, 0, stats.FlexGamesPlayed)
	assert.Equal(t, "", stats.FlexTier)
}

func TestGetRankedStatsRiotError(t *testing.T) {
	ctx := context.WithValue(context.Background(), commands.CONTEXT_KEY, commands.CommandContext{
		Command:   "test-command",
		RequestId: uuid.New(),
		DiscordId: "123456",
	})

	summoner := models.Summoner{
		ID:         -1,
		Name:       "TestSummoner",
		TagLine:    "NA1",
		PlayerUuid: "test-puuid-123",
	}

	mockClient := &mockRiotClient{
		t:        t,
		queueErr: assert.AnError,
	}

	adapter := NewHttpRiotAdapter(mockClient)

	_, err := adapter.GetRankedStats(ctx, summoner)

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

// Test with a real HTTP server to test actual integration with riot.Client
func TestHttpRiotAdapterIntegration(t *testing.T) {
	// Create a test server to mock Riot API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/riot/account/v1/accounts/by-riot-id/TestSummoner/NA1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"puuid": "test-puuid-integration",
				"gameName": "TestSummoner",
				"tagLine": "NA1"
			}`))
		case "/lol/league/v4/entries/by-puuid/test-puuid-integration":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"leagueId": "test-league-flex",
					"queueType": "RANKED_FLEX_SR",
					"tier": "GOLD",
					"rank": "I",
					"puuid": "test-puuid-integration",
					"leaguePoints": 50,
					"wins": 15,
					"losses": 10
				},
				{
					"leagueId": "test-league-solo",
					"queueType": "RANKED_SOLO_5x5",
					"tier": "PLATINUM",
					"rank": "II",
					"puuid": "test-puuid-integration",
					"leaguePoints": 75,
					"wins": 25,
					"losses": 20
				}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create a riot client pointing to our test server
	riotClient, err := riot.NewClient(server.URL, riot.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("X-Riot-Token", "test-token")
		return nil
	}))
	require.NoError(t, err)

	adapter := NewHttpRiotAdapter(riotClient)

	// Test GetSummoner
	ctx := context.WithValue(context.Background(), commands.CONTEXT_KEY, commands.CommandContext{
		Command:   "test-command",
		RequestId: uuid.New(),
		DiscordId: "123456",
	})
	summoner, err := adapter.GetSummoner(ctx, "TestSummoner", "NA1")

	require.NoError(t, err)
	assert.Equal(t, "TestSummoner", summoner.Name)
	assert.Equal(t, "NA1", summoner.TagLine)
	assert.Equal(t, "test-puuid-integration", summoner.PlayerUuid)

	// Test GetRankedStats
	stats, err := adapter.GetRankedStats(ctx, summoner)

	require.NoError(t, err)
	assert.Equal(t, summoner, stats.Summoner)
	assert.Equal(t, 25, stats.FlexGamesPlayed)
	assert.Equal(t, "GOLD", stats.FlexTier)
	assert.Equal(t, "I", stats.FlexRank)
	assert.Equal(t, 15, stats.FlexWins)
	assert.Equal(t, 50, stats.FlexLp)
	assert.Equal(t, 45, stats.SoloDuoGamesPlayed)
	assert.Equal(t, "PLATINUM", stats.SoloDuoTier)
	assert.Equal(t, "II", stats.SoloDuoRank)
	assert.Equal(t, 25, stats.SoloDuoWins)
	assert.Equal(t, 75, stats.SoloDuoLp)
}
