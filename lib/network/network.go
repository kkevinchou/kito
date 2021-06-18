package network

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	messageQueueBufferSize = 1024
)

func queueIncomingMessages(conn net.Conn, messageQueue chan *Message) {
	decoder := json.NewDecoder(conn)
	for {
		message := Message{}
		err := decoder.Decode(&message)
		if err != nil {
			fmt.Println("error reading:", err.Error())
			continue
		}

		select {
		case messageQueue <- &message:
		default:
			fmt.Println("message queue full")
		}
	}
}
