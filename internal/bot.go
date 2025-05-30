package internal

import (
	"fmt"

	"github.com/nolifejavadeveloper/disgo/internal/event"
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

func (b *Bot) Start(token string) error {
	b.conn.token = token
	err := b.conn.connect(discordGateway)
	if err != nil {
		return fmt.Errorf("error connecting to gateway: %s", err.Error())
	}
	b.conn.startReading()

	return nil
}

func (b *Bot) Subscribe(handler ievent.EventHandler) {
	event.Subscribe(b.eventBus, handler)
}
