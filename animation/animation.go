package animation

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/veandco/go-sdl2/sdl"
)

type AnimationState struct {
	animationDefinition *assets.AnimationDefinition
	currentFrame        int
	timeCounter         float64
	secondsPerFrame     float64
}

func (a *AnimationState) GetFrame() *sdl.Texture {
	return a.animationDefinition.GetFrame(a.currentFrame)
}

func (a *AnimationState) Update(delta time.Duration) {
	a.timeCounter += delta.Seconds()
	for a.timeCounter >= a.secondsPerFrame {
		a.timeCounter -= a.secondsPerFrame
		a.currentFrame += 1
	}
	a.currentFrame = a.currentFrame % a.animationDefinition.NumFrames()
}

func CreateStateFromAnimationDef(animationDefinition *assets.AnimationDefinition) *AnimationState {
	secondsPerFrame := 1 / float64(animationDefinition.Fps())
	return &AnimationState{
		animationDefinition: animationDefinition,
		secondsPerFrame:     secondsPerFrame,
	}
}
