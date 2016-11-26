package ant

import "github.com/kkevinchou/ant/lib/math/vector"

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
	game.MoveCamera(vector.Vector3{X: c.X, Y: c.Y, Z: c.Z})
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
	game.CameraViewChange(vector.Vector{X: c.X, Y: c.Y})
}
