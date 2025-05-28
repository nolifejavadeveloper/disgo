package pkg

const (
	ActivityTypePlaying = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)

type Activity struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Url   string `json:"url"`
	State string `json:"state"`
}

func NewActivity(name string, typ int, url string, state string) *Activity {
	return &Activity{
		Name: name,
		Type: typ,
		Url: url,
		State: state,
	}
}