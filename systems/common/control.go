package common

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/types"
)

func GetControlVector(keyboardInput types.KeyboardInput) mgl64.Vec3 {
	var controlVector mgl64.Vec3
	if key, ok := keyboardInput[types.KeyboardKeyW]; ok && key.Event == types.KeyboardEventDown {
		controlVector[2]++
	}
	if key, ok := keyboardInput[types.KeyboardKeyS]; ok && key.Event == types.KeyboardEventDown {
		controlVector[2]--
	}
	if key, ok := keyboardInput[types.KeyboardKeyA]; ok && key.Event == types.KeyboardEventDown {
		controlVector[0]--
	}
	if key, ok := keyboardInput[types.KeyboardKeyD]; ok && key.Event == types.KeyboardEventDown {
		controlVector[0]++
	}
	if key, ok := keyboardInput[types.KeyboardKeyLShift]; ok && key.Event == types.KeyboardEventDown {
		controlVector[1]--
	}
	if key, ok := keyboardInput[types.KeyboardKeySpace]; ok && key.Event == types.KeyboardEventDown {
		controlVector[1]++
	}
	return controlVector
}
