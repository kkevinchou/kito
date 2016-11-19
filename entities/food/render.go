package food

import (
	"fmt"
	"time"

	"github.com/kkevinchou/ant/lib"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderComponent struct {
	entity  Food
	texture *sdl.Texture
}

const (
	textWidth  = 64
	textHeight = 28
)

func (r *RenderComponent) Render(assetManager *lib.AssetManager, renderer *sdl.Renderer) {
	if r.entity.Owned() {
		return
	}

	position := r.entity.Position()
	// renderer.Copy(r.texture, nil, &sdl.Rect{int32(position.X) - 32, int32(position.Y) - 64, 64, 64})

	font := assetManager.GetFont("courier_new.ttf")
	renderer.SetDrawColor(17, 72, 0, 255)
	surface, err := font.RenderUTF8_Solid("FOOD", sdl.Color{100, 100, 0, 100})
	if err != nil {
		fmt.Println(err)
		return
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		texture.Destroy()
	}()

	_, _, w, h, err := texture.Query()
	if err != nil {
		fmt.Println(err)
		return
	}

	renderer.Copy(texture, nil, &sdl.Rect{int32(position.X) - w/2, int32(position.Y) - h, w, h})
}

func (r *RenderComponent) UpdateRenderComponent(delta time.Duration) {
}

func (r *RenderComponent) GetRenderPriority() int {
	return 1000
}

func (r *RenderComponent) GetY() float64 {
	return r.entity.Position().Y
}
