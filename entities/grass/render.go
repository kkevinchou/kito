package grass

import (
	"time"

	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
}

type RenderComponent struct {
	entity         Renderable
	animationState *animation.AnimationState
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	texture := r.animationState.GetCurrentFrame()
	renderer.Copy(texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 64, 64, 64})
}

func (r *RenderComponent) UpdateRenderComponent(delta time.Duration) {
	r.animationState.Update(delta)
}

func (r *RenderComponent) GetRenderPriority() int {
	return 1000
}

func (r *RenderComponent) GetY() float64 {
	return r.entity.Position().Y
}
