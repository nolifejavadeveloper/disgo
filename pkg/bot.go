package pkg

import (
	"github.com/nolifejavadeveloper/disgo/internal"
	"github.com/rs/zerolog"
)

type DiscordBot struct {
	bot *internal.Bot
}

func NewDiscordBot(token string, logger *zerolog.Logger) *DiscordBot {
	return &DiscordBot{
		bot: internal.NewBot(token, logger),
	}
}



func (db *DiscordBot) Start(token string) error {
	return db.bot.Start()
}
