package internal

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	gerrors "github.com/nolifejavadeveloper/disgo/internal/errors"
	"github.com/nolifejavadeveloper/disgo/internal/event"
	"github.com/nolifejavadeveloper/disgo/internal/model"
	"github.com/rs/zerolog"
)

//https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes

const discordGateway = "wss://gateway.discord.gg/?v=10&encoding=json"

type messageHandler func(wc *websocketConn, data []byte, t string) error

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

type websocketConn struct {
	logger zerolog.Logger
	conn   *websocket.Conn

	bus *event.Bus

	token   string
	intents int

	os      string
	browser string
	device  string

	heartbeatInterval int64
	lastHeartbeat     int64
	quitHeartbeat     chan struct{}

	lastSeq *int

	heartbeatAckReceived bool

	sessionId string
	resumeUrl string

	ready        bool
	shouldResume bool
}

func makeWebsocketConn(logger *zerolog.Logger, bus *event.Bus) *websocketConn {
	newLogger := logger.With().Str("address", discordGateway).Logger()

	return &websocketConn{
		logger: newLogger,

		bus: bus,

		intents: 1,

		os:      "Linux",
		browser: "disgo",
		device:  "disgo",
	}
}

func (wc *websocketConn) connect(addr string) error {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		wc.logger.Error().Err(err).Msg("Error while dialing websocket")
		return err
	}

	wc.logger.Info().Msg("Successfully connected to Discord gateway")

	wc.conn = conn

	return nil
}

func (wc *websocketConn) startReading() {
	go func() {
		for {
			wc.read()
		}
	}()
}

func (wc *websocketConn) read() {
	messageType, data, err := wc.conn.ReadMessage()
	if err != nil {
		if closeErr, ok := err.(*websocket.CloseError); ok {
			wc.handleDisconnect(closeErr.Code, closeErr.Text)
			wc.conn.Close()
			return
		}

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

func (wc *websocketConn) handleDisconnect(code int, msg string) {
	err, ok := gerrors.GetGatewayErrorByCode(code)
	if !ok {
		wc.logger.Error().Msgf("Gateway connection closed with unknown error code: %d", code)
		return
	}

	wc.logger.Warn().Str("gateway_nessage", msg).Int("error_code", code).Str("error_code_definition", err.Message).Msg("Connection closed by discord gateway")
	if err.ShouldReconnect {
		wc.logger.Info().Msg("Reconnecting to resume gateway")
		wc.stopHeartbeat()
		wc.connect(wc.resumeUrl)
		wc.shouldResume = true
		return
	}
}

func (wc *websocketConn) startHeartbeat() {
	go func() {
		wc.receivedHeatbeat()
		wc.logger.Debug().Msg("Starting heartbeat task")
		jitter := rand.Float64()
		firstHeartbeat := int64(float64(wc.heartbeatInterval) * jitter)

		wc.logger.Debug().Msgf("First heartbeat computed: %d", firstHeartbeat)

		wc.lastHeartbeat = time.Now().Local().UnixMilli() - firstHeartbeat

		for {
			select {
			case <-wc.quitHeartbeat:
				return
			default:
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
		}
	}()
}

func (wc *websocketConn) sendHeartbeat() error {
	wc.lastHeartbeat = time.Now().Local().UnixMilli()
	wc.heartbeatAckReceived = false

	return wc.writeEvent(wc.lastSeq, OpCodeHeartbeat, "")
}

func (wc *websocketConn) stopHeartbeat() {
	wc.quitHeartbeat <- struct{}{}
}

func (wc *websocketConn) sendIdentify() error {
	wc.logger.Debug().Msg("Identify sent")
	payload := &model.IdentifyEvent{
		Token: wc.token,
		Properties: &model.ConnectionProperties{
			Os:      wc.os,
			Browser: wc.browser,
			Device:  wc.os,
		},
		Intents: wc.intents,
	}

	return wc.writeEvent(payload, OpCodeIdentify, "")
}

func (wc *websocketConn) receivedHeatbeat() {
	wc.heartbeatAckReceived = true
}

func (wc *websocketConn) write(e *outgoingWebsocketEvent) error {
	s, err := json.Marshal(e)
	if err != nil {
		wc.logger.Error().Err(err).Msgf("Error marshalling JSON")
	}

	wc.logger.Debug().Str("content", string(s)).Msg("Sent packet")
	return wc.conn.WriteMessage(websocket.TextMessage, s)
}

func (wc *websocketConn) sendResume() {
	e := &model.ResumeEvent{
		Token:     wc.token,
		SessionId: wc.sessionId,
		Seq:       *wc.lastSeq,
	}

	wc.writeEvent(e, OpCodeResume, "")
}

func (wc *websocketConn) writeEvent(v any, op int, t string) error {
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
