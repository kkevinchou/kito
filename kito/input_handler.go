package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/lib/input"
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

func (g *Game) HandleInput(frameInput input.Input) {
	singleton := g.GetSingleton()
	singleton.PlayerInput[singleton.PlayerID] = frameInput

	keyboardInput := frameInput.KeyboardInput
	if _, ok := keyboardInput[input.KeyboardKeyEscape]; ok {
		// move this into a system maybe
		g.GameOver()
	}

	for _, cmd := range frameInput.Commands {
		if _, ok := cmd.(input.QuitCommand); ok {
			g.GameOver()
		} else {
			panic(fmt.Sprintf("UNEXPECTED COMMAND %v", cmd))
		}
	}
}
