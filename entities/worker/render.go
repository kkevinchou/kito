package worker

import "github.com/kkevinchou/ant/lib"

type RenderComponent struct {
	entity         Worker
	animationState *lib.AnimationState
}

func (r *RenderComponent) Texture() string {
	return "worker"
}
