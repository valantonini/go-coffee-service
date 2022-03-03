package service

import (
	"encoding/json"
	"github.com/matryer/is"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"testing"
	"time"
)

func Test_Outbox(t *testing.T) {
	Is := is.New(t)

	t.Run("should add entry to outbox", func(t *testing.T) {
		p := &mockPublisher{}
		db := data.NewInMemoryOutbox()
		outbox := NewOutbox(&db, p)
		msgData := struct {
			foo string
			baz int
		}{
			"bar",
			7,
		}

		msg, _ := json.Marshal(msgData)

		id, _ := outbox.Send("sample-message", msg)

		Is.Equal(p.messages[0].topic, "sample-message")

		unsent := db.GetUnsent()
		Is.Equal(len(unsent), 1)
		Is.Equal(unsent[0].Id, id)
		Is.Equal(unsent[0].Sent, false)
	})

	t.Run("background polling should send unsent entries in outbox", func(t *testing.T) {
		p := &mockPublisher{}
		db := data.NewInMemoryOutbox()
		outbox := NewOutbox(&db, p)
		cancel := outbox.StartBackgroundPolling(10 * time.Millisecond)
		defer cancel()

		msgData := struct {
			foo string
			baz int
		}{
			"bar",
			7,
		}

		msg, _ := json.Marshal(msgData)

		outbox.Send("sample-message", msg)

		time.Sleep(13 * time.Millisecond)
		Is.Equal(len(db.GetUnsent()), 0)
	})
}
