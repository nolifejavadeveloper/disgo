package internal

import (
	"fmt"

	ievent "github.com/nolifejavadeveloper/disgo/internal/event"
	"github.com/rs/zerolog"
)

type Bot struct {
	logger   *zerolog.Logger
	conn     *websocketConn
	eventBus *ievent.Bus
}

func NewBot(logger *zerolog.Logger) *Bot {
	bus := ievent.NewBus()
	return &Bot{
		logger: logger,
		conn:   makeWebsocketConn(logger, bus),
	}
}

func (db *Bot) Start(token string) error {
	db.conn.token = token
	err := db.conn.connect()
	if err != nil {
		return fmt.Errorf("error connecting to gateway: %s", err.Error())
	}
	db.conn.startReading()

	return nil
}
