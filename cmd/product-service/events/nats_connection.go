package events

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

func NewNatsConnection(natsAddress string) (*nats.Conn, error) {
	startTime := time.Now()
	backoff := 1 * time.Second      // this should be an exponential backoff
	maxWaitTime := 45 * time.Second // max time to wait of the DB connection

	for {
		var err error
		nc, err := nats.Connect(natsAddress)

		if err != nil {
			if time.Now().Sub(startTime) > maxWaitTime {
				return nil, err
			}
			fmt.Println("error connecting to nats. backing off")
			time.Sleep(backoff)
			continue
		}

		return nc, nil
	}
}
