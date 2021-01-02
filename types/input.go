package types

type MouseWheelDirection string

const (
	MouseWheelDirectionUp      MouseWheelDirection = "UP"
	MouseWheelDirectionDown    MouseWheelDirection = "DOWN"
	MouseWheelDirectionNeutral MouseWheelDirection = "NEUTRAL"
)

type MouseMotionEvent struct {
	XRel float64
	YRel float64
}

type MouseInput struct {
	MouseWheel       MouseWheelDirection
	MouseMotionEvent *MouseMotionEvent
	LeftButtonDown   bool
	MiddleButtonDown bool
	RightButtonDown  bool
}

type KeyboardKey string
type KeyboardEvent int

const (
	KeyboardKeyW KeyboardKey = "W"
	KeyboardKeyA KeyboardKey = "A"
	KeyboardKeyS KeyboardKey = "S"
	KeyboardKeyD KeyboardKey = "D"

	KeyboardKeyUp    KeyboardKey = "Up"
	KeyboardKeyDown  KeyboardKey = "Down"
	KeyboardKeyLeft  KeyboardKey = "Left"
	KeyboardKeyRight KeyboardKey = "Right"

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

type QuitCommand struct {
}
