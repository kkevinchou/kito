package ant

import (
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/entities/food"
	"github.com/kkevinchou/ant/entities/grass"
	"github.com/kkevinchou/ant/entities/worker"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/managers/item"
	"github.com/kkevinchou/ant/managers/path"
	"github.com/kkevinchou/ant/pathing"
	"github.com/kkevinchou/ant/systems"
	"github.com/kkevinchou/ant/systems/movement"
	"github.com/kkevinchou/ant/systems/render"
	"github.com/veandco/go-sdl2/sdl"
)

func setupGrass() {
	grass.New(366, 450)
	grass.New(386, 450)
	grass.New(406, 450)
	grass.New(406, 350)
	grass.New(436, 350)
}

func setupSystems(renderer *sdl.Renderer) *systems.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := assets.NewAssetManager(renderer, "assets")
	renderSystem := render.NewRenderSystem(renderer, assetManager)
	movementSystem := movement.NewMovementSystem()

	d := systems.GetDirectory()
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
	worker    *worker.Worker
	pathIndex int
}

func (g *Game) Init(renderer *sdl.Renderer) {
	setupSystems(renderer)
	setupGrass()
	food.New(150, 100)
	g.worker = worker.New()
	g.worker.SetPosition(vector.Vector{400, 350})
}

func (g *Game) MoveAnt(x, y float64) {
	position := g.worker.Position()
	pathManager := systems.GetDirectory().PathManager()
	g.path = pathManager.FindPath(
		geometry.Point{X: position.X, Y: position.Y},
		geometry.Point{X: x, Y: y},
	)
	if g.path != nil {
		g.pathIndex = 1
		g.worker.SetTarget(g.path[1].Vector())
	}
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
