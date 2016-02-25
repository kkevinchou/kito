package assets

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

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
