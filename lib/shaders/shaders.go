package shaders

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexShaderPath, fragmentShaderPath string) (*Shader, error) {
	vertexShaderSource, err := ioutil.ReadFile(vertexShaderPath)
	if err != nil {
		return nil, err
	}

	fragmentShaderSource, err := ioutil.ReadFile(fragmentShaderPath)
	if err != nil {
		return nil, err
	}

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csource, free := gl.Strs(string(vertexShaderSource) + "\x00")
	gl.ShaderSource(vertexShader, 1, csource, nil)
	free()

	gl.CompileShader(vertexShader)

	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))
		panic("Failed to compile vertex shader:\n" + log)
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource, free = gl.Strs(string(fragmentShaderSource) + "\x00")
	gl.ShaderSource(fragmentShader, 1, csource, nil)
	free()

	gl.CompileShader(fragmentShader)

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(log))
		panic("Failed to compile fragment shader:\n" + log)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic("Failed to link program:\n" + log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{ID: shaderProgram}, nil
}

func (s *Shader) SetUniformMat4(uniform string, value mgl32.Mat4) {
	uniformLocation := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", uniform)))
	gl.UniformMatrix4fv(uniformLocation, 1, false, &value[0])
}

func (s *Shader) SetUniformVec3(uniform string, value mgl32.Vec3) {
	floats := []float32{value.X(), value.Y(), value.Z()}
	uniformLocation := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", uniform)))
	gl.Uniform3fv(uniformLocation, 1, &floats[0])
}

func (s *Shader) SetUniformInt(uniform string, value int32) {
	uniformLocation := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", uniform)))
	gl.Uniform1i(uniformLocation, value)
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}
