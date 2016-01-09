package assets

import (
	"fmt"
	"io/ioutil"

	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Manager struct {
	icons map[string]*sdl.Texture
	fonts map[string]*ttf.Font
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *Manager {
	ttf.Init()

	assetManager := Manager{
		icons: loadTextures(renderer, filepath.Join(directory, "icons")),
		fonts: loadFonts(filepath.Join(directory, "fonts")),
	}
	return &assetManager
}

func (assetManager *Manager) GetTexture(filename string) *sdl.Texture {
	return assetManager.icons[filename]
}

func loadFonts(directory string) map[string]*ttf.Font {
	fonts := make(map[string]*ttf.Font)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		fontPath := filepath.Join(directory, file.Name())

		font, err := ttf.OpenFont(fontPath, 30)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fonts[file.Name()] = font
	}

	return fonts
}

func loadTextures(renderer *sdl.Renderer, directory string) map[string]*sdl.Texture {
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
