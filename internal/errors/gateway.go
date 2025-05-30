package errors

type GatewayDisconnectError struct {
	Message         string
	ShouldReconnect bool
}

func newEntry(message string, reconnect bool) GatewayDisconnectError {
	return GatewayDisconnectError{
		Message:         message,
		ShouldReconnect: reconnect,
	}
}

var gatewayDisconnectErrors map[int]GatewayDisconnectError = make(map[int]GatewayDisconnectError)

func init() {
	// register error codes
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
	gatewayDisconnectErrors[4000] = newEntry("Unknown error", true)
	gatewayDisconnectErrors[4001] = newEntry("Unknown opcode", true)
	gatewayDisconnectErrors[4002] = newEntry("Decode error", true)
	gatewayDisconnectErrors[4003] = newEntry("Not authenticated", true)
	gatewayDisconnectErrors[4004] = newEntry("Authentication failed", false)
	gatewayDisconnectErrors[4005] = newEntry("Already authenticated", true)
	gatewayDisconnectErrors[4007] = newEntry("Invalid seq", true)
	gatewayDisconnectErrors[4008] = newEntry("Rate limited", true)
	gatewayDisconnectErrors[4009] = newEntry("Session timed out", true)
	gatewayDisconnectErrors[4010] = newEntry("Invalid shard", false)
	gatewayDisconnectErrors[4011] = newEntry("Sharding required", false)
	gatewayDisconnectErrors[4012] = newEntry("Invalid API version", false)
	gatewayDisconnectErrors[4013] = newEntry("Invalid intent(s)", false)
	gatewayDisconnectErrors[4014] = newEntry("Disallowed intent(s)", false)
}

func GetGatewayErrorByCode(errorCode int) (GatewayDisconnectError, bool) {
	disconnectError, ok := gatewayDisconnectErrors[errorCode]

	return disconnectError, ok
}
