package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/NChitty/lol-discord-bot/cmd/bot"
	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/bwmarrin/discordgo"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/phsym/console-slog"
)


func main() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	slog.SetDefault(logger)

	ctx := context.Background()

	app, err := bot.NewApp(ctx)
	if err != nil {
		slog.Error("Failed to build application", err)
		os.Exit(1)
	}

	if err = app.Start(); err != nil {
		slog.Error("Failed to start application", err)
		os.Exit(1)
	}

	app.DiscordSession.Identify.Intents = discordgo.IntentsGuilds

	app.DiscordSession.AddHandler(func(s *discordgo.Session, i *discordgo.GuildCreate) {
		slog.Info("A new guild has added EkkoBot", "name", i.Guild.Name, "id", i.Guild.ID)
		guildId := pgtype.Text{
			String: i.Guild.ID,
			Valid:  true,
		}
		_, err := app.Queries.GetGuildByDiscordId(ctx, guildId)
		if err != nil && err.Error() == "no rows in result set" {
			_, err := app.Queries.CreateGuild(ctx, guildId)
			if err != nil {
				slog.Error("Failed to insert guild", "name", i.Guild.Name, "id", i.Guild.ID)
			}
		}
	})

	if err = app.DiscordSession.Open(); err != nil {
		slog.Error("Could not open discord session.", "error", err)
		os.Exit(1)
	}

	commands.CommandRegistry.CreateCommands(app.DiscordSession)

	defer app.DiscordSession.Close()
	defer app.Connection.Close(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
