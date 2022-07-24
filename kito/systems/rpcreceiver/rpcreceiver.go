package rpcreceiver

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	GetEventBroker() eventbroker.EventBroker
	GetEntityByID(id int) entities.Entity
	GetPlayerEntityByID(id int) entities.Entity
}

type RPCReceiverSystem struct {
	*base.BaseSystem
	world  World
	events []events.Event
}

func NewRPCReceiverSystem(world World) *RPCReceiverSystem {
	rpcSystem := &RPCReceiverSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}

	eventBroker := world.GetEventBroker()
	eventBroker.AddObserver(rpcSystem, []events.EventType{
		events.EventTypeRPC,
	})
	return rpcSystem
}

func (s *RPCReceiverSystem) Observe(event events.Event) {
	if event.Type() == events.EventTypeRPC {
		s.events = append(s.events, event)
	}
}

func (s *RPCReceiverSystem) clearEvents() {
	s.events = nil
}
func (s *RPCReceiverSystem) Update(delta time.Duration) {
	defer s.clearEvents()

	for _, event := range s.events {
		if e, ok := event.(*events.RPCEvent); ok {
			tokens := strings.Split(e.Command, " ")
			if len(tokens) == 0 {
				continue
			}

			if len(tokens) != 3 {
				continue
			}

			command := tokens[0]
			if command == "position" {
				var entity entities.Entity
				if tokens[1] == "self" || tokens[1] == "me" {
					entity = s.world.GetPlayerEntityByID(e.PlayerID)
				} else {
					entityID, err := strconv.Atoi(tokens[1])
					if err != nil {
						continue
					}
					entity = s.world.GetEntityByID(entityID)
				}

				if entity == nil {
					continue
				}

				vec := strings.Split(tokens[2], ",")
				x, err := strconv.Atoi(vec[0])
				if err != nil {
					continue
				}
				y, err := strconv.Atoi(vec[1])
				if err != nil {
					continue
				}
				z, err := strconv.Atoi(vec[2])
				if err != nil {
					continue
				}

				cc := entity.GetComponentContainer()

				positionVec := mgl64.Vec3{float64(x), float64(y), float64(z)}
				cc.TransformComponent.Position = positionVec
				if cc.ThirdPersonControllerComponent != nil {
					cc.ThirdPersonControllerComponent.BaseVelocity = mgl64.Vec3{}
				}

				fmt.Println("executed rpc", e.Command)
			}
		}
	}
}

func (s *RPCReceiverSystem) Name() string {
	return "RPCSystem"
}