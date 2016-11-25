package models

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/kkevinchou/ant/lib/math/vector"
)

// Source: https://gist.github.com/davemackintosh/67959fa9dfd9018d79a4
// https://github.com/alanmacleod/parse-obj3d/blob/master/parse.go

type Face struct {
	Verticies []*FaceVertex
	Material  string
}

type FaceVertex struct {
	Vertex  *vector.Vector3
	Texture *vector.Vector
	Normal  *vector.Vector3
}

type Model struct {
	normals   []*vector.Vector3
	textures  []*vector.Vector
	verticies []*vector.Vector3
	Faces     []*Face
}

// NewModel reads an OBJ model file and creates a Model from its contents
func NewModel(file string) (*Model, error) {
	objFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	model := Model{}

	scanner := bufio.NewScanner(objFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if len(line) == 1 {
			return nil, errors.New("line of length 1 unexpected")
		}

		if strings.HasPrefix(line, "vt") {
			var err error
			var x, y float64

			split := strings.Split(line, " ")
			if x, err = strconv.ParseFloat(split[1], 64); err != nil {
				return nil, err
			}
			if y, err = strconv.ParseFloat(split[2], 64); err != nil {
				return nil, err
			}
			model.textures = append(model.textures, &vector.Vector{X: x, Y: y})
		} else if strings.HasPrefix(line, "vn") {
			var err error
			var v *vector.Vector3

			if v, err = parseVector(line[3:]); err != nil {
				return nil, err
			}
			model.normals = append(model.normals, v)
		} else if strings.HasPrefix(line, "v") {
			var err error
			var v *vector.Vector3

			if v, err = parseVector(line[2:]); err != nil {
				return nil, err
			}
			model.verticies = append(model.verticies, v)

		} else if strings.HasPrefix(line, "f") {
			split := strings.Split(line[2:], " ")

			face := &Face{}

			for _, vertConfig := range split {
				var vertexIndex int64
				var vertexNormalIndex int64
				var vertexTextureIndex int64
				var err error

				faceVertex := &FaceVertex{}

				indicies := strings.Split(vertConfig, "/")
				if vertexIndex, err = strconv.ParseInt(indicies[0], 10, 32); err != nil {
					return nil, err
				}
				faceVertex.Vertex = model.getVertex(int(vertexIndex))

				if indicies[1] != "" {
					if vertexTextureIndex, err = strconv.ParseInt(indicies[1], 10, 32); err != nil {
						return nil, err
					}
					faceVertex.Texture = model.getVertexTextures(int(vertexTextureIndex))
				}

				if indicies[2] != "" {
					if vertexNormalIndex, err = strconv.ParseInt(indicies[2], 10, 32); err != nil {
						return nil, err
					}
					faceVertex.Normal = model.getVertexNormal(int(vertexNormalIndex))
				}

				face.Verticies = append(face.Verticies, faceVertex)
			}
			model.Faces = append(model.Faces, face)
		}
	}

	return &model, nil
}

func (model *Model) getVertex(i int) *vector.Vector3 {
	return model.verticies[i-1]
}

func (model *Model) getVertexNormal(i int) *vector.Vector3 {
	return model.normals[i-1]
}

func (model *Model) getVertexTextures(i int) *vector.Vector {
	return model.textures[i-1]
}

// Parse a vector (space separated floats) from a string
func parseVector(str string) (*vector.Vector3, error) {
	var err error
	var x, y, z float64

	split := strings.Split(str, " ")
	if x, err = strconv.ParseFloat(split[0], 64); err != nil {
		return nil, err
	}
	if y, err = strconv.ParseFloat(split[1], 64); err != nil {
		return nil, err
	}
	if z, err = strconv.ParseFloat(split[2], 64); err != nil {
		return nil, err
	}

	return &vector.Vector3{X: x, Y: y, Z: z}, nil
}
