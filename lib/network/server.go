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

	incomingConnections chan *Connection
}

func NewServer(host, port, connectionType string) *Server {
	return &Server{
		host:           host,
		port:           port,
		connectionType: connectionType,

		nextPlayerID: 70000,

		incomingConnections: make(chan *Connection, incomingConnectionsBufferSize),
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

			playerID := s.generateNextPlayerID()

			message, err := s.createAcceptMessage(playerID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			s.SendMessage(conn, message)
			if err != nil {
				fmt.Println("error sending accept message:", err.Error())
				continue
			}

			select {
			case s.incomingConnections <- &Connection{Connection: conn, PlayerID: playerID}:
			default:
				panic("incomingConnections queue full")
			}
		}
	}()

	return nil
}

func (s *Server) PullIncomingConnections() []*Connection {
	connections := []*Connection{}

	for i := 0; i < incomingConnectionsBufferSize; i++ {
		select {
		case connection := <-s.incomingConnections:
			connections = append(connections, connection)
		default:
			return connections
		}
	}

	return connections
}

func (s *Server) SendMessage(connection net.Conn, message *Message) error {
	encoder := json.NewEncoder(connection)
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
