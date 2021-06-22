package player

import "net"

type Player struct {
	ID         int
	Connection net.Conn
}

type PlayerManager struct {
	players map[int]*Player
}

func (p *PlayerManager) RegisterPlayer(id int, connection net.Conn) {
	p.players[id] = &Player{ID: id, Connection: connection}
}

func (p *PlayerManager) GetPlayer(id int) *Player {
	return p.players[id]
}
