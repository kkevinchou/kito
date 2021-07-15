package player

import (
	"net"

	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/settings"
)

type World interface {
	CommandFrame() int
}

type Player struct {
	ID     int
	Client *network.Client
}

type PlayerManager struct {
	players map[int]*Player
	world   World
}

func NewPlayerManager(world World) *PlayerManager {
	return &PlayerManager{
		players: map[int]*Player{},
		world:   world,
	}
}

func (p *PlayerManager) RegisterPlayerWithConnection(id int, connection net.Conn) {
	client := network.NewClient(settings.ServerID, connection)
	client.SetCommandFrameFunction(p.world.CommandFrame)
	p.RegisterPlayerWithClient(id, client)
}

func (p *PlayerManager) RegisterPlayerWithClient(id int, client *network.Client) {
	p.players[id] = &Player{ID: id, Client: client}
}

func (p *PlayerManager) GetPlayer(id int) *Player {
	return p.players[id]
}

func (p *PlayerManager) GetPlayers() map[int]*Player {
	return p.players
}
