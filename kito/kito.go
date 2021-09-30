package kito

import (
	"fmt"
	"math"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"

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

	eventBroker         eventbroker.EventBroker
	commandFrameHistory *commandframe.CommandFrameHistory
}

func (g *Game) runCommandFrame(delta time.Duration) {
	g.singleton.CommandFrame++
	for _, system := range g.systems {
		system.Update(delta)
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

func (g *Game) GetSingleton() *singleton.Singleton {
	return g.singleton
}

func (g *Game) GetEntityByID(id int) (entities.Entity, error) {
	if entity, ok := g.entities[id]; ok {
		return entity, nil
	}

	stack := debug.Stack()

	return nil, fmt.Errorf("%sfailed to find entity with ID %d", string(stack), id)
}

func (g *Game) RegisterEntities(entityList []entities.Entity) {
	for _, entity := range entityList {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range entityList {
		for _, system := range g.systems {
			system.RegisterEntity(entity)
		}
	}
}

func (g *Game) CommandFrame() int {
	return g.singleton.CommandFrame
}

func (g *Game) GetEventBroker() eventbroker.EventBroker {
	return g.eventBroker
}

func (g *Game) GetCommandFrameHistory() *commandframe.CommandFrameHistory {
	return g.commandFrameHistory
}

func compileShaders() {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	if err := shaderManager.CompileShaderProgram("basic", "basic", "basic"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("basicShadow", "basicshadow", "basicshadow"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("skybox", "skybox", "skybox"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model", "model", "model"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model_static", "model_static", "model"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("depthDebug", "basictexture", "depthvalue"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("material", "model", "material"); err != nil {
		panic(err)
	}
}

func initSeed() {
	seed := settings.Seed
	fmt.Printf("Initializing with seed %d ...\n", seed)
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
