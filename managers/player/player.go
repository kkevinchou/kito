package player

import (
	"github.com/kkevinchou/kito/lib/network"
)

type World interface {
	CommandFrame() int
}

type Player struct {
	ID     int
	Client *network.Client
}

type PlayerManager struct {
	players   []*Player
	playerMap map[int]*Player
	world     World
}

func NewPlayerManager(world World) *PlayerManager {
	return &PlayerManager{
		players:   []*Player{},
		playerMap: map[int]*Player{},
		world:     world,
	}
}

func (p *PlayerManager) RegisterPlayer(id int, client *network.Client) {
	player := &Player{ID: id, Client: client}
	p.playerMap[id] = player
	p.players = append(p.players, player)
}

func (p *PlayerManager) GetPlayer(id int) *Player {
	return p.playerMap[id]
}

func (p *PlayerManager) GetPlayers() []*Player {
	return p.players
}
