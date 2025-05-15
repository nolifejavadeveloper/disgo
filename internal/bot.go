package internal

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Bot struct {
	logger *zerolog.Logger
	conn   *websocketConn
}

func NewBot(token string, logger *zerolog.Logger) *Bot {
	return &Bot{
		logger: logger,
		conn:   makeWebsocketConn(logger, token),
	}
}

func (db *Bot) Start() error {
	err := db.conn.connect()
	if err != nil {
		return fmt.Errorf("error connecting to gateway: %s", err.Error())
	}
	db.conn.startReading()

	return nil
}
