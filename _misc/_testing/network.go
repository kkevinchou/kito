package main

import (
	"time"

	"github.com/kkevinchou/kito/lib/network"
)

func main() {
	go func() {
		s := network.Server{}
		s.Start()
	}()

	m := &network.Message{SenderID: 13, MessageType: 1}
	c := network.Client{}
	c.Connect()

	for i := 0; i < 10; i++ {
		err := c.SendMessage(m)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}
