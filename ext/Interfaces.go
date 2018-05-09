package ext

import "github.com/PaulSonOfLars/gotgbot/types"

type Sendable interface {
	send() (*types.Message, error)
}
