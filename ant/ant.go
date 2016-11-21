package ant

import (
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
	grass.New(5, 4)
	grass.New(1, 2)
	grass.New(6, 1)
	grass.New(6, 7)
	grass.New(4, 2)
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

	// renderSystem.Register(pathManager.NavMesh())

	return d
}

type Game struct {
	path      []pathing.Node
	worker    *worker.WorkerImpl
	pathIndex int
}

func (g *Game) Init(window *sdl.Window) {
	setupSystems(window)
	setupGrass()
	// food.New(150, 100)
	// food.New(150, 150)
	// food.New(300, 450)
	// g.worker = worker.New()
	// g.worker.SetPosition(vector.Vector{400, 350})
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
		g.worker.SetTarget(g.path[1].Vector())
	}
}

func (g *Game) PlaceFood(x, y float64) {
	food.New(x, y)
}

func (g *Game) Update() {
	if g.path != nil {
		if g.worker.Position().Sub(g.path[g.pathIndex].Vector()).Length() <= 2 {
			g.pathIndex += 1
			if g.pathIndex == len(g.path) {
				g.path = nil
				g.worker.SetSeekActive(false)
				g.worker.SetVelocity(vector.Zero())
			} else {
				g.worker.SetTarget(g.path[g.pathIndex].Vector())
			}
		}
	}
}