package conversation

type KeyStrategy int64

// Note: If you add a new keystrategy here, make sure to add it to the getStateKey method!
const (
	// KeyStrategySenderAndChat ensures that each sender get a unique conversation in each chats.
	KeyStrategySenderAndChat KeyStrategy = iota
	// KeyStrategySender gives a unique conversation to each sender, but that conversation is available in all chats.
	KeyStrategySender
	// KeyStrategyChat gives a unique conversation to each chat, which all senders can interact in together.
	KeyStrategyChat
)
