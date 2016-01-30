package pathing

import (
	"fmt"
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderComponent struct {
	polygons []*Polygon
}

func (r *RenderComponent) UpdateRenderComponent(delta time.Duration) {}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	font := assetManager.GetFont("courier_new.ttf")
	for _, poly := range r.polygons {
		renderer.SetDrawColor(17, 72, 0, 255)
		// renderer.SetDrawColor(17, 72, 0, 112)
		points := []sdl.Point{}
		for _, point := range poly.Points() {
			points = append(points, sdl.Point{X: int32(point.X), Y: int32(point.Y)})

			surface, err := font.RenderUTF8_Solid("test text abcdefghijklmnopqrstuvwxyz", sdl.Color{100, 100, 100, 100})
			if err != nil {
				fmt.Println(err)
				return
			}

			_ = surface

			// texture, err := renderer.CreateTextureFromSurface(surface)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// surface.Free()

			// _, _, width, height, err := texture.Query()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// err = renderer.Copy(texture, nil, &sdl.Rect{int32(point.X), int32(point.Y), width, height})
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
		}

		points = append(points, points[0])
		renderer.DrawLines(points)
	}
}

func (r *RenderComponent) GetRenderPriority() int {
	return 1
}

func (r *RenderComponent) GetY() float64 {
	return -1
}
