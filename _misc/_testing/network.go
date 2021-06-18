package main

import (
	"time"

	"github.com/kkevinchou/kito/lib/network"
)

func main() {
	host := "localhost"
	port := "8080"
	connectionType := "tcp"

	server := network.NewServer(host, port, connectionType)
	err := server.Start()
	if err != nil {
		panic(err)
	}

	client := network.NewClient()
	acceptMessage, err := client.Connect(host, port, connectionType)
	if err != nil {
		panic(err)
	}

	if acceptMessage.PlayerID == 0 {
		panic(err)
	}

	err = client.SendMessage(&network.Message{MessageType: network.MessageTypeInput})
	if err != nil {
		panic(err)
	}

	time.Sleep(10000 * time.Second)
}
