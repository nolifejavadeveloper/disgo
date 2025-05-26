package pkg

import (
	"github.com/nolifejavadeveloper/disgo/internal/event"
)

type DiscordBot interface {
	Start(string) error
	Subscribe(event.Event)
}



