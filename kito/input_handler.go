package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/types"
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

func (g *Game) HandleInput(command interface{}) {
	singleton := g.GetSingleton()

	if _, ok := command.(*types.QuitCommand); ok {
		g.GameOver()
	} else if c, ok := command.(*types.KeyboardInput); ok {
		if _, ok := (*c)[types.KeyboardKeyEscape]; ok {
			// move this into a system maybe
			g.GameOver()
		}
		singleton.SetKeyboardInputSet(c)
	} else if c, ok := command.(*types.MouseInput); ok {
		singleton.SetMouseInput(c)
	} else {
		panic(fmt.Sprintf("UNEXPECTED COMMAND %v", command))
	}
}
