package entity

import (
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Position() vector.Vector
	Velocity() vector.Vector
}

type RenderComponent struct {
	entity   Renderable
	iconName string
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	position := r.entity.Position()
	texture := assetManager.GetTexture(r.iconName)
	renderer.Copy(texture, &sdl.Rect{0, 0, 64, 64}, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 32, 64, 64})

	heading := r.entity.Velocity().Normalize()
	lineStart := heading.Scale(40).Add(position)
	lineEnd := heading.Scale(55).Add(position)

	renderer.SetDrawColor(0, 255, 255, 255)
	renderer.DrawLine(int(lineStart.X), int(lineStart.Y), int(lineEnd.X), int(lineEnd.Y))
}
