package networkupdate

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/events"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/systems/base"
)

const (
	commandFramesPerUpdate = 20
)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEventBroker() eventbroker.EventBroker
}

type NetworkUpdateSystem struct {
	*base.BaseSystem
	world         World
	entities      []entities.Entity
	elapsedFrames int
	events        []events.Event
}

func NewNetworkUpdateSystem(world World) *NetworkUpdateSystem {
	networkUpdateSystem := &NetworkUpdateSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}

	eventBroker := world.GetEventBroker()
	eventBroker.AddObserver(networkUpdateSystem, []events.EventType{
		events.EventTypeCreateEntity,
	})

	return networkUpdateSystem
}

func (s *NetworkUpdateSystem) Observe(event events.Event) {
	if event.Type() == events.EventTypeCreateEntity {
		s.events = append(s.events, event)
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

	// TODO: perhaps it's better to have a stack of events that we pop off so we don't accidentally lose
	// events?
	defer s.clearEvents()

	s.elapsedFrames %= commandFramesPerUpdate

	gameStateUpdate := &network.GameStateUpdateMessage{
		Entities: map[int]network.EntitySnapshot{},
	}

	for _, entity := range s.entities {
		gameStateUpdate.Entities[entity.GetID()] = constructEntitySnapshot(entity)
	}

	for _, event := range s.events {
		bytes, err := event.Serialize()
		if err != nil {
			fmt.Println("failed to serial event", err)
			continue
		}

		serializedEvent := network.Event{Type: int(event.Type()), Bytes: bytes}
		gameStateUpdate.Events = append(gameStateUpdate.Events, serializedEvent)
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		player.Client.SendMessage(network.MessageTypeGameStateUpdate, gameStateUpdate)
	}
}

func (s *NetworkUpdateSystem) clearEvents() {
	s.events = []events.Event{}
}

func constructEntitySnapshot(entity entities.Entity) network.EntitySnapshot {
	componentContainer := entity.GetComponentContainer()

	return network.EntitySnapshot{
		ID:          entity.GetID(),
		Position:    componentContainer.TransformComponent.Position,
		Orientation: componentContainer.TransformComponent.Orientation,
	}
}
