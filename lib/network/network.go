package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

const (
	messageQueueBufferSize        = 1024
	incomingConnectionsBufferSize = 1024

	UnsetClientID = -1
)

type Connection struct {
	ID         int
	Connection net.Conn
}

func queueIncomingMessages(conn net.Conn, messageQueue chan *Message) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		message := Message{}
		err := decoder.Decode(&message)
		if err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println("error reading:", err.Error())
			fmt.Println("closing connection")
			return
		}

		select {
		case messageQueue <- &message:
		default:
			fmt.Println("message queue full")
		}
	}
}
