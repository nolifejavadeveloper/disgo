package pkg

import (
	"os"

	"github.com/nolifejavadeveloper/disgo/internal"

	"github.com/rs/zerolog"
)

type DiscordBot struct {
	conn *internal.WebsocketConn
}

func NewDiscordBot() *DiscordBot {
	return &DiscordBot{
		
	}
}

func (db *DiscordBot) Setup(token string) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db.conn = internal.MakeWebsocketConn(&logger, token)

	db.conn.StartReading()
}
