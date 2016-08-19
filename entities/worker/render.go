package worker

import (
	"time"

	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
	Velocity() vector.Vector
	Heading() vector.Vector
}

type RenderComponent struct {
	entity         Renderable
	animationState *lib.AnimationState
}

func (r *RenderComponent) Render(assetManager *lib.AssetManager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	texture := r.animationState.GetCurrentFrame()
	renderer.Copy(texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 64, 64, 64})

	heading := r.entity.Heading().Normalize()
	// heading := r.entity.Velocity().Normalize()
	lineStart := heading.Scale(40).Add(position).Add(vector.Vector{0, -32})
	lineEnd := heading.Scale(55).Add(position).Add(vector.Vector{0, -32})

	renderer.SetDrawColor(0, 255, 255, 255)
	renderer.DrawLine(int(lineStart.X), int(lineStart.Y), int(lineEnd.X), int(lineEnd.Y))
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
