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

type UpdateViewCommand struct {
	Value vector.Vector
}

type ToggleCameraControlCommand struct {
	Value bool
}

type QuitCommand struct{}

type KeyboardKey string
type KeyboardEvent int

const (
	KeyboardKeyW KeyboardKey = "W"
	KeyboardKeyA KeyboardKey = "A"
	KeyboardKeyS KeyboardKey = "S"
	KeyboardKeyD KeyboardKey = "D"

	KeyboardKeyLShift KeyboardKey = "Left Shift"
	KeyboardKeySpace  KeyboardKey = "Space"
	KeyboardKeyEscape KeyboardKey = "Escape"

	KeyboardEventUp   = iota
	KeyboardEventDown = iota
)

type KeyboardInput struct {
	Key    KeyboardKey
	Repeat bool
	Event  KeyboardEvent
}

type KeyboardInputSet map[KeyboardKey]KeyboardInput
