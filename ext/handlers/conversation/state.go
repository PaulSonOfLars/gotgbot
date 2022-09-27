package conversation

// State stores all the variables relevant to the current conversation state.
//
// Note: More keys may be added in the future to support additional features.
// As such, any storage implementations should be flexible, and allow for storing the entire struct rather than
// individual fields.
type State struct {
	// Key represents the name of the current state, as defined in the States map of handlers.Conversation.
	Key string
}
