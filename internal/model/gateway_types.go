package model

type Activity struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Url   string `json:"url"`
	State string `json:"state"`
}