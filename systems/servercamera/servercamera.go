package servercamera

import (
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities/singleton"

	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/kkevinchou/kito/systems/common"
)

const (
	zoomSpeed float64 = 100
	moveSpeed float64 = 25
)

type Singleton interface {
	GetKeyboardInputSet() *input.KeyboardInput
}

type World interface {
	GetSingleton() *singleton.Singleton
	GetCamera() entities.Entity
	GetEntityByID(id int) (entities.Entity, error)
}

type ServerCameraSystem struct {
	*base.BaseSystem
	world World
}

func NewServerCameraSystem(world World) *ServerCameraSystem {
	s := ServerCameraSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
	return &s
}

func (s *ServerCameraSystem) RegisterEntity(entity entities.Entity) {
}

func (s *ServerCameraSystem) Update(delta time.Duration) {
	camera := s.world.GetCamera()
	if camera == nil {
		return
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()
	singleton := s.world.GetSingleton()

	// TODO: This should be player specific
	for _, player := range playerManager.GetPlayers() {
		common.HandleFollowCameraControls(camera, s.world, singleton.PlayerInput[player.ID])
	}
}
