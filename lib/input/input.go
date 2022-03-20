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
	Buttons          [3]bool // left, right, middle
}

type KeyboardKey string
type KeyboardEvent int

const (
	KeyboardKeyW KeyboardKey = "W"
	KeyboardKeyA KeyboardKey = "A"
	KeyboardKeyS KeyboardKey = "S"
	KeyboardKeyD KeyboardKey = "D"
	KeyboardKeyQ KeyboardKey = "Q"

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

// Input represents the input provided by a user during a command frame
// Input should be only constructed by the input poller and should not be
// written to by any systems, only read. Input is stored in a client side
// command frame history which will copy the KeyboardInput by reference
type Input struct {
	KeyboardInput KeyboardInput
	MouseInput    MouseInput
	Commands      []any
}

// func (i Input) Copy() Input {
// 	keyboardInput := KeyboardInput{}
// 	for k, v := range i.KeyboardInput {
// 		keyboardInput[k] = v
// 	}

// 	return Input{
// 		KeyboardInput: keyboardInput,
// 		MouseInput:    i.MouseInput,
// 	}
// }
