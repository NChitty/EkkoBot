package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NChitty/lol-discord-bot/cmd/ports/db"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateTrackCommand(q *db.Queries) {
	command := &discordgo.ApplicationCommand{
		Name:        "track",
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
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			name := i.ApplicationCommandData().GetOption("name").StringValue()
			tag := i.ApplicationCommandData().GetOption("tag").StringValue()
			ctx := context.Background()

			_, err := q.CreateSummoner(ctx, db.CreateSummonerParams{
				Name:       pgtype.Text{String: name, Valid: true},
				Tagline:    pgtype.Text{String: tag, Valid: true},
				Playeruuid: pgtype.Text{String: "nonsense", Valid: true},
			})

			if err != nil {
				slog.Error(
					"Could not insert new summoner",
					"summoner", fmt.Sprintf("%s#%s", name, tag),
					"error", err,
				)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to add your summoner, check the name and tag and try again.",
					},
				})
				if err != nil {
					slog.Error(
						"Could not respond to interaction",
						"command", command.Name,
						"error", err,
					)
				}

				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Added summoner, we will start tracking your GAINS",
				},
			})
			if err != nil {
				slog.Warn("Could not respond to interaction", "command", command.Name, "error", err)
			}
		})
}
