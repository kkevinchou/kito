package kito

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/metrics"

	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/types"
)

type System interface {
	Update(delta time.Duration)
	RegisterEntity(entity entities.Entity)
}

type RenderFunction func(delta time.Duration)

func emptyRenderFunction(delta time.Duration) {}

type Game struct {
	gameOver bool
	gameMode types.GameMode

	singleton *singleton.Singleton
	systems   []System
	entities  map[int]entities.Entity

	eventBroker     eventbroker.EventBroker
	metricsRegistry *metrics.MetricsRegistry

	// Client
	commandFrameHistory *commandframe.CommandFrameHistory
}

func NewGame() *Game {
	return &Game{
		gameMode:        types.GameModePlaying,
		singleton:       singleton.NewSingleton(),
		entities:        map[int]entities.Entity{},
		eventBroker:     eventbroker.NewEventBroker(),
		metricsRegistry: metrics.New(),
	}
}

func (g *Game) Start(pollInputFunc input.InputPoller) {
	var accumulator float64
	var renderAccumulator float64
	var fpsAccumulator float64

	msPerFrame := float64(1000) / settings.FPS
	previousTimeStamp := float64(time.Now().UnixNano()) / 1000000

	frameCount := 0
	renderFunction := getRenderFunction()
	for !g.gameOver {
		now := float64(time.Now().UnixNano()) / 1000000
		delta := math.Min(now-previousTimeStamp, settings.MaxTimeStepMS)
		if delta == settings.MaxTimeStepMS {
			fmt.Println("hit settings.MaxTimeStepMS - simulation time has been lost")
		}
		previousTimeStamp = now

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= settings.MSPerCommandFrame {
			// input is handled once per command frame
			g.HandleInput(pollInputFunc())
			g.runCommandFrame(time.Duration(settings.MSPerCommandFrame) * time.Millisecond)
			accumulator -= settings.MSPerCommandFrame
		}

		if renderAccumulator >= msPerFrame {
			frameCount++
			g.metricsRegistry.Inc("fps", 1)
			renderFunction(time.Duration(renderAccumulator) * time.Millisecond)
			for renderAccumulator >= msPerFrame {
				renderAccumulator -= msPerFrame
			}
		}

		fpsAccumulator += delta
		if fpsAccumulator > 1000 {
			// fmt.Println("FPS:", frameCount)
			frameCount = 0
			fpsAccumulator -= 1000
		}
	}
}

func (g *Game) runCommandFrame(delta time.Duration) {
	g.singleton.CommandFrame++
	for _, system := range g.systems {
		system.Update(delta)
	}
}

func initSeed() {
	seed := settings.Seed
	fmt.Printf("initializing with seed %d ...\n", seed)
	rand.Seed(seed)
}

func getRenderFunction() RenderFunction {
	renderFunction := emptyRenderFunction
	d := directory.GetDirectory()
	renderSystem := d.RenderSystem()
	if renderSystem != nil {
		renderFunction = renderSystem.Render
	}

	return renderFunction
}
