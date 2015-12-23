package pathing

import (
	"github.com/kkevinchou/ant/assets"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderComponent struct {
	polygons []*Polygon
}

func (r *RenderComponent) Render(assetManager *assets.Manager, renderer *sdl.Renderer) {
	for _, poly := range r.polygons {
		renderer.SetDrawColor(17, 72, 0, 255)
		// renderer.SetDrawColor(17, 72, 0, 112)
		points := []sdl.Point{}
		for _, point := range poly.Points() {
			points = append(points, sdl.Point{X: int32(point.X), Y: int32(point.Y)})
		}
		points = append(points, points[0])
		renderer.DrawLines(points)
	}
}
