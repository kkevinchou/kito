package ability

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) (entities.Entity, error)
	RegisterEntities([]entities.Entity)
}

type AbilitySystem struct {
	*base.BaseSystem
	world World
}

func NewAbilitySystem(world World) *AbilitySystem {
	return &AbilitySystem{
		world: world,
	}
}

func (s *AbilitySystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerManager := directory.GetDirectory().PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		playerInput := singleton.PlayerInput[player.ID]
		entity, err := s.world.GetEntityByID(player.EntityID)
		if err != nil {
			fmt.Printf("ability system couldn't find player %d\n", player.ID)
			continue
		}
		if entity != nil {
			if key, ok := playerInput.KeyboardInput[input.KeyboardKeyQ]; ok && key.Event == input.KeyboardEventUp {
				projSpeed := 50
				cc := entity.GetComponentContainer()
				direction := cc.TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
				proj := entities.NewProjectile(cc.TransformComponent.Position.Add(mgl64.Vec3{0, 10, 0}))
				projcc := proj.GetComponentContainer()
				projcc.PhysicsComponent.Velocity = direction.Mul(float64(projSpeed))
				s.world.RegisterEntities([]entities.Entity{proj})
			}
		}
	}
}
