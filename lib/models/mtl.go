package models

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl64"
)

type Material struct {
	Ambient  *mgl64.Vec3
	Diffuse  *mgl64.Vec3
	Specular *mgl64.Vec3
}

func parseMaterials(file string) (map[string]*Material, error) {
	mtlFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer mtlFile.Close()

	var currentMaterial *Material
	materials := map[string]*Material{}
	scanner := bufio.NewScanner(mtlFile)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "newmtl") {
			split := strings.Split(line, " ")
			currentMaterial = &Material{}
			materials[split[1]] = currentMaterial
		} else if strings.HasPrefix(line, "Ka") {
			var err error
			var v *mgl64.Vec3
			if v, err = parseVector(line[3:]); err != nil {
				return nil, err
			}

			currentMaterial.Ambient = v
		} else if strings.HasPrefix(line, "Kd") {
			var err error
			var v *mgl64.Vec3
			if v, err = parseVector(line[3:]); err != nil {
				return nil, err
			}

			currentMaterial.Diffuse = v
		} else if strings.HasPrefix(line, "Ks") {
			var err error
			var v *mgl64.Vec3
			if v, err = parseVector(line[3:]); err != nil {
				return nil, err
			}

			currentMaterial.Specular = v
		}
	}

	return materials, nil
}
