package grass

type RenderComponent struct {
	entity Grass
}

func (r *RenderComponent) Texture() string {
	return "high-grass"
}

func (r *RenderComponent) Visible() bool {
	return true
}
