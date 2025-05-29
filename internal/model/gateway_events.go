package model

type IdentifyEvent struct {
	Token          string                `json:"token"`
	Properties     *ConnectionProperties `json:"properties"`
	Compress       *bool                 `json:"compress,omitempty"`
	LargeThreshold int                   `json:"large_threshold,omitempty"`
	Shard          []int                 `json:"shard,omitempty"`
	Presence       *UpdatePresenceEvent  `json:"presence,omitempty"`
	Intents        int                   `json:"intents,omitempty"`
}

type ConnectionProperties struct {
	Os      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type ResumeEvent struct {
	Token     string `json:"token"`
	SessionId string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type ResumedEvent struct {
}

type ReadyEvent struct {
	V                int    `json:"v"`
	User             *User  `json:"user"`
	Guilds           int    `json:"guilds"`
	SessionId        string `json:"session_id"`
	ResumeGatewayUrl string `json:"resume_gateway_url"`
	Shard            []int  `json:"shard,omitempty"`
	Application      string `json:"application"`
}

type UpdatePresenceEvent struct {
	Since      *int        `json:"since,omitempty"`
	Activities []*Activity `json:"activities"`
	Status     string      `json:"status"`
	Afk        bool        `json:"afk"`
}
