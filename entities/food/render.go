package food

import (
	"time"

	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
}

type RenderComponent struct {
	entity  Renderable
	texture *sdl.Texture
}

func (r *RenderComponent) Render(assetManager *lib.AssetManager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	renderer.Copy(r.texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 64, 64, 64})
}

func (r *RenderComponent) UpdateRenderComponent(delta time.Duration) {
}

func (r *RenderComponent) GetRenderPriority() int {
	return 1000
}

func (r *RenderComponent) GetY() float64 {
	return r.entity.Position().Y
}
