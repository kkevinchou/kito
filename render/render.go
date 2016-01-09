package render

import (
	"fmt"
	"log"
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Renderable interface {
	Render(*assets.Manager, *sdl.Renderer)
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	assetManager *assets.Manager
	renderables  []Renderable
}

func initFont() *ttf.Font {
	ttf.Init()

	font, err := ttf.OpenFont("assets/fonts/courier_new.ttf", 30)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Font not found")
	}

	return font

	// surface, err := font.RenderUTF8_Solid("test text abcdefghijklmnopqrstuvwxyz", sdl.Color{100, 100, 100, 100})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// texture, err := renderer.CreateTextureFromSurface(surface)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// surface.Free()
	// font.Close()

	// err = renderer.Copy(texture, nil, &sdl.Rect{0, 0, 648, 35})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

func NewRenderSystem(renderer *sdl.Renderer, assetManager *assets.Manager) RenderSystem {
	renderSystem := RenderSystem{
		renderer:     renderer,
		assetManager: assetManager,
	}

	// initFont(renderer)

	return renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
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
