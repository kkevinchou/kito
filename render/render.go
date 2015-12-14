package render

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/physics"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderComponent struct {
	entity   physics.PhysicsI
	iconName string
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	physicsComponent := r.entity.GetPhysicsComponent()

	position := physicsComponent.Position
	texture := assetManager.GetTexture(r.iconName)
	renderer.Copy(texture, &sdl.Rect{0, 0, 64, 64}, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 32, 64, 64})

	heading := physicsComponent.Velocity.Normalize()
	lineStart := heading.Scale(40).Add(position)
	lineEnd := heading.Scale(55).Add(position)

	renderer.SetDrawColor(0, 255, 255, 255)
	renderer.DrawLine(int(lineStart.X), int(lineStart.Y), int(lineEnd.X), int(lineEnd.Y))
}

func (r *RenderComponent) Initialize(iconName string, entity physics.PhysicsI) {
	r.entity = entity
	r.iconName = iconName
}

type RenderI interface {
	GetRenderComponent() *RenderComponent
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	assetManager *assets.Manager
	renderables  []RenderI
}

func NewRenderSystem(renderer *sdl.Renderer, assetManager *assets.Manager) RenderSystem {
	renderSystem := RenderSystem{
		renderer:     renderer,
		assetManager: assetManager,
	}

	return renderSystem
}

func (r *RenderSystem) Register(renderable RenderI) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) Update(delta time.Duration) {
	r.renderer.Clear()
	r.renderer.SetDrawColor(151, 117, 170, 255)
	r.renderer.FillRect(&sdl.Rect{0, 0, 800, 600})

	for _, renderable := range r.renderables {
		renderable.GetRenderComponent().Render(r.assetManager, r.renderer)
	}

	r.renderer.Present()
}
