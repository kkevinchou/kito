package camera

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/types"
)

type Singleton interface {
	GetKeyboardInputSet() *commands.KeyboardInputSet
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
	keyboardInputSet := *singleton.GetKeyboardInputSet()

	var controlVector vector.Vector3
	if key, ok := keyboardInputSet[commands.KeyboardKeyW]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Z--
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyS]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Z++
	}
	// Left
	if key, ok := keyboardInputSet[commands.KeyboardKeyA]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.X--
	}
	// Right
	if key, ok := keyboardInputSet[commands.KeyboardKeyD]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.X++
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyLShift]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Y--
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeySpace]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Y++
	}

	mouseInput := singleton.GetMouseInput()

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
