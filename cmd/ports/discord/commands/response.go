package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// Sends a chat to close a command
func SendCommandResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	command *discordgo.ApplicationCommand,
	msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		slog.Error(
			"Could not respond to interaction",
			"command", command.Name,
			"error", err,
		)
	}
}
