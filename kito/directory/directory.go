package directory

import (
	"sync"
	"time"

	"github.com/kkevinchou/kito/kito/managers/item"
	"github.com/kkevinchou/kito/kito/managers/path"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/modelspec"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/lib/textures"
)

type IAssetManager interface {
	GetTexture(name string) *textures.Texture
	GetAnimatedModel(name string) *modelspec.ModelSpecification
}

type IRenderSystem interface {
	Render(time.Duration)
}

type IShaderManager interface {
	CompileShaderProgram(name, vertexShader, fragmentShader string) error
	GetShaderProgram(name string) *shaders.ShaderProgram
}

type IPlayerManager interface {
	RegisterPlayer(id int, client types.NetworkClient)
	GetPlayer(id int) *player.Player
	GetPlayers() []*player.Player
}

type Directory struct {
	renderSystem  IRenderSystem
	assetManager  IAssetManager
	itemManager   *item.Manager
	pathManager   *path.Manager
	shaderManager IShaderManager
	playerManager IPlayerManager
}

var instance *Directory
var once sync.Once

func GetDirectory() *Directory {
	once.Do(func() {
		instance = &Directory{}
	})
	return instance
}

func (d *Directory) RegisterRenderSystem(system IRenderSystem) {
	d.renderSystem = system
}

func (d *Directory) RenderSystem() IRenderSystem {
	return d.renderSystem
}

func (d *Directory) RegisterAssetManager(manager IAssetManager) {
	d.assetManager = manager
}

func (d *Directory) AssetManager() IAssetManager {
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

func (d *Directory) RegisterShaderManager(manager IShaderManager) {
	d.shaderManager = manager
}

func (d *Directory) ShaderManager() IShaderManager {
	return d.shaderManager
}

func (d *Directory) RegisterPlayerManager(manager IPlayerManager) {
	d.playerManager = manager
}

func (d *Directory) PlayerManager() IPlayerManager {
	return d.playerManager
}
