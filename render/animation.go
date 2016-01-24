package render

import (
	"time"

	"github.com/kkevinchou/ant/assets"
)

type AnimationState struct {
	currentFrame    int
	timeCounter     float64
	secondsPerFrame int

	MetaData assets.MetaData
}

func (a *AnimationState) GetFrame() int {
	return a.currentFrame
}

func (a *AnimationState) Update(delta time.Duration) {
	secondsPerFrame := 1 / float64(a.MetaData.Fps)
	a.timeCounter += delta.Seconds()
	for a.timeCounter >= secondsPerFrame {
		a.timeCounter -= secondsPerFrame
		a.currentFrame += 1
	}
	a.currentFrame = a.currentFrame % a.MetaData.NumFrames
}
