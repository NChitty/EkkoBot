package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var CommandRegistry *commandRegistry

func init() {
	CommandRegistry = new()
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command interface {
	getCommand() CommandHandler
}

type commandRegistry struct {
	commands           []*discordgo.ApplicationCommand
	commandHandlers    map[string]CommandHandler
}

func new() *commandRegistry {
	return &commandRegistry{
		commands:        make([]*discordgo.ApplicationCommand, 0),
		commandHandlers: make(map[string]CommandHandler),
	}
}

func (cr *commandRegistry) registerHandler(
	command *discordgo.ApplicationCommand,
	handler CommandHandler,
) {
	cr.commands = append(cr.commands, command)
	cr.commandHandlers[command.Name] = handler
	slog.Debug("Registered command", "command", command.Name)
}

func (cr *commandRegistry) AddHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := cr.commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (cr *commandRegistry) CreateCommands(s *discordgo.Session) {
	for _, command := range cr.commands {
		slog.Debug("Creating command with discord", "command", command.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			slog.Error("Cannot create command", "command", command.Name, "error", err)
		}
	}
}
