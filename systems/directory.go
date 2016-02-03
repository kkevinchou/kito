package systems

import (
	"sync"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/systems/movement"
	"github.com/kkevinchou/ant/systems/render"
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
// 	Register(interfaces.Renderable)
// 	EventHandlers() []EventHandler
// }

type Directory struct {
	renderSystem   *render.RenderSystem
	movementSystem *movement.MovementSystem
	assetManager   *assets.Manager
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

func (d *Directory) RegisterAssetManager(manager *assets.Manager) {
	d.assetManager = manager
}

func (d *Directory) AssetManager() *assets.Manager {
	return d.assetManager
}

// func (d *Directory) Publish(event Event) {
// }
