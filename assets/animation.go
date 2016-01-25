package assets

import "github.com/veandco/go-sdl2/sdl"

type AnimationDefinition struct {
	numFrames int
	fps       int
	frames    []*sdl.Texture
}

func (a *AnimationDefinition) NumFrames() int {
	return a.numFrames
}

func (a *AnimationDefinition) Fps() int {
	return a.fps
}

func (a *AnimationDefinition) GetFrame(frame int) *sdl.Texture {
	return a.frames[frame]
}
