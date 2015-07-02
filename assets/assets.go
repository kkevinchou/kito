package assets

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"io/ioutil"
	"path/filepath"
)

type Manager struct {
	assets map[string]*sdl.Texture
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *Manager {
	assetManager := Manager{
		assets: loadAssets(renderer, directory),
	}
	return &assetManager
}

func (assetManager *Manager) GetTexture(filename string) *sdl.Texture {
	return assetManager.assets[filename]
}

func loadAssets(renderer *sdl.Renderer, directory string) map[string]*sdl.Texture {
	m := make(map[string]*sdl.Texture)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		imagePath := filepath.Join(directory, file.Name())

		image, err := img.Load(imagePath)
		if err != nil {
			fmt.Println("Failed to load \"%s\": %s", imagePath, err)
			continue
		}

		texture, err := renderer.CreateTextureFromSurface(image)
		if err != nil {
			fmt.Println("Failed to create texture \"%s\": %s", imagePath, err)
			continue
		}

		m[file.Name()] = texture
	}

	return m
}

// loadAssets(renderer, "assets/icons")
