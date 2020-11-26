package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/systems/render"
)

type CommandPoller func(game *Game) []Command

type Command interface {
	Execute(game *Game)
}

type SetCameraSpeed struct {
	X float64
	Y float64
	Z float64
}

func (c *SetCameraSpeed) Execute(game *Game) {
	game.SetCameraCommandHeading(vector.Vector3{X: c.X, Y: c.Y, Z: c.Z})
}

type QuitCommand struct {
}

func (c *QuitCommand) Execute(game *Game) {
	game.GameOver()
}

type CameraViewCommand struct {
	X float64
	Y float64
}

func (c *CameraViewCommand) Execute(game *Game) {
	if !game.camera.controlled {
		return
	}
	game.CameraViewChange(vector.Vector{X: c.X, Y: c.Y})
}

type CameraRaycastCommand struct {
	X float64
	Y float64
}

func (c *CameraRaycastCommand) Execute(game *Game) {
	renderSystem := directory.GetDirectory().RenderSystem()
	worldPoint := renderSystem.GetWorldPoint(c.X, c.Y)
	dir := worldPoint.Sub(game.camera.Position()).Normalize()
	render.LineStart = game.camera.Position()
	render.LineEnd = game.camera.Position().Add(dir.Scale(3))
	fmt.Println(worldPoint)
	fmt.Println("Camera position:", game.camera.Position(), "Direction:", dir)
}

type SetCameraControlCommand struct {
	Value bool
}

func (c *SetCameraControlCommand) Execute(game *Game) {
	game.camera.controlled = c.Value
	// sdl.SetRelativeMouseMode(c.Value)
}
