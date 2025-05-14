package internal

import (
	"encoding/json"
)

var dispatchHandler messageHandler = func(wc *WebsocketConn, data []byte, t string) error {
	wc.logger.Debug().Msg("Dispatch recieved from Discord with T: " + t)
	return nil
}

var heartbeatHandler messageHandler = func(wc *WebsocketConn, data []byte, t string) error {
	return wc.sendHeartbeat()
}

var helloHandler messageHandler = func(wc *WebsocketConn, data []byte, t string) error {
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

	wc.StartHeartbeat()

	wc.sendIdentify()

	return nil
}

var heartbeatAckHandler messageHandler = func(wc *WebsocketConn, data []byte, t string) error {
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
