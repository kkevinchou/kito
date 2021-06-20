package components

import (
	"github.com/kkevinchou/kito/lib/network"
)

type ConnectionComponent struct {
	Client   *network.Client
	PlayerID int
}

func NewConnectionComponent(host, port, connectionType string) (*ConnectionComponent, error) {
	client := network.NewClient()
	acceptMessage, err := client.Connect(host, port, connectionType)
	if err != nil {
		return nil, err
	}

	return &ConnectionComponent{
		Client:   client,
		PlayerID: acceptMessage.PlayerID,
	}, nil
}

func (c *ConnectionComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ConnectionComponent = c
}
