package directory

import (
	"sync"

	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/movement"
	"github.com/kkevinchou/kito/systems/render"
)

// type EventType int

// const (
// 	EntityCreated EventType = iota
// )

// type Event interface {
// 	Data() map[string]string
// 	Type() EventType
// }

// type System interface {
// }

// type EventHandler struct {
// 	Type    EventType
// 	Handler func(Event)
// }

// type RenderSystemI interface {
// 	System
// 	Register(types.Renderable)
// 	EventHandlers() []EventHandler
// }

type Directory struct {
	renderSystem   *render.RenderSystem
	movementSystem *movement.MovementSystem
	assetManager   *lib.AssetManager
	itemManager    *item.Manager
	pathManager    *path.Manager
	// aiManager      *ai.Manager
}

var instance *Directory
var once sync.Once

func GetDirectory() *Directory {
	once.Do(func() {
		instance = &Directory{}
	})
	return instance
}

func (d *Directory) RegisterRenderSystem(system *render.RenderSystem) {
	d.renderSystem = system
}

func (d *Directory) RenderSystem() *render.RenderSystem {
	return d.renderSystem
}

func (d *Directory) RegisterMovementSystem(system *movement.MovementSystem) {
	d.movementSystem = system
}

func (d *Directory) MovementSystem() *movement.MovementSystem {
	return d.movementSystem
}

func (d *Directory) RegisterAssetManager(manager *lib.AssetManager) {
	d.assetManager = manager
}

func (d *Directory) AssetManager() *lib.AssetManager {
	return d.assetManager
}

func (d *Directory) RegisterItemManager(manager *item.Manager) {
	d.itemManager = manager
}

func (d *Directory) ItemManager() *item.Manager {
	return d.itemManager
}

func (d *Directory) RegisterPathManager(manager *path.Manager) {
	d.pathManager = manager
}

func (d *Directory) PathManager() *path.Manager {
	return d.pathManager
}

// func (d *Directory) RegisterAIManager(manager *ai.Manager) {
// 	d.aiManager = manager
// }

// func (d *Directory) AIManager() *ai.Manager {
// 	return d.aiManager
// }

// func (d *Directory) Publish(event Event) {
// }
