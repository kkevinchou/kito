package render

import (
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type Light struct {
	id       uint32
	ambient  []float32
	diffuse  []float32
	specular []float32

	position vector.Vector3
}

var direction int = 1

func NewLight(id uint32) *Light {
	light := &Light{
		id:       id,
		position: vector.Vector3{0, 20, 1},
		ambient:  []float32{0.25, 0.25, 0.25, 1},
		diffuse:  []float32{1, 1, 1, 1},
		specular: []float32{1, 1, 1, 1},
	}

	// gl.Lightfv(id, gl.AMBIENT, &light.ambient[0])
	gl.Lightfv(id, gl.DIFFUSE, &light.diffuse[0])
	gl.Lightfv(id, gl.SPECULAR, &light.specular[0])
	gl.Enable(id)

	return light
}

func (l *Light) SetPosition(position vector.Vector3) {
	l.position = position
}

func (l *Light) Update(delta time.Duration) {
	if l.position.Y < 0 {
		direction = -1
	}

	if l.position.Y > 20 {
		direction = 1
	}

	l.position = l.position.Sub(vector.Vector3{0, 5 * float64(delta.Seconds()) * float64(direction), 0})
	position := []float32{float32(l.position.X), float32(l.position.Y), float32(l.position.Z), 1}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &position[0])
}

func (l *Light) Position() vector.Vector3 {
	return l.position
}
