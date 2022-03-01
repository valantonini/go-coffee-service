package entities

type OutboxEntry struct {
	Id      string `bson:"_id"`
	Topic   string `bson:"topic"`
	Message []byte `bson:"message"`
	Sent    bool   `bson:"sent"`
}
