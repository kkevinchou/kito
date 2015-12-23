package render

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderI interface {
	Render(*assets.Manager, *sdl.Renderer)
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
		renderable.Render(r.assetManager, r.renderer)
	}

	r.renderer.Present()
}
