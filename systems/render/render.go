package render

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/kkevinchou/ant/lib"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Renderable interface {
	Render(*lib.AssetManager, *sdl.Renderer)
	UpdateRenderComponent(time.Duration)
	GetRenderPriority() int
	GetY() float64
}

type Renderables []Renderable

func (r Renderables) Len() int {
	return len(r)
}

func (r Renderables) Less(i, j int) bool {
	if r[i].GetRenderPriority() == r[j].GetRenderPriority() {
		return r[i].GetY() < r[j].GetY()
	}
	return r[i].GetRenderPriority() < r[j].GetRenderPriority()
}

func (r Renderables) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	assetManager *lib.AssetManager
	renderables  Renderables
}

func initFont() *ttf.Font {
	ttf.Init()

	font, err := ttf.OpenFont("_assets/fonts/courier_new.ttf", 30)
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

func NewRenderSystem(renderer *sdl.Renderer, assetManager *lib.AssetManager) *RenderSystem {
	renderSystem := RenderSystem{
		renderer:     renderer,
		assetManager: assetManager,
	}

	_ = initFont()

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) Update(delta time.Duration) {
	r.renderer.Clear()
	r.renderer.SetDrawColor(151, 117, 170, 255)
	r.renderer.FillRect(&sdl.Rect{0, 0, 800, 600})
	sort.Stable(r.renderables)

	// TODO: have the render system know how to render as opposed to the render component
	// the component should provide the data necessary for rendering
	for _, renderable := range r.renderables {
		renderable.UpdateRenderComponent(delta)
		renderable.Render(r.assetManager, r.renderer)
	}

	r.renderer.Present()
}

// func (r *RenderSystem) EventHandlers() []systems.EventHandler {
// 	return []systems.EventHandler{
// 		systems.EventHandler{
// 			Type:    systems.EntityCreated,
// 			Handler: r.HandleEntityCreated,
// 		},
// 	}
// }

// func (r *RenderSystem) HandleEntityCreated(event systems.Event) {
// }
