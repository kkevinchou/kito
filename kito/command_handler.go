package kito

import (
	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
)

// type CameraRaycastCommand struct {
// 	X float64
// 	Y float64
// }

// func (c *CameraRaycastCommand) Execute(game *Game) {
// 	renderSystem := directory.GetDirectory().RenderSystem()
// 	worldPoint := renderSystem.GetWorldPoint(c.X, c.Y)
// 	dir := worldPoint.Sub(game.camera.Position()).Normalize()
// 	render.LineStart = game.camera.Position()
// 	render.LineEnd = game.camera.Position().Add(dir.Scale(3))
// 	fmt.Println(worldPoint)
// 	fmt.Println("Camera position:", game.camera.Position(), "Direction:", dir)
// }

func (g *Game) GameOver() {
	g.gameOver = true
}

func (g *Game) MoveCommand(vector vector.Vector3) {
	g.viewer.SetVelocityDirection(vector)
}

func (g *Game) UpdateViewCommand(vector vector.Vector) {
	if g.viewControlled {
		g.viewer.UpdateView(vector)
	}
}

func (g *Game) ToggleCameraControlCommand(value bool) {
	g.viewControlled = value
}

func (g *Game) Handle(command interface{}) {
	if _, ok := command.(*commands.QuitCommand); ok {
		g.GameOver()
	} else if c, ok := command.(*commands.MoveCommand); ok {
		g.MoveCommand(c.Value)
	} else if c, ok := command.(*commands.UpdateViewCommand); ok {
		g.UpdateViewCommand(c.Value)
	} else if c, ok := command.(*commands.ToggleCameraControlCommand); ok {
		g.ToggleCameraControlCommand(c.Value)
	}
}