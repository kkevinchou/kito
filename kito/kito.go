package kito

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"

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
}

type RenderFunction func(delta time.Duration)

type Game struct {
	gameOver bool
	camera   entities.Entity
	gameMode types.GameMode

	singleton *singleton.Singleton
	systems   []System
	entities  map[int]entities.Entity

	renderFunction RenderFunction
}

func (g *Game) update(delta time.Duration) {
	for _, system := range g.systems {
		system.Update(delta)
	}
}

func (g *Game) Start(pollInputFunc InputPoller) {
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
			g.update(time.Duration(msPerCommandFrame) * time.Millisecond)
			accumulator -= msPerCommandFrame
		}

		if renderAccumulator >= msPerFrame {
			frameCount++
			g.renderFunction(time.Duration(renderAccumulator) * time.Millisecond)
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
