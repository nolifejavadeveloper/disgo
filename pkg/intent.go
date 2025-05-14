package pkg

const (
	IntentGuilds                      = 0
	IntentGuildMembers                = 1
	IntentGuildModeration             = 2
	IntentGuildExpressions            = 3
	IntentGuildIntergrations          = 4
	IntentGuildWebhooks               = 5
	IntentGuildInvites                = 6
	IntentGuildVoiceStates            = 7
	IntentGuildPresences              = 8
	IntentGuildMessages               = 9
	IntentGuildMessageReactions       = 10
	IntentGuildMessageTyping          = 11
	IntentDirectMessages              = 12
	IntentDirectMessageReactions      = 13
	IntentDirectMessageTyping         = 14
	IntentDirectContent               = 15
	IntentGuildScheduledEvents        = 16
	IntentAutoModerationConfiguration = 20
	IntentAutoModerationExecution     = 21
	IntentGuildMessagePolls           = 24
	IntentDirectMessagePolls          = 25
)

var allIntents = []int{IntentGuilds, IntentGuildMembers, IntentGuildModeration, 
	IntentGuildExpressions, IntentGuildIntergrations, IntentGuildWebhooks, 
	IntentGuildInvites, IntentGuildVoiceStates, IntentGuildPresences,
		IntentGuildMessages, IntentGuildMessageReactions, IntentGuildMessageTyping, 
		IntentDirectMessages, IntentDirectMessageReactions, IntentDirectMessageTyping, 
		IntentDirectContent, IntentGuildScheduledEvents, IntentAutoModerationConfiguration, 
		IntentAutoModerationExecution, IntentGuildMessagePolls, IntentDirectMessagePolls,
	}

func BuildIntent(intents ...int) int {
	var val int = 0
	for _, intent := range intents {
		val |= 1 << intent
	}

	return val
}

func AllIntents() int {
	return BuildIntent(allIntents...)
}
