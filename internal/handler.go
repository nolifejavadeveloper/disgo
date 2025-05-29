package internal

import (
	"encoding/json"

	"github.com/nolifejavadeveloper/disgo/internal/model"
)

var dispatchHandler messageHandler = func(wc *websocketConn, data []byte, t string) error {
	wc.logger.Debug().Msg("Dispatch recieved from Discord with T: " + t)

	var e any

	switch t {
	case "READY":
		e = &model.ReadyEvent{}
	}

	err := json.Unmarshal(data, e)
	if err != nil {
		return err
	}

	handleDispatch(e, t)

	return nil
}

func handleDispatch(wc *websocketConn, e any, t string) {
	switch t {
	case "READY":
		e = e.(model.ReadyEvent)
		wc.resumeUrl = e.ResumeGatewayUrl
		wc.sessionId = e.SessionId
	}
}

var heartbeatHandler messageHandler = func(wc *websocketConn, data []byte, t string) error {
	return wc.sendHeartbeat()
}

var helloHandler messageHandler = func(wc *websocketConn, data []byte, t string) error {
	wc.logger.Debug().Msg("Hello received from Discord")
	type helloMessage struct {
		HeartbeatInterval int64 `json:"heartbeat_interval"`
	}

	msg := &helloMessage{}

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	wc.heartbeatInterval = msg.HeartbeatInterval
	wc.logger.Debug().Msgf("Heartbeat interval received: %d", msg.HeartbeatInterval)

	wc.startHeartbeat()

	wc.sendIdentify()

	return nil
}

var heartbeatAckHandler messageHandler = func(wc *websocketConn, data []byte, t string) error {
	wc.logger.Debug().Msg("Heartbeat acked")
	wc.receivedHeatbeat()
	return nil
}

func registerHanlders() {
	registerHandler(OpCodeDispatch, dispatchHandler)
	registerHandler(OpCodeHello, helloHandler)
	registerHandler(OpCodeHeartbeat, heartbeatHandler)
	registerHandler(OpCodeAckHeartbeat, heartbeatAckHandler)
}

func init() {
	registerHanlders()
}
