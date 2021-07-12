package camera

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

type CameraSystem struct {
	*base.BaseSystem
	world World
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
	return &s
}

func (s *CameraSystem) RegisterEntity(entity entities.Entity) {
}

func (s *CameraSystem) Update(delta time.Duration) {
	// TODO: we will need to support multiple cameras (one for each player)
	camera := s.world.GetCamera()
	if camera == nil {
		return
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()
	singleton := s.world.GetSingleton()

	for _, player := range playerManager.GetPlayers() {
		common.HandleFollowCameraControls(camera, s.world, singleton.PlayerInput[player.ID])
	}
}
