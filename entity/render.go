package entity

import (
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/render"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
	Velocity() vector.Vector
	Heading() vector.Vector
}

type RenderComponent struct {
	entity         Renderable
	animationState render.AnimationState
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	texture := assetManager.GetAnimation(r.animationState.MetaData.Name, 0)
	renderer.Copy(texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 32, 64, 64})

	// heading := r.entity.Heading().Normalize()
	heading := r.entity.Velocity().Normalize()
	lineStart := heading.Scale(40).Add(position)
	lineEnd := heading.Scale(55).Add(position)

	renderer.SetDrawColor(0, 255, 255, 255)
	renderer.DrawLine(int(lineStart.X), int(lineStart.Y), int(lineEnd.X), int(lineEnd.Y))
}
