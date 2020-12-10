package camera

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/types"
)

type Singleton interface {
	GetKeyboardInputSet() *types.KeyboardInput
}

type World interface {
	GetSingleton() types.Singleton
	GetCamera() types.Viewer
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
	mouseInput := *singleton.GetMouseInput()

	var controlVector vector.Vector3
	if key, ok := keyboardInput[types.KeyboardKeyW]; ok && key.Event == types.KeyboardEventDown {
		controlVector.Z--
	}
	if key, ok := keyboardInput[types.KeyboardKeyS]; ok && key.Event == types.KeyboardEventDown {
		controlVector.Z++
	}
	// Left
	if key, ok := keyboardInput[types.KeyboardKeyA]; ok && key.Event == types.KeyboardEventDown {
		controlVector.X--
	}
	// Right
	if key, ok := keyboardInput[types.KeyboardKeyD]; ok && key.Event == types.KeyboardEventDown {
		controlVector.X++
	}
	if key, ok := keyboardInput[types.KeyboardKeyLShift]; ok && key.Event == types.KeyboardEventDown {
		controlVector.Y--
	}
	if key, ok := keyboardInput[types.KeyboardKeySpace]; ok && key.Event == types.KeyboardEventDown {
		controlVector.Y++
	}

	if mouseInput.LeftButtonDown && mouseInput.MouseMotionEvent != nil {
		camera.UpdateView(vector.Vector{X: mouseInput.MouseMotionEvent.XRel, Y: mouseInput.MouseMotionEvent.YRel})
	}

	zoomValue := 0
	if mouseInput.MouseWheel == types.MouseWheelDirectionNeutral {
		zoomValue = 0
	} else if mouseInput.MouseWheel == types.MouseWheelDirectionUp {
		zoomValue = -1
	} else if mouseInput.MouseWheel == types.MouseWheelDirectionDown {
		zoomValue = 1
	} else {
		panic(fmt.Sprintf("unexpected mousewheel value %v", mouseInput.MouseWheel))
	}

	if controlVector.IsZero() && zoomValue == 0 {
		return
	}

	forwardVector := camera.Forward()
	zoomVector := forwardVector.Scale(float64(zoomValue))

	forwardVector = forwardVector.Scale(controlVector.Z)
	forwardVector.Y = 0

	rightVector := camera.Right()
	rightVector = rightVector.Scale(-controlVector.X)

	impulse := &types.Impulse{}
	if !controlVector.IsZero() {
		moveSpeed := camera.MaxSpeed()
		impulse.Vector = forwardVector.Add(rightVector).Add(vector.Vector3{X: 0, Y: controlVector.Y, Z: 0}).Normalize().Scale(moveSpeed)
		impulse.DecayRate = 2.5
		camera.ApplyImpulse("cameraMove", impulse)
	}

	if zoomValue != 0 {
		zoomSpeed := 2 * camera.MaxSpeed()
		impulse.Vector = zoomVector.Scale(zoomSpeed)
		impulse.DecayRate = 2.5
		camera.ApplyImpulse("cameraZoom", impulse)
	}
}
