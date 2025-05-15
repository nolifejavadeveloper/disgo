package internal

import (
	"github.com/rs/zerolog"
)

type Bot struct {
	logger *zerolog.Logger
	conn   *websocketConn
}

func NewBot(token string, logger *zerolog.Logger) *Bot {
	return &Bot{
		logger: logger,
		conn: makeWebsocketConn(logger, token),
	}
}

func (db *Bot) Start() {
	db.conn.connect()
	db.conn.startReading()
}
