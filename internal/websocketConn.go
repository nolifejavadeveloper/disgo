package internal

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const discordGateway = "wss://gateway.discord.gg/?v=10&encoding=json"

type messageHandler func(wc *WebsocketConn, data []byte, t string) error

type websocketEvent struct {
	Op int    `json:"op"`
	S  *int   `json:"s,omitempty"`
	T  string `json:"t,omitempty"`
}

type outgoingWebsocketEvent struct {
	websocketEvent
	D interface{} `json:"d"`
}

type incomingWebsocketEvent struct {
	websocketEvent
	D json.RawMessage `json:"d"`
}

var messageHandlers = make(map[int]messageHandler)

type WebsocketConn struct {
	logger zerolog.Logger
	conn   *websocket.Conn

	token   string
	intents int

	os      string
	browser string
	device  string

	heartbeatInterval int64
	lastHeartbeat     int64

	lastSeq *int

	heartbeatAckReceived bool

	resumeUrl string

	ready bool
}

func MakeWebsocketConn(logger *zerolog.Logger, token string) *WebsocketConn {
	newLogger := logger.With().Str("address", discordGateway).Logger()
	conn, _, err := websocket.DefaultDialer.Dial(discordGateway, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Error while dialing websocket")
		return nil
	}

	logger.Info().Msg("Successfully connected to Discord gateway")

	return &WebsocketConn{
		logger: newLogger,
		conn:   conn,

		token:   token,
		intents: 1,

		os:      "Linux",
		browser: "disgo",
		device:  "disgo",
	}
}

func (wc *WebsocketConn) StartReading() {
	go func() {
		for {
			wc.Read()
		}
	}()
}

func (wc *WebsocketConn) Read() {
	messageType, data, err := wc.conn.ReadMessage()
	if err != nil {
		wc.logger.Error().Err(err).Msg("Error reading from websocket")
		return
	}

	if messageType == websocket.TextMessage {
		msg := &incomingWebsocketEvent{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			wc.logger.Error().Err(err).Msg("Error unmarshalling JSON")
			return
		}

		handler, ok := messageHandlers[msg.Op]
		if !ok {
			wc.logger.Error().Int("opcode", msg.Op).Msg("No message handler found")
			return
		}

		if msg.S != nil {
			wc.lastSeq = msg.S
		}

		err = handler(wc, msg.D, msg.T)

		if err != nil {
			wc.logger.Error().Int("opcode", msg.Op).Err(err).Msg("Failed to handle message")
		}
	}
}

func (wc *WebsocketConn) StartHeartbeat() {
	go func() {
		wc.receivedHeatbeat()

		// TODO: move to crypto/rand for more random float64 generating
		wc.logger.Debug().Msg("Starting heartbeat task")
		jitter := rand.Float64()
		firstHeartbeat := int64(float64(wc.heartbeatInterval) * jitter)

		wc.logger.Debug().Msgf("First heartbeat computed: %d", firstHeartbeat)

		wc.lastHeartbeat = time.Now().Local().UnixMilli() - firstHeartbeat

		for {
			if time.Now().Local().UnixMilli()-wc.lastHeartbeat > wc.heartbeatInterval {
				if !wc.heartbeatAckReceived {
					// reconnect
					wc.logger.Warn().Msg("Last heartbeat did not receive ack, reconnecting...")
					return
				}

				err := wc.sendHeartbeat()
				if err != nil {
					wc.logger.Error().Err(err).Msg("Error sending timmed heartbeat")
					continue
				}
				wc.logger.Debug().Msg("Sucessfully sent timed heartbeat")
			}
		}
	}()
}

func (wc *WebsocketConn) sendHeartbeat() error {
	wc.lastHeartbeat = time.Now().Local().UnixMilli()
	wc.heartbeatAckReceived = false

	return wc.writeEvent(wc.lastSeq, OpCodeHeartbeat, "")
}

func (wc *WebsocketConn) sendIdentify() error {
	wc.logger.Debug().Msg("Identify sent")
	payload := &IdentifyPayload{
		Token: wc.token,
		Properties: &ConnectionProperties{
			Os:      wc.os,
			Browser: wc.browser,
			Device:  wc.os,
		},
		Intents: wc.intents,
	}

	return wc.writeEvent(payload, OpCodeIdentify, "")
}

func (wc *WebsocketConn) receivedHeatbeat() {
	wc.heartbeatAckReceived = true
}

func (wc *WebsocketConn) write(e *outgoingWebsocketEvent) error {
	s, err := json.Marshal(e)
	if err != nil {
		wc.logger.Error().Err(err).Msgf("Error marshalling JSON")
	}

	wc.logger.Debug().Str("content", string(s)).Msg("Sent packet")
	return wc.conn.WriteMessage(websocket.TextMessage, s)
}

func (wc *WebsocketConn) writeEvent(v any, op int, t string) error {

	e := &outgoingWebsocketEvent{
		websocketEvent: websocketEvent{
			Op: op,
		},
		D: v,
	}

	if t != "" {
		e.T = t
	}

	return wc.write(e)
}

func registerHandler(op int, handler messageHandler) {
	messageHandlers[op] = handler
}
