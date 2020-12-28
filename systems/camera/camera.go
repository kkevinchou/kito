package camera

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/types"
)

const (
	zoomSpeed float64 = 100
	moveSpeed float64 = 25
)

type Singleton interface {
	GetKeyboardInputSet() *types.KeyboardInput
}

type World interface {
	GetSingleton() types.Singleton
	GetCamera() types.Camera
}

type CameraSystem struct {
	world World
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		world: world,
	}
	return &s
}

func (s *CameraSystem) Update(delta time.Duration) {
	camera := s.world.GetCamera()

	if !camera.Controlled() {
		return
	}

	singleton := s.world.GetSingleton()
	keyboardInput := *singleton.GetKeyboardInputSet()

	var controlVector mgl64.Vec3
	if key, ok := keyboardInput[types.KeyboardKeyW]; ok && key.Event == types.KeyboardEventDown {
		controlVector[2]--
	}
	if key, ok := keyboardInput[types.KeyboardKeyS]; ok && key.Event == types.KeyboardEventDown {
		controlVector[2]++
	}
	// Left
	if key, ok := keyboardInput[types.KeyboardKeyA]; ok && key.Event == types.KeyboardEventDown {
		controlVector[0]--
	}
	// Right
	if key, ok := keyboardInput[types.KeyboardKeyD]; ok && key.Event == types.KeyboardEventDown {
		controlVector[0]++
	}
	if key, ok := keyboardInput[types.KeyboardKeyLShift]; ok && key.Event == types.KeyboardEventDown {
		controlVector[1]--
	}
	if key, ok := keyboardInput[types.KeyboardKeySpace]; ok && key.Event == types.KeyboardEventDown {
		controlVector[1]++
	}

	zoomValue := 0
	if singleton.GetMouseInput() != nil {
		mouseInput := *singleton.GetMouseInput()
		if mouseInput.LeftButtonDown && mouseInput.MouseMotionEvent != nil {
			camera.UpdateView(mgl64.Vec2{mouseInput.MouseMotionEvent.XRel, mouseInput.MouseMotionEvent.YRel})
		}

		if mouseInput.MouseWheel == types.MouseWheelDirectionNeutral {
			zoomValue = 0
		} else if mouseInput.MouseWheel == types.MouseWheelDirectionUp {
			zoomValue = -1
		} else if mouseInput.MouseWheel == types.MouseWheelDirectionDown {
			zoomValue = 1
		} else {
			panic(fmt.Sprintf("unexpected mousewheel value %v", mouseInput.MouseWheel))
		}
	}
	if utils.Vec3IsZero(controlVector) && zoomValue == 0 {
		return
	}

	forwardVector := camera.Forward()
	zoomVector := forwardVector.Mul(float64(zoomValue))

	forwardVector = forwardVector.Mul(controlVector.Z())
	forwardVector[1] = 0

	rightVector := camera.Right()
	rightVector = rightVector.Mul(-controlVector.X())

	impulse := &types.Impulse{}
	if !utils.Vec3IsZero(controlVector) {
		impulse.Vector = forwardVector.Add(rightVector).Add(mgl64.Vec3{0, controlVector.Y(), 0}).Normalize().Mul(moveSpeed)
		impulse.DecayRate = 2.5
		camera.ApplyImpulse("cameraMove", impulse)
	}

	if zoomValue != 0 {
		impulse.Vector = zoomVector.Mul(zoomSpeed)
		impulse.DecayRate = 2.5
		camera.ApplyImpulse("cameraZoom", impulse)
	}
}
