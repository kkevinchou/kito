package main

import (
	"time"

	"github.com/kkevinchou/kito/lib/network"
)

func main() {
	host := "localhost"
	port := "8080"
	connectionType := "tcp"

	server := network.NewServer(host, port, connectionType, 18)
	err := server.Start()
	if err != nil {
		panic(err)
	}

	client, err := network.Connect(host, port, connectionType)
	if err != nil {
		panic(err)
	}

	if client.ID() == 0 {
		panic(err)
	}

	err = client.SendMessage(network.MessageTypeInput, nil)
	if err != nil {
		panic(err)
	}

	time.Sleep(10000 * time.Second)
}
