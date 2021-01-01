package render

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Quad struct {
	vao uint32
}

func NewQuad(vertixAttributes []float32) *Quad {
	vertixAttributes = []float32{
		// bottom
		-0.5, 0, -0.5, 0, 1, 0,
		0.5, 0, 0.5, 0, 1, 0,
		0.5, 0, -0.5, 0, 1, 0,
		0.5, 0, 0.5, 0, 1, 0,
		-0.5, 0, -0.5, 0, 1, 0,
		-0.5, 0, 0.5, 0, 1, 0,
	}

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertixAttributes)*4, gl.Ptr(vertixAttributes), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	q := Quad{
		vao: vao,
	}
	return &q
}

func (q *Quad) GetVAO() uint32 {
	return q.vao
}
