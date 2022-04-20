package kito

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entitymanager"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/metrics"

	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/types"
)

type System interface {
	Update(delta time.Duration)
}

type RenderFunction func(delta time.Duration)

func emptyRenderFunction(delta time.Duration) {}

type Game struct {
	gameOver bool
	gameMode types.GameMode

	singleton     *singleton.Singleton
	entityManager *entitymanager.EntityManager
	systems       []System

	eventBroker     eventbroker.EventBroker
	metricsRegistry *metrics.MetricsRegistry

	inputPollingFn input.InputPoller

	// Client
	commandFrameHistory *commandframe.CommandFrameHistory
}

func NewBaseGame() *Game {
	return &Game{
		gameMode:        types.GameModePlaying,
		singleton:       singleton.NewSingleton(),
		entityManager:   entitymanager.NewEntityManager(),
		eventBroker:     eventbroker.NewEventBroker(),
		metricsRegistry: metrics.New(),
		inputPollingFn:  input.NullInputPoller,
	}
}

func (g *Game) Start() {
	var accumulator float64
	var renderAccumulator float64

	msPerFrame := float64(1000) / float64(settings.FPS)
	previousTimeStamp := float64(time.Now().UnixNano()) / 1000000

	frameCount := 0
	renderFunction := getRenderFunction()
	for !g.gameOver {
		now := float64(time.Now().UnixNano()) / 1000000
		delta := now - previousTimeStamp
		previousTimeStamp = now

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= float64(settings.MSPerCommandFrame) {
			// input is handled once per command frame
			g.HandleInput(g.inputPollingFn())
			g.runCommandFrame(time.Duration(settings.MSPerCommandFrame) * time.Millisecond)
			accumulator -= float64(settings.MSPerCommandFrame)
		}

		if renderAccumulator >= msPerFrame {
			frameCount++
			g.metricsRegistry.Inc("fps", 1)
			// renderFunction(time.Duration(renderAccumulator) * time.Millisecond)
			renderFunction(time.Duration(delta) * time.Millisecond)
			renderAccumulator -= msPerFrame
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
