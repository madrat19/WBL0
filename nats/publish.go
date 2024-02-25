package nats

import (
	"log"

	"github.com/nats-io/stan.go"
)

// Публикация в канал
func Publish(clusterID, clientID, channel, message string) error { // test-cluster
	con, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatal(err)
	} else {
		err = con.Publish(channel, []byte(message))
	}

	defer con.Close()
	return err
}
