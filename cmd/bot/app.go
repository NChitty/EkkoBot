package bot

import (
	"context"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/bot/environment"
	"github.com/NChitty/lol-discord-bot/cmd/domain/adapters"
	"github.com/NChitty/lol-discord-bot/cmd/domain/services"
	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/NChitty/lol-discord-bot/cmd/ports/riot"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	// Base objects

	Context        context.Context
	DiscordSession *discordgo.Session
	Connection     *pgx.Conn
	Queries        *db.Queries

	// Adapters

	GuildRepository services.GuildRepository

	// Services

	GuildService    commands.GuildServicer
	SummonerService commands.SummonerServicer
}

func NewApp(ctx context.Context) (*App, error) {
	// Base objects
	discordToken := environment.GetEnvironmentFileValue("DISCORD_TOKEN_FILE")
	discordSession, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}

	discordSession.StateEnabled = true

	conn, err := pgx.Connect(ctx, "")
	if err != nil {
		return nil, err
	}
	queries := db.New(conn)

	riotToken := environment.GetEnvironmentFileValue("RIOT_API_TOKEN_FILE")
	client, err := riot.NewClient("https://%s.api.riotgames.com/", riot.WithRequestEditorFn(riot.RiotTokenHeader(riotToken)))
	if err != nil {
		return nil, err
	}

	// Adapters
	guildRepository := adapters.NewGuildRepository(queries)
	riotAdapter := adapters.NewHttpRiotAdapter(client)
	summonerRepository := adapters.NewSummonerRepository(queries, conn)

	guildService := services.NewGuildService(guildRepository)
	summonerService := services.NewSummonerService(riotAdapter, summonerRepository)

	return &App{
		Context:         ctx,
		DiscordSession:  discordSession,
		Connection:      conn,
		Queries:         queries,
		GuildService:    guildService,
		SummonerService: summonerService,
	}, nil
}

func (a *App) Start() error {
	cfg, err := pgx.ParseConfig("")
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(stdlib.OpenDB(*cfg), &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"ekkobot",
		driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}

	slog.Debug("Creating commands")
	commands.CreateTrackCommand(a.Context, a.GuildService, a.SummonerService)
	commands.CreateInfoCommand(a.Context, a.SummonerService)
	commands.CommandRegistry.AddHandlers(a.DiscordSession)

	return nil
}
