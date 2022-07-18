package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/input"
)

func (g *Game) GameOver() {
	g.gameOver = true
}

func (g *Game) HandleInput(frameInput input.Input) {

	keyboardInput := frameInput.KeyboardInput
	if _, ok := keyboardInput[input.KeyboardKeyEscape]; ok {
		// move this into a system maybe
		g.GameOver()
	}

	if keyEvent, ok := keyboardInput[input.KeyboardKeyTick]; ok {
		if keyEvent.Event == input.KeyboardEventUp {
			g.ToggleWindowVisibility(types.WindowConsole)
		}
	}

	if keyEvent, ok := keyboardInput[input.KeyboardKeyF1]; ok {
		if keyEvent.Event == input.KeyboardEventUp {
			g.ToggleWindowVisibility(types.WindowDebug)
		}
	}

	for _, cmd := range frameInput.Commands {
		if _, ok := cmd.(input.QuitCommand); ok {
			g.GameOver()
		} else {
			panic(fmt.Sprintf("UNEXPECTED COMMAND %v", cmd))
		}
	}

	if g.GetFocusedWindow() == types.WindowGame {
		singleton := g.GetSingleton()
		singleton.PlayerInput[singleton.PlayerID] = frameInput
	}
}
