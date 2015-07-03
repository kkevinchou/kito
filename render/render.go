package render

import (
	// "fmt"
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/physics"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type RenderComponent struct {
	physicsComponent *physics.PhysicsComponent
	iconName         string
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	position := r.physicsComponent.Position
	texture := assetManager.GetTexture(r.iconName)
	renderer.Copy(texture, &sdl.Rect{0, 0, 64, 64}, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 32, 64, 64})

	heading := r.physicsComponent.Velocity.Normalize()
	lineStart := heading.Scale(40).Add(position)
	lineEnd := heading.Scale(55).Add(position)

	renderer.SetDrawColor(0, 255, 255, 255)
	renderer.DrawLine(int(lineStart.X), int(lineStart.Y), int(lineEnd.X), int(lineEnd.Y))
}

func (r *RenderComponent) Initialize(iconName string, p *physics.PhysicsComponent) {
	r.physicsComponent = p
	r.iconName = iconName
}

type Renderable interface {
	GetRenderComponent() *RenderComponent
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	assetManager *assets.Manager
	renderables  []Renderable
}

func NewRenderSystem(renderer *sdl.Renderer, assetManager *assets.Manager) RenderSystem {
	renderSystem := RenderSystem{
		renderer:     renderer,
		assetManager: assetManager,
	}

	return renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) Update(delta time.Duration) {
	r.renderer.SetDrawColor(151, 117, 170, 255)
	r.renderer.FillRect(&sdl.Rect{0, 0, 800, 600})
	// r.renderer.FillRect(&sdl.Rect{1, 1, 1, 1})

	for _, renderable := range r.renderables {
		renderable.GetRenderComponent().Render(r.assetManager, r.renderer)
	}

	r.renderer.Present()
}
