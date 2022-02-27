package service

type message struct {
	topic string
	msg   []byte
}

type mockPublisher struct {
	messages []message
}

func (m *mockPublisher) Publish(topic string, msg []byte) error {
	m.messages = append(m.messages, message{topic, msg})
	return nil
}
