package networkupdate

import (
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/systems/base"
)

const (
	commandFramesPerUpdate = 20
)

type World interface {
	RegisterEntities([]entities.Entity)
}

type NetworkUpdateSystem struct {
	*base.BaseSystem
	world         World
	entities      []entities.Entity
	elapsedFrames int
}

func NewNetworkUpdateSystem(world World) *NetworkUpdateSystem {
	return &NetworkUpdateSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *NetworkUpdateSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.TransformComponent != nil && componentContainer.NetworkComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *NetworkUpdateSystem) Update(delta time.Duration) {
	s.elapsedFrames++
	if s.elapsedFrames < commandFramesPerUpdate {
		return
	}

	s.elapsedFrames %= commandFramesPerUpdate

	snapshot := &network.GameStateSnapshotMessage{
		Entities: map[int]network.EntitySnapshot{},
	}

	for _, entity := range s.entities {
		snapshot.Entities[entity.GetID()] = constructEntitySnapshot(entity)
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		player.Client.SendMessage(network.MessageTypeGameStateSnapshot, snapshot)
	}
}

func constructEntitySnapshot(entity entities.Entity) network.EntitySnapshot {
	componentContainer := entity.GetComponentContainer()

	return network.EntitySnapshot{
		ID:          entity.GetID(),
		Position:    componentContainer.TransformComponent.Position,
		Orientation: componentContainer.TransformComponent.Orientation,
	}
}
