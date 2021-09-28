package networkupdate

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/network"
)

const (
	commandFramesPerUpdate = 10
)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEventBroker() eventbroker.EventBroker
	GetSingleton() *singleton.Singleton
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
			fmt.Println("failed to serialize event", err)
			continue
		}

		serializedEvent := network.Event{Type: int(event.Type()), Bytes: bytes}
		gameStateUpdate.Events = append(gameStateUpdate.Events, serializedEvent)
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		gameStateUpdate.LastInputCommandFrame = player.LastInputCommandFrame
		gameStateUpdate.LastInputGlobalCommandFrame = player.LastInputGlobalCommandFrame
		gameStateUpdate.CurrentGlobalCommandFrame = s.world.GetSingleton().CommandFrame
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
		Type:        int(entity.Type()),
		Position:    componentContainer.TransformComponent.Position,
		Orientation: componentContainer.TransformComponent.Orientation,
	}
}
