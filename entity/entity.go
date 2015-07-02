package entity

import (
	// "fmt"
	// "github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/physics"
	"github.com/kkevinchou/ant/render"
	"github.com/kkevinchou/ant/steering"
	"time"
)

type Entity struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*render.RenderComponent
}

func New() Entity {
	physicsComponent := &physics.PhysicsComponent{
		MaxSpeed: 200,
		Mass:     10,
	}

	seekComponent := &steering.SeekComponent{}
	seekComponent.Initialize(physicsComponent)

	renderComponent := &render.RenderComponent{}
	renderComponent.Initialize("stag-head.png", physicsComponent)

	entity := Entity{
		PhysicsComponent: physicsComponent,
		SeekComponent:    seekComponent,
		RenderComponent:  renderComponent,
	}
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
