package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/NChitty/lol-discord-bot/cmd/bot/environment"
	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/NChitty/lol-discord-bot/cmd/ports/riot"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/phsym/console-slog"
)

var (
	Connection     *pgx.Conn
	DiscordToken   string
	DiscordSession *discordgo.Session
	Queries        *db.Queries
)

func init() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	slog.SetDefault(logger)

	DiscordToken = environment.GetEnvironmentFileValue("DISCORD_TOKEN_FILE")
}

func init() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		slog.Error("Could not create a discord session, check the DISCORD_TOKEN_FILE and try again.")
		os.Exit(1)
	}
}

func init() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "")
	if err != nil {
		slog.Error("Could not establish a connection to the database", "error", err)
		os.Exit(1)
	}
	Connection = conn
	Queries = db.New(Connection)
}

func init() {
	cfg, err := pgx.ParseConfig("")
	if err != nil {
		slog.Error("Could not parse database configuration.", "error", err)
		os.Exit(1)
	}
	driver, err := postgres.WithInstance(stdlib.OpenDB(*cfg), &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"ekkobot",
		driver,
	)
	if err != nil {
		slog.Error("Could not create a new migrate instance", "error", err)
		os.Exit(1)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		slog.Error("Could not migrate database", "error", err)
		os.Exit(1)
	}
}

func init() {
	slog.Debug("Creating commands")
	commands.CreateTrackCommand(Queries, riot.RiotClient)
	commands.CommandRegistry.AddHandlers(DiscordSession)
}

func main() {
	ctx := context.Background()
	err := DiscordSession.Open()
	if err != nil {
		slog.Error("Could not open discord session.", "error", err)
		os.Exit(1)
	}

	commands.CommandRegistry.CreateCommands(DiscordSession)

	defer DiscordSession.Close()
	defer Connection.Close(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
