package messaging

type Message interface {
	Identifier() string
}
