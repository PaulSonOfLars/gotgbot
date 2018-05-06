package ext

import "gotgbot/types"

type Sendable interface {
	send() (*types.Message, error)
}
