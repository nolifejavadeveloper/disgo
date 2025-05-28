package model

type Activity struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Url   string `json:"url,omitempty"`
	State string `json:"state,omitempty"`
}