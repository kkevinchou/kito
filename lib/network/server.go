package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

const (
	messageQueueSize = 1024
)

type Server struct {
	host           string
	port           string
	connectionType string

	nextPlayerID      int
	nextPlayerIDMutex sync.Mutex

	messageQueue chan Message
}

func NewServer(host, port, connectionType string) *Server {
	return &Server{
		host:           host,
		port:           port,
		connectionType: connectionType,

		nextPlayerID: 1,

		messageQueue: make(chan Message, messageQueueSize),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.connectionType, s.host+":"+s.port)
	if err != nil {
		return err
	}
	fmt.Println("listening on " + s.host + ":" + s.port)

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("error accepting a connection on the listener:", err.Error())
				continue
			}

			s.nextPlayerIDMutex.Lock()
			id := s.nextPlayerID
			s.nextPlayerID++
			s.nextPlayerIDMutex.Unlock()

			acceptMessage := AcceptMessage{
				PlayerID: id,
			}
			bodyBytes, err := json.Marshal(acceptMessage)
			if err != nil {
				fmt.Println("error marshaling accept message:", err.Error())
				continue
			}

			encoder := json.NewEncoder(conn)
			message := Message{
				SenderID:    0,
				MessageType: MessageTypeAcceptConnection,
				Body:        bodyBytes,
			}

			err = encoder.Encode(message)
			if err != nil {
				fmt.Println("error sending accept message:", err.Error())
				continue
			}

			go s.handleRequest(conn)
		}
	}()

	return nil
}

func (s *Server) handleRequest(conn net.Conn) {
	for {
		decoder := json.NewDecoder(conn)

		message := Message{}
		err := decoder.Decode(&message)
		if err != nil {
			fmt.Println("error reading:", err.Error())
			continue
		}

		select {
		case s.messageQueue <- message:
		default:
			fmt.Println("message queue full")
		}
	}
}

func (s *Server) MessageQueue() chan Message {
	return s.messageQueue
}
