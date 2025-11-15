package commands

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func CreateTrackCommand() {
	command := &discordgo.ApplicationCommand{
		Name: "ping",
		Description: "Pong",
	}
	slog.Debug(fmt.Sprintf("Creating \"%v\" command", command.Name))
	CommandRegistry.registerHandler(
		command,
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong",
				},
			})
			if err != nil {
				slog.Warn("Could not respond to interaction", "command", command.Name, "error", err)
			}
		})
}
