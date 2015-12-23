package entity

import (
	// "fmt"
	// "github.com/kkevinchou/ant/math/vector"
	"time"

	"github.com/kkevinchou/ant/physics"
	"github.com/kkevinchou/ant/steering"
)

type Entity struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
}

func New() *Entity {
	entity := &Entity{}

	entity.PhysicsComponent = &physics.PhysicsComponent{MaxSpeed: 100, Mass: 10}
	entity.SeekComponent = &steering.SeekComponent{Entity: entity}
	entity.RenderComponent = &RenderComponent{entity: entity, iconName: "stag-head.png"}

	return entity
}

func (e *Entity) GetPhysicsComponent() *physics.PhysicsComponent {
	return e.PhysicsComponent
}

func (e *Entity) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
