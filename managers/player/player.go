package player

import (
	"net"

	"github.com/kkevinchou/kito/lib/network"
)

type Player struct {
	ID     int
	Client *network.Client
}

type PlayerManager struct {
	players map[int]*Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: map[int]*Player{},
	}
}

func (p *PlayerManager) RegisterPlayerWithConnection(id int, connection net.Conn) {
	client := network.NewClientFromConnection(connection)
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
