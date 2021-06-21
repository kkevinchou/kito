package kito

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/settings"

	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/types"
)

const (
	fps               float64 = 60
	msPerCommandFrame float64 = 16
	maxTimeStep       float64 = 250 // in milliseconds
)

var (
	cameraStartPosition = mgl64.Vec3{0, 10, 30}
	cameraStartView     = mgl64.Vec2{0, 0}
)

type System interface {
	Update(delta time.Duration)
	RegisterEntity(entity entities.Entity)
	UpdateOnCommandFrame() bool
}

type RenderFunction func(delta time.Duration)

func emptyRenderFunction(delta time.Duration) {}

type Game struct {
	gameOver bool
	camera   entities.Entity
	gameMode types.GameMode

	singleton *singleton.Singleton
	systems   []System
	entities  map[int]entities.Entity
}

func (g *Game) runCommandFrame(delta time.Duration) {
	for _, system := range g.systems {
		if system.UpdateOnCommandFrame() {
			system.Update(delta)
		}
	}
}

func (g *Game) Start(pollInputFunc InputPoller) {
	renderFunction := emptyRenderFunction
	d := directory.GetDirectory()
	renderSystem := d.RenderSystem()
	if renderSystem != nil {
		renderFunction = renderSystem.Update
	}

	var accumulator float64
	var renderAccumulator float64

	msPerFrame := float64(1000) / fps

	var fpsAccumulator float64

	previousTimeStamp := float64(time.Now().UnixNano()) / 1000000
	frameCount := 0
	for !g.gameOver {
		now := float64(time.Now().UnixNano()) / 1000000
		delta := math.Min(now-previousTimeStamp, maxTimeStep)
		previousTimeStamp = now

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= msPerCommandFrame {
			// input is handled once per command frame
			inputList := pollInputFunc()
			for _, input := range inputList {
				g.HandleInput(input)
			}
			g.runCommandFrame(time.Duration(msPerCommandFrame) * time.Millisecond)
			accumulator -= msPerCommandFrame
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
			// fmt.Println(fmt.Sprintf("%d frames rendered last second", frameCount))
			frameCount = 0
			fpsAccumulator -= 1000
		}
	}
}

func (g *Game) GetCamera() entities.Entity {
	return g.camera
}

func (g *Game) GetSingleton() types.Singleton {
	return g.singleton
}

func (g *Game) GetEntityByID(id int) (entities.Entity, error) {
	if entity, ok := g.entities[id]; ok {
		return entity, nil
	}

	return nil, fmt.Errorf("failed to find entity with ID %d", id)
}

func (g *Game) RegisterEntities(entityList []entities.Entity) {
	g.entities = map[int]entities.Entity{}
	for _, entity := range entityList {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range entityList {
		for _, system := range g.systems {
			system.RegisterEntity(entity)
		}
	}
}

func (g *Game) SetCamera(camera entities.Entity) {
	g.camera = camera
}

func compileShaders() {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	shaderManager.CompileShaderProgram("basic", "basic", "basic")
	shaderManager.CompileShaderProgram("basicShadow", "basicshadow", "basicshadow")
	shaderManager.CompileShaderProgram("skybox", "skybox", "skybox")
	shaderManager.CompileShaderProgram("model", "model", "model")
	shaderManager.CompileShaderProgram("depth", "depth", "depth")
	shaderManager.CompileShaderProgram("depthDebug", "basictexture", "depthvalue")
}

func initSeed() {
	seed := settings.Seed
	fmt.Printf("Initializing with seed %d ...\n", seed)
	rand.Seed(seed)
}
