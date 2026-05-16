package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/ports/discord"
	"github.com/bwmarrin/discordgo"
)

const TRACK_COMMAND string = "track"

func trackCommand(ctx context.Context, guildService GuildServicer, summonerService SummonerServicer, command *discordgo.ApplicationCommand) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Debug(fmt.Sprintf("Received %s", i.Type.String()), "id", i.GuildID)

		name := i.ApplicationCommandData().GetOption("name").StringValue()
		tag := i.ApplicationCommandData().GetOption("tag").StringValue()
		cmdCtx := context.WithValue(ctx, CONTEXT_KEY, newCommandCtxValue(TRACK_COMMAND, i.GuildID))
		if _, err := guildService.GetGuild(cmdCtx, i.GuildID); err != nil {
			slog.ErrorContext(cmdCtx, "Failed to execute command", "error", err.Error())
			discord.SendCommandResponse(cmdCtx, s, i, command, "Could not complete the request. Reach out to your system administrator for details.")
			return
		} else {
			if stats, err := summonerService.GetSummonerStats(cmdCtx, name, tag); err == nil {
				discord.SendSummonerResponse(cmdCtx, s, i, command, stats)
				return
			} else {
				slog.ErrorContext(cmdCtx, "Failed to execute command", "error", err.Error())
				discord.SendCommandResponse(cmdCtx, s, i, command, "Could not complete the request. Reach out to your system administrator for details.")
				return
			}
		}
	}
}

func CreateTrackCommand(ctx context.Context, guildService GuildServicer, summonerService SummonerServicer) {
	command := &discordgo.ApplicationCommand{
		Name:        TRACK_COMMAND,
		Description: "Start tracking the LP changes of a summoner.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Your summoner name",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tag",
				Description: "Your summoner's tag",
				Required:    true,
			},
		},
	}
	slog.Debug(fmt.Sprintf("Creating \"%v\" command", command.Name))
	CommandRegistry.registerHandler(
		command,
		trackCommand(ctx, guildService, summonerService, command),
	)
}
