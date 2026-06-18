package commands

import "github.com/google/uuid"

type commandCtxKeyType string

const CONTEXT_KEY commandCtxKeyType = "command"

type CommandContext struct {
	Command   string
	RequestId uuid.UUID
	DiscordId string
}

func newCommandCtxValue(command string, discordId string) CommandContext {
	if id, err := uuid.NewV7(); err != nil {
		return CommandContext{
			Command: command,
			RequestId: uuid.New(),
			DiscordId: discordId,
		}
	} else {
		return CommandContext{
			Command: command,
			RequestId: id,
			DiscordId: discordId,
		}
	}
}
