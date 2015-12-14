package entity

import (
	// "fmt"
	// "github.com/kkevinchou/ant/math/vector"
	"time"

	"github.com/kkevinchou/ant/physics"
	"github.com/kkevinchou/ant/render"
	"github.com/kkevinchou/ant/steering"
)

type Entity struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*render.RenderComponent
}

func New() *Entity {
	entity := &Entity{}

	physicsComponent := &physics.PhysicsComponent{
		MaxSpeed: 100,
		Mass:     10,
	}

	seekComponent := &steering.SeekComponent{}
	seekComponent.Initialize(entity)

	renderComponent := &render.RenderComponent{}
	renderComponent.Initialize("stag-head.png", entity)

	entity.PhysicsComponent = physicsComponent
	entity.SeekComponent = seekComponent
	entity.RenderComponent = renderComponent

	return entity
}

func (e *Entity) GetPhysicsComponent() *physics.PhysicsComponent {
	return e.PhysicsComponent
}

func (e *Entity) GetRenderComponent() *render.RenderComponent {
	return e.RenderComponent
}

func (e *Entity) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
