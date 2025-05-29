package pkg

import "github.com/nolifejavadeveloper/disgo/internal/model"

const (
	ActivityTypePlaying = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)

type Activity struct {
	name string
	typ  int
	url  string
}

func NewActivity(name string, typ int, url string, state string) *Activity {
	return &Activity{
		name: name,
		typ:  typ,
		url:  url,
	}
}

func ActivityPlaying(name string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypePlaying,
	}
}

func ActivityStreaming(name string, url string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypeStreaming,
		url:  url,
	}
}

func ActivityListening(name string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypeListening,
	}
}

func ActivityWatching(name string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypeWatching,
	}
}

func ActivityCustom(name string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypeCustom,
	}
}

func ActivityCompeting(name string) *Activity {
	return &Activity{
		name: name,
		typ:  ActivityTypeCompeting,
	}
}

func (a *Activity) toModel() *model.Activity {
	return &model.Activity{
		Name:  a.name,
		Type: a.typ,
		Url: a.url,
	}
}