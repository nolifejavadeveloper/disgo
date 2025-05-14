package internal

type IdentifyPayload struct {
	Token          string               `json:"token"`
	Properties     *ConnectionProperties `json:"properties"`
	Compress       bool                 `json:"compress,omitempty"`
	LargeThreshold int                  `json:"large_threshold,omitempty"`
	Shard          []int                `json:"shard,omitempty"`
	Presence       any                  `json:"presence,omitempty"`
	Intents        int                  `json:"intents,omitempty"`
}

type ConnectionProperties struct {
	Os      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type ReadyPayload struct {
	
}