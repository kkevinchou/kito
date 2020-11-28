package commands

import (
	"github.com/kkevinchou/kito/lib/math/vector"
)

type Command interface{}

// func (c *CameraRaycastCommand) Execute(game *Game) {
// 	renderSystem := directory.GetDirectory().RenderSystem()
// 	worldPoint := renderSystem.GetWorldPoint(c.X, c.Y)
// 	dir := worldPoint.Sub(game.camera.Position()).Normalize()
// 	render.LineStart = game.camera.Position()
// 	render.LineEnd = game.camera.Position().Add(dir.Scale(3))
// 	fmt.Println(worldPoint)
// 	fmt.Println("Camera position:", game.camera.Position(), "Direction:", dir)
// }

type MoveCommand struct {
	Value vector.Vector3
}

type UpdateViewCommand struct {
	Value vector.Vector
}

type ToggleCameraControlCommand struct {
	Value bool
}

type QuitCommand struct{}
