package assets

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AssetManager struct {
	textures map[string]*Texture
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *AssetManager {
	_ = initFont()

	assetManager := AssetManager{
		textures: loadTextures(directory, renderer),
	}

	return &assetManager
}

func (a *AssetManager) GetTexture(name string) *Texture {
	if _, ok := a.textures[name]; !ok {
		panic(fmt.Sprintf("could not find texture %s", name))
	}
	return a.textures[name]
}

func initFont() *ttf.Font {
	ttf.Init()

	font, err := ttf.OpenFont("_assets/fonts/courier_new.ttf", 30)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Font not found")
	}

	return font
}
