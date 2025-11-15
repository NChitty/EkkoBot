package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/NChitty/lol-discord-bot/cmd/ports/discord/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/phsym/console-slog"
)

var (
	DiscordToken   string
	DiscordSession *discordgo.Session
)

func init() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	slog.SetDefault(logger)

	discordToken, isPresent := os.LookupEnv("DISCORD_TOKEN")
	if !isPresent {
		slog.Error("DISCORD_TOKEN environment variable is unset.")
		os.Exit(1)
	}
	if discordToken == "" {
		slog.Error("DISCORD_TOKEN environment variable is empty.")
		os.Exit(1)
	}
	DiscordToken = discordToken
}

func init() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		slog.Error("Could not create a discord session, check the DISCORD_TOKEN and try again.")
		os.Exit(1)
	}
}

func init() {
	slog.Debug("Creating commands")
	commands.CreateTrackCommand()
	commands.CommandRegistry.AddHandlers(DiscordSession)
}

func main() {
	err := DiscordSession.Open()
	if err != nil {
		slog.Error("Could not open discord session.", "error", err)
		os.Exit(1)
	}

	commands.CommandRegistry.CreateCommands(DiscordSession)

	defer DiscordSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
