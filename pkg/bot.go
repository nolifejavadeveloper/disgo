package pkg

import (
	"github.com/nolifejavadeveloper/disgo/internal"
	"github.com/nolifejavadeveloper/disgo/internal/event"
	"github.com/rs/zerolog"
)

type DiscordBot interface {
	Start(string) error
	Subscribe(event.EventHandler)
}

func NewDiscordBot(logger *zerolog.Logger) DiscordBot {
	return internal.NewBot(logger)
}

var _ DiscordBot = (*internal.Bot)(nil)