package nats

import (
	"L0/storage"
	"log"

	"github.com/nats-io/stan.go"
)

// Вызывается при публикации в канал и сохраняет данные в БД
func MessageCallBack(message *stan.Msg) {
	result := storage.Store(string(message.Data))
	log.Println(result)

}

// Подписка на канал
func Subscribe(clusterID, clientID, channel string, messageCallBack func(*stan.Msg)) (stan.Conn, error) {
	con, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Println(err)
	} else {
		_, err = con.Subscribe(channel, messageCallBack, stan.DeliverAllAvailable())
	}
	return con, err
}
