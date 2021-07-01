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

func (p *PlayerManager) RegisterPlayer(id int, connection net.Conn) {
	p.players[id] = &Player{ID: id, Client: network.NewClientFromConnection(connection)}
}

func (p *PlayerManager) GetPlayer(id int) *Player {
	return p.players[id]
}

func (p *PlayerManager) GetPlayers() map[int]*Player {
	return p.players
}
