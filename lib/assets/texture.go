package assets

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	textureFolders []string = []string{"images", "collada", "icons"}
)

type Texture struct {
	ID uint32
}

func loadTextures(directory string, renderer *sdl.Renderer) map[string]*Texture {
	textureMap := map[string]*Texture{}

	var subDirectories []string
	for _, subDir := range textureFolders {
		subDirectories = append(subDirectories, path.Join(directory, subDir))
	}

	for _, subDir := range subDirectories {
		files, err := os.ReadDir(subDir)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		for _, file := range files {
			extension := filepath.Ext(file.Name())
			if extension != ".png" {
				continue
			}

			imagePath := filepath.Join(subDir, file.Name())
			textureID := newTexture(imagePath)
			extensionLength := len(extension)

			textureMap[file.Name()[0:len(file.Name())-extensionLength]] = &Texture{ID: textureID}
		}
	}

	return textureMap
}
