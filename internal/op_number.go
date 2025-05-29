package internal

// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes

const (
	OpCodeDispatch = 0
	OpCodeHeartbeat = 1
	OpCodeIdentify = 2
	OpCodePresenceUpdate = 3
	OpCodeVoiceStateUpdate = 4
	OpCodeResume = 6
	OpCodeReconnect = 7
	OpCodeRequestGuildMembers = 8
	OpCodeInvalidSession = 9
	OpCodeHello = 10
	OpCodeAckHeartbeat = 11

	OpCodeRequestSoundboardSounds = 31
)