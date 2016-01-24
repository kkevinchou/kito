package grass

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/render"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
}

type RenderComponent struct {
	entity         Renderable
	animationState render.AnimationState
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	texture := assetManager.GetAnimation(r.animationState.MetaData.Name, r.animationState.GetFrame())
	renderer.Copy(texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 32, 64, 64})
}

func (r *RenderComponent) UpdateRenderComponent(delta time.Duration) {
	r.animationState.Update(delta)
}
