package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AssetManager struct {
	icons map[string]*sdl.Texture
	fonts map[string]*ttf.Font
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *AssetManager {
	return nil
	// ttf.Init()

	// assetManager := AssetManager{
	// 	icons:      loadTextures(filepath.Join(directory, "icons"), renderer),
	// 	fonts:      loadFonts(filepath.Join(directory, "fonts")),
	// 	animations: loadAnimations(filepath.Join(directory, "animations"), renderer),
	// }

	// return &assetManager
}

func (assetManager *AssetManager) GetTexture(filename string) *sdl.Texture {
	return assetManager.icons[filename]
}

func (assetManager *AssetManager) GetFont(filename string) *ttf.Font {
	return assetManager.fonts[filename]
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

		font, err := ttf.OpenFont(fontPath, 24)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fonts[file.Name()] = font
	}

	return fonts
}

func loadTextures(directory string, renderer *sdl.Renderer) map[string]*sdl.Texture {
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

		extensionLength := len(filepath.Ext(file.Name()))
		m[file.Name()[0:len(file.Name())-extensionLength]] = texture
	}

	return m
}
