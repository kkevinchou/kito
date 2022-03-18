package player

import (
	"github.com/kkevinchou/kito/kito/types"
)

type World interface {
	CommandFrame() int
}

// TODO: make this an entity with components
type Player struct {
	ID     int
	Client types.NetworkClient

	LastInputLocalCommandFrame  int // the player's last command frame
	LastInputGlobalCommandFrame int // the gcf when this input was received
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

func (p *PlayerManager) RegisterPlayer(id int, client types.NetworkClient) {
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
