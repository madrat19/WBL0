package run

import (
	"L0/nats"
	"L0/server"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func Run() error {

	// Запуск http сервера
	go server.RunServer()
	log.Println("Http server is running")

	// Запуск nats-streaming
	filePath := `C:\Users\chemo\go\src\github.com\nats-io\nats-streaming-server\nats-streaming-server.go`
	cmd := exec.Command("go", "run", filePath)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	// Ожидание запуска nats-streaming
	time.Sleep(5 * time.Second)
	log.Println("Nats is running")

	// Подписка на канал в nats
	conn, err := nats.Subscribe("test-cluster", "subscriber", "channel-1", nats.MessageCallBack)
	if err != nil {
		return err
	}

	fmt.Scanln()

	defer conn.Close()
	return err
}
