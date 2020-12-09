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
