package parser

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
}

type FaceVertex struct {
	Vertex  *vector.Vector3
	Texture *vector.Vector
	Normal  *vector.Vector3
}

type Model struct {
	Normals   []*vector.Vector3
	Textures  []*vector.Vector
	Verticies []*vector.Vector3
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
		var lineType string
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if len(line) == 1 {
			return nil, errors.New("line of length 1 unexpected")
		}

		if line[0] == 'v' {
			lineType = "v"
			if line[1] == 'n' {
				lineType = "vn"
			} else if line[1] == 't' {
				lineType = "vt"
			}
		} else if line[0] == 'f' {
			lineType = "f"
		} else {
			continue
		}

		switch lineType {
		case "v": // VERTICES
			var err error
			var x, y, z float64

			split := strings.Split(strings.TrimSpace(line[2:]), " ")
			if x, err = strconv.ParseFloat(split[0], 64); err != nil {
				return nil, err
			}
			if y, err = strconv.ParseFloat(split[1], 64); err != nil {
				return nil, err
			}
			if z, err = strconv.ParseFloat(split[2], 64); err != nil {
				return nil, err
			}
			model.Verticies = append(model.Verticies, &vector.Vector3{X: x, Y: y, Z: z})
		case "vn": // NORMALS
			var err error
			var x, y, z float64

			split := strings.Split(strings.TrimSpace(line[3:]), " ")
			if x, err = strconv.ParseFloat(split[0], 64); err != nil {
				return nil, err
			}
			if y, err = strconv.ParseFloat(split[1], 64); err != nil {
				return nil, err
			}
			if z, err = strconv.ParseFloat(split[2], 64); err != nil {
				return nil, err
			}
			model.Normals = append(model.Normals, &vector.Vector3{X: x, Y: y, Z: z})
		case "vt": // TEXTURE VERTICES
			var err error
			var x, y float64

			split := strings.Split(strings.TrimSpace(line[3:]), " ")
			if x, err = strconv.ParseFloat(split[0], 64); err != nil {
				return nil, err
			}
			if y, err = strconv.ParseFloat(split[1], 64); err != nil {
				return nil, err
			}
			model.Textures = append(model.Textures, &vector.Vector{X: x, Y: y})
		case "f": // INDICES
			split := strings.Split(strings.TrimSpace(line[2:]), " ")

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
	return model.Verticies[i-1]
}

func (model *Model) getVertexNormal(i int) *vector.Vector3 {
	return model.Normals[i-1]
}

func (model *Model) getVertexTextures(i int) *vector.Vector {
	return model.Textures[i-1]
}
