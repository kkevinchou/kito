package render

import (
	"time"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	maxY float64 = 25
	minY float64 = 5
)

type Light struct {
	id       uint32
	ambient  []float32
	diffuse  []float32
	specular []float32

	position mgl64.Vec3
}

var direction int = 1

func NewLight(id uint32) *Light {
	light := &Light{
		id:       id,
		position: mgl64.Vec3{0, maxY, 0},
		ambient:  []float32{0.25, 0.25, 0.25, 1},
		diffuse:  []float32{1, 1, 1, 1},
		specular: []float32{1, 1, 1, 1},
	}

	return light
}

func (l *Light) SetPosition(position mgl64.Vec3) {
	l.position = position
}

func (l *Light) Update(delta time.Duration) {
	if l.position.Y() < minY {
		direction = -1
	}

	if l.position.Y() > maxY {
		direction = 1
	}

	l.position = l.position.Sub(mgl64.Vec3{0, 5 * delta.Seconds() * float64(direction), 0})
	position := []float32{float32(l.position.X()), float32(l.position.Y()), float32(l.position.Z()), 1}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &position[0])
}

func (l *Light) Position() mgl64.Vec3 {
	return l.position
}
