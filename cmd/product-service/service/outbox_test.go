package service

import (
	"encoding/json"
	"github.com/matryer/is"
	"testing"
)

func Test_Outbox(t *testing.T) {
	Is := is.New(t)

	t.Run("should add entry to outbox", func(t *testing.T) {
		p := &mockPublisher{}
		db := NewInMemoryOutbox()
		outbox := NewOutbox(&db, p)

		data := struct {
			foo string
			baz int
		}{
			"bar",
			7,
		}

		msg, _ := json.Marshal(data)

		id, _ := outbox.Send("sample-message", msg)

		Is.Equal(p.messages[0].topic, "sample-message")
		Is.Equal((*db.entries)[id].id, id)
		Is.Equal((*db.entries)[id].sent, true)
	})
}
