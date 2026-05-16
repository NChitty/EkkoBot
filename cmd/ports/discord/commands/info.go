package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/ports/discord"
	"github.com/bwmarrin/discordgo"
)

const INFO_COMMAND string = "stats"

func infoCommand(ctx context.Context, guildService GuildServicer, summonerService SummonerServicer, command *discordgo.ApplicationCommand) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Debug(fmt.Sprintf("Received %s", i.Type.String()), "id", i.GuildID)

		name := i.ApplicationCommandData().GetOption("name").StringValue()
		tag := i.ApplicationCommandData().GetOption("tag").StringValue()
		cmdCtx := context.WithValue(ctx, CONTEXT_KEY, newCommandCtxValue(INFO_COMMAND, i.GuildID))
		// Add guild
		err := guildService.CreateGuild(cmdCtx, i.GuildID)
		if err != nil {
			slog.ErrorContext(cmdCtx, "Could not create guild", "error", err.Error())
		}

		if stats, err := summonerService.GetSummonerStats(cmdCtx, name, tag, i.GuildID); err == nil {
			discord.SendSummonerResponse(cmdCtx, s, i, command, stats)
		} else {
			slog.ErrorContext(cmdCtx, "Failed to execute command", "error", err.Error())
			discord.SendCommandResponse(cmdCtx, s, i, command, "Could not complete the request. Reach out to your system administrator for details.")
			return
		}
	}
}

func CreateInfoCommand(ctx context.Context, guildService GuildServicer, summonerService SummonerServicer) {
	command := &discordgo.ApplicationCommand{
		Name:        INFO_COMMAND,
		Description: "Get the current ranked stats for the summoner",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Summoner name",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tag",
				Description: "Summoner's tag",
				Required:    true,
			},
		},
	}
	slog.Debug(fmt.Sprintf("Creating \"%v\" command", command.Name))
	CommandRegistry.registerHandler(
		command,
		infoCommand(ctx, guildService, summonerService, command),
	)
}
