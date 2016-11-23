package ant

import (
	"math/rand"
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/entities/food"
	"github.com/kkevinchou/ant/entities/grass"
	"github.com/kkevinchou/ant/entities/worker"
	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/pathing"
	"github.com/kkevinchou/ant/managers/item"
	"github.com/kkevinchou/ant/managers/path"
	"github.com/kkevinchou/ant/systems/movement"
	"github.com/kkevinchou/ant/systems/render"
	"github.com/veandco/go-sdl2/sdl"
)

func setupGrass() {
	grass.New(5, 0, 4)
	grass.New(2, 0, 2)
	grass.New(6, 0, 1)
	grass.New(6, 0, 7)
	grass.New(4, 0, 2)
	grass.New(0, 0, 0)
}

func setupSystems(window *sdl.Window) *directory.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := lib.NewAssetManager(nil, "_assets")
	renderSystem := render.NewRenderSystem(window, assetManager)
	movementSystem := movement.NewMovementSystem()

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterMovementSystem(movementSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	renderSystem.Register(pathManager.NavMesh())

	return d
}

type Game struct {
	path      []pathing.Node
	worker    *worker.WorkerImpl
	pathIndex int
}

func (g *Game) Init(window *sdl.Window) {
	rand.Seed(int64(time.Now().Nanosecond()))
	setupSystems(window)
	setupGrass()
	food.New(0, 0, 0)
	g.worker = worker.New()
	g.worker.SetPosition(vector.Vector3{19, 12, -10})
}

func (g *Game) MoveAnt(x, y float64) {
	position := g.worker.Position()
	pathManager := directory.GetDirectory().PathManager()
	g.path = pathManager.FindPath(
		geometry.Point{X: position.X, Y: position.Y},
		geometry.Point{X: x, Y: y},
	)
	if g.path != nil {
		g.pathIndex = 1
		g.worker.SetTarget(g.path[1].Vector3())
	}
}

func (g *Game) PlaceFood(x, y float64) {
	food.New(x, 0, y)
}

func (g *Game) CameraView(x, y int) {
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.CameraView(x, y)
}

func (g *Game) MoveCamera(v vector.Vector3) {
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.MoveCamera(v)
}

func (g *Game) Update() {
	if g.path != nil {
		if g.worker.Position().Sub(g.path[g.pathIndex].Vector3()).Length() <= 2 {
			g.pathIndex += 1
			if g.pathIndex == len(g.path) {
				g.path = nil
				g.worker.SetSeekActive(false)
				g.worker.SetVelocity(vector.Vector3{})
			} else {
				g.worker.SetTarget(g.path[g.pathIndex].Vector3())
			}
		}
	}
}
