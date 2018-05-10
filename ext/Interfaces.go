package ext

type Sendable interface {
	send() (*Message, error)
}
