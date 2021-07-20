package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/lib/input"
)

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
