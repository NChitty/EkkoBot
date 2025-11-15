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

			slog.Debug("Checking if already tracking summoner", "name", name, "tag", tag)

			summoner, err := q.GetSummonerByNameAndTag(ctx, db.GetSummonerByNameAndTagParams{
				Name:    pgtype.Text{String: name, Valid: true},
				Tagline: pgtype.Text{String: tag, Valid: true},
			})

			// brand new summoner
			if err != nil && err.Error() == "no rows in result set" {
				summoner, err = q.CreateSummoner(ctx, db.CreateSummonerParams{
					Name:       pgtype.Text{String: name, Valid: true},
					Tagline:    pgtype.Text{String: tag, Valid: true},
					Playeruuid: pgtype.Text{String: "nonsense", Valid: true},
				})

				slog.Info("Started tracking new summoner", "summoner", summoner)

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
					slog.Error(
						"Could not respond to interaction",
						"command", command.Name,
						"error", err,
					)
				}
				return
			}

			// failed to check if summoner exists
			if err != nil {
				slog.Error("Could not check if summoner exists", "name", name, "tag", tag, "error", err)
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "We failed to check if we are already tracking your summoner.",
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

			slog.Debug("Summoner already exists")
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "TODO: SEND CURRENT DATA, WE ARE ALREADY TRACKING",
				},
			})
			if err != nil {
				slog.Error("Could not respond to interaction", "command", command.Name, "error", err)
			}
		})
}
