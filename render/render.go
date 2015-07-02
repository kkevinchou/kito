package render

import (
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
	// renderer.SetDrawColor(38, 3, 57, 255)
	// renderer.FillRect(&sdl.Rect{int32(position.X), int32(position.Y), 50, 50})
	renderer.Copy(texture, &sdl.Rect{0, 0, 64, 64}, &sdl.Rect{int32(position.X), int32(position.Y), 50, 50})
}

func (r *RenderComponent) Initialize(iconName string, p *physics.PhysicsComponent) {
	r.physicsComponent = p
	r.iconName = iconName
}

type RenderComposed interface {
	GetRenderComponent() *RenderComponent
}

type Renderable interface {
	RenderComposed
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

	for _, renderable := range r.renderables {
		renderable.GetRenderComponent().Render(r.assetManager, r.renderer)
	}

	r.renderer.Present()
}
