package components

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

type EasingFunction func(x float64) float64

type EasingComponent struct {
	active      bool
	currentTime time.Duration
	duration    time.Duration

	// easingFunction should sample the function at x and returns (which should return a value of 1)
	easingFunction EasingFunction
}

func NewEasingComponent(duration time.Duration, f EasingFunction) *EasingComponent {
	return &EasingComponent{easingFunction: f, duration: duration}
}

func (c *EasingComponent) AddToComponentContainer(container *ComponentContainer) {
	container.EasingComponent = c
}

func (c *EasingComponent) Update(delta time.Duration) {
	if !c.active {
		return
	}
	c.currentTime += delta
}

func (c *EasingComponent) Start() {
	c.active = true
	c.currentTime = 0
}

func (c *EasingComponent) Stop() {
	c.active = false
}

func (c *EasingComponent) Active() bool {
	return c.active
}

func (c *EasingComponent) GetValue() float64 {
	xFactor := float64(c.duration.Milliseconds())
	x := float64(c.currentTime.Milliseconds()) / xFactor
	return c.easingFunction(x)
}

func EaseInOutCubic(x float64) float64 {
	x = mgl64.Clamp(x, 0, 1)

	if x < 0.5 {
		return 4 * math.Pow(x, 3)
	}

	return 1 - math.Pow(-2*x+2, 3)/2
}

func EaseOutSine(x float64) float64 {
	x = mgl64.Clamp(x, 0, 1)

	return math.Sin((x * math.Pi) / 2)
}

func EaseInOutSine(x float64) float64 {
	x = mgl64.Clamp(x, 0, 1)

	return -(math.Cos(x*math.Pi) - 1) / 2
}

func EaseInOutCirc(x float64) float64 {
	x = mgl64.Clamp(x, 0, 1)

	if x < 0.5 {
		return (1 - math.Sqrt(1-math.Pow(2*x, 2))) / 2
	}

	return (math.Sqrt(1-math.Pow(-2*x+2, 2)) + 1) / 2
}
