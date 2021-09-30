package input

type InputPoller func() Input

type MouseWheelDirection int

type MouseMotionEvent struct {
	XRel float64
	YRel float64
}

func (m MouseMotionEvent) IsZero() bool {
	return m.XRel == 0 && m.YRel == 0
}

type MouseInput struct {
	MouseWheelDelta  int
	MouseMotionEvent MouseMotionEvent
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

	KeyboardEventUp = iota
	KeyboardEventDown
)

type KeyState struct {
	Key   KeyboardKey
	Event KeyboardEvent
}

type KeyboardInput map[KeyboardKey]KeyState

type QuitCommand struct {
}

type Input struct {
	NewInput      bool
	KeyboardInput KeyboardInput
	MouseInput    MouseInput
	Commands      []interface{}
}

func (i Input) Copy() Input {
	keyboardInput := KeyboardInput{}
	for k, v := range i.KeyboardInput {
		keyboardInput[k] = v
	}

	return Input{
		KeyboardInput: keyboardInput,
		MouseInput:    i.MouseInput,
	}
}
