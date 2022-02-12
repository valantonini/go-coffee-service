package events

// Publisher interacts with the message broker
type Publisher interface {
	Publish(topic string, data []byte) error
}
