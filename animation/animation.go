package animation

import "time"

type AnimationState struct {
	currentFrame    int
	timeCounter     float64
	secondsPerFrame int
	NumFrames       int
	Fps             int
	Name            string
}

func (a *AnimationState) GetFrame() int {
	return a.currentFrame
}

func (a *AnimationState) Update(delta time.Duration) {
	secondsPerFrame := 1 / float64(a.Fps)
	a.timeCounter += delta.Seconds()
	for a.timeCounter >= secondsPerFrame {
		a.timeCounter -= secondsPerFrame
		a.currentFrame += 1
	}
	a.currentFrame = a.currentFrame % a.NumFrames
}
