package food

import "github.com/veandco/go-sdl2/sdl"

type RenderComponent struct {
	entity  Food
	texture *sdl.Texture
}

const (
	textWidth  = 64
	textHeight = 28
)

func (r *RenderComponent) Texture() string {
	return "mushroom-gills"
}
