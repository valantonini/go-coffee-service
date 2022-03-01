package entities

type OutboxEntry struct {
	Id      string
	Topic   string
	Message []byte
	Sent    bool
}
