package pathing

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type NavMeshRenderData struct {
	ID      string
	Visible bool
}

func (n *NavMeshRenderData) IsVisible() bool {
	return true
}

type RenderComponent struct {
	RenderData *NavMeshRenderData
}

func (r *RenderComponent) GetRenderData() components.RenderData {
	return r.RenderData
}

func (n *RenderComponent) Position() vector.Vector3 {
	return vector.Vector3{}
}

func (n *RenderComponent) SetPosition(v vector.Vector3) {
}
