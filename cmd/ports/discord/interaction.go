package discord

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
)

// Sends a chat to close a command
func SendCommandResponse(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	command *discordgo.ApplicationCommand,
	msg string,
) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Could not respond to interaction",
			"error", err,
		)
	}
}

// SummonerEmbed creates and sends a Discord embed message from SummonerStats
func SendSummonerResponse(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	command *discordgo.ApplicationCommand,
	stats models.SummonerStats,
) {
	fields := make([]*discordgo.MessageEmbedField, 0)

	// Add Solo/Duo queue info
	if stats.SoloDuoGamesPlayed > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Solo/Duo",
			Value:  fmt.Sprintf("%s %s", stats.SoloDuoTier, stats.SoloDuoRank),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Solo/Duo Record",
			Value:  fmt.Sprintf("%dW %dL", stats.SoloDuoWins, stats.SoloDuoGamesPlayed-stats.SoloDuoWins),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Solo/Duo LP",
			Value:  fmt.Sprintf("%d", stats.SoloDuoLp),
			Inline: true,
		})
	}

	// Add Flex queue info
	if stats.FlexGamesPlayed > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Flex",
			Value:  fmt.Sprintf("%s %s", stats.FlexTier, stats.FlexRank),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Flex Record",
			Value:  fmt.Sprintf("%dW %dL", stats.FlexWins, stats.FlexGamesPlayed-stats.FlexWins),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Flex LP",
			Value:  fmt.Sprintf("%d", stats.FlexLp),
			Inline: true,
		})
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Summoner Stats",
					Description: fmt.Sprintf("Statistics for **%s#%s**", stats.Summoner.Name, stats.Summoner.TagLine),
					Fields:      fields,
				},
			},
		},
	})
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Could not respond to interaction",
			"error", err,
		)
	}
}
