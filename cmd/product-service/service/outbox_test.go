package service

import (
	"encoding/json"
	"github.com/matryer/is"
	"testing"
)

func Test_Outbox(t *testing.T) {
	Is := is.New(t)

	t.Run("should add entry to outbox", func(t *testing.T) {
		var mockPublisher mockPublisher
		outbox := NewOutbox(&mockPublisher)

		data := struct {
			foo string
			baz int
		}{
			"bar",
			7,
		}

		msg, _ := json.Marshal(data)

		outbox.Send("sample-message", msg)

		Is.Equal(mockPublisher.messages[0].topic, "sample-message")
	})
}
