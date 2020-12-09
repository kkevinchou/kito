package types

type MouseWheelDirection string

const (
	MouseWheelDirectionUp      MouseWheelDirection = "UP"
	MouseWheelDirectionDown    MouseWheelDirection = "DOWN"
	MouseWheelDirectionNeutral MouseWheelDirection = "NEUTRAL"
)

type MouseInput struct {
	MouseWheel MouseWheelDirection
}

type KeyboardKey string
type KeyboardEvent int

const (
	KeyboardKeyW KeyboardKey = "W"
	KeyboardKeyA KeyboardKey = "A"
	KeyboardKeyS KeyboardKey = "S"
	KeyboardKeyD KeyboardKey = "D"

	KeyboardKeyLShift KeyboardKey = "Left Shift"
	KeyboardKeySpace  KeyboardKey = "Space"
	KeyboardKeyEscape KeyboardKey = "Escape"

	KeyboardEventUp   = iota
	KeyboardEventDown = iota
)

type KeyState struct {
	Key   KeyboardKey
	Event KeyboardEvent
}

type KeyboardInput map[KeyboardKey]KeyState
