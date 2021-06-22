package networklistener

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/systems/base"
)

type World interface{}

type NetworkListenerSystem struct {
	*base.BaseSystem

	world   World
	nserver *network.Server
}

func NewNetworkListenerSystem(world World, host, port, connectionType string) *NetworkListenerSystem {
	nserver := network.NewServer(host, port, connectionType)
	nserver.Start()

	return &NetworkListenerSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		nserver:    nserver,
	}
}

func (s *NetworkListenerSystem) RegisterEntity(entity entities.Entity) {
}

func (s *NetworkListenerSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	incomingConnections := s.nserver.PullIncomingConnections()
	for _, incomingConnection := range incomingConnections {
		fmt.Println("New player connected with id", incomingConnection.PlayerID)
		playerManager.RegisterPlayer(incomingConnection.PlayerID, incomingConnection.Connection)
	}
}
