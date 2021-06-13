package directory

import (
	"sync"

	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/render"
)

type Directory struct {
	renderSystem *render.RenderSystem
	assetManager *assets.AssetManager
	itemManager  *item.Manager
	pathManager  *path.Manager
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

func (d *Directory) RegisterAssetManager(manager *assets.AssetManager) {
	d.assetManager = manager
}

func (d *Directory) AssetManager() *assets.AssetManager {
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
