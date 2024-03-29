package networkupdate

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils/entityutils"
	"github.com/kkevinchou/kito/lib/metrics"
)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEventBroker() eventbroker.EventBroker
	GetSingleton() *singleton.Singleton
	CommandFrame() int
	QueryEntity(componentFlags int) []entities.Entity
	MetricsRegistry() *metrics.MetricsRegistry
}

type NetworkUpdateSystem struct {
	*base.BaseSystem
	world         World
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
		events.EventTypeUnregisterEntity,
	})

	return networkUpdateSystem
}

func (s *NetworkUpdateSystem) Observe(event events.Event) {
	if event.Type() == events.EventTypeUnregisterEntity {
		s.events = append(s.events, event)
	}
}

func (s *NetworkUpdateSystem) Update(delta time.Duration) {
	s.elapsedFrames++
	if s.elapsedFrames < settings.CommandFramesPerServerUpdate {
		return
	}

	s.elapsedFrames %= settings.CommandFramesPerServerUpdate

	serverStats := map[string]string{
		"fps":       fmt.Sprintf("%d", int(s.world.MetricsRegistry().GetOneSecondSum("fps"))),
		"frametime": fmt.Sprintf("%d", int(s.world.MetricsRegistry().GetOneSecondAverage("frametime"))),
	}

	gameStateUpdate := &knetwork.GameStateUpdateMessage{
		Entities:    map[int]knetwork.EntitySnapshot{},
		ServerStats: serverStats,
	}

	for _, entity := range s.world.QueryEntity(components.ComponentFlagTransform | components.ComponentFlagNetwork) {
		if entity.Type() == types.EntityTypeCamera {
			continue
		}
		gameStateUpdate.Entities[entity.GetID()] = entityutils.ConstructEntitySnapshot(entity)
	}

	defer s.clearEvents()
	for _, event := range s.events {
		bytes, err := knetwork.Serialize(event)
		if err != nil {
			fmt.Println("failed to serialize event", err)
			continue
		}

		networkEvent := knetwork.Event{Type: event.Type(), Bytes: bytes}
		gameStateUpdate.Events = append(gameStateUpdate.Events, networkEvent)
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		gameStateUpdate.LastInputCommandFrame = player.LastInputLocalCommandFrame
		gameStateUpdate.LastInputGlobalCommandFrame = player.LastInputGlobalCommandFrame
		gameStateUpdate.CurrentGlobalCommandFrame = s.world.CommandFrame()
		player.Client.SendMessage(knetwork.MessageTypeGameStateUpdate, gameStateUpdate)
	}
}

func (s *NetworkUpdateSystem) clearEvents() {
	s.events = nil
}

func (s *NetworkUpdateSystem) Name() string {
	return "NetworkUpdateSystem"
}
