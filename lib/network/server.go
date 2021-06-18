package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	host           string
	port           string
	connectionType string

	nextPlayerID      int
	nextPlayerIDMutex sync.Mutex

	messageQueue chan *Message
	connections  map[int]net.Conn
}

func NewServer(host, port, connectionType string) *Server {
	return &Server{
		host:           host,
		port:           port,
		connectionType: connectionType,

		nextPlayerID: 1,

		messageQueue: make(chan *Message, messageQueueBufferSize),
		connections:  map[int]net.Conn{},
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

			id := s.generateNextPlayerID()
			s.connections[id] = conn

			message, err := s.createAcceptMessage(id)
			if err != nil {
				fmt.Println(err)
				continue
			}

			s.SendMessage(id, message)
			if err != nil {
				fmt.Println("error sending accept message:", err.Error())
				continue
			}

			go queueIncomingMessages(conn, s.messageQueue)
		}
	}()

	return nil
}

func (s *Server) IncomingMessageQueue() chan *Message {
	return s.messageQueue
}

func (s *Server) SendMessage(playerID int, message *Message) error {
	encoder := json.NewEncoder(s.connections[playerID])
	err := encoder.Encode(message)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) generateNextPlayerID() int {
	s.nextPlayerIDMutex.Lock()
	id := s.nextPlayerID
	s.nextPlayerID++
	s.nextPlayerIDMutex.Unlock()

	return id
}

func (s *Server) createAcceptMessage(playerID int) (*Message, error) {
	acceptMessage := AcceptMessage{
		PlayerID: playerID,
	}
	bodyBytes, err := json.Marshal(acceptMessage)
	if err != nil {
		fmt.Println("error marshaling accept message:", err.Error())
		return nil, err
	}
	return &Message{
		SenderID:    0,
		MessageType: MessageTypeAcceptConnection,
		Body:        bodyBytes,
	}, nil

}
