package assets

import (
	"fmt"

	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/assets/loaders"
	"github.com/kkevinchou/kito/lib/textures"
	"github.com/veandco/go-sdl2/ttf"
)

type AssetManager struct {
	textures       map[string]*textures.Texture
	animatedModels map[string]*animation.ModelSpecification
}

func NewAssetManager(directory string) *AssetManager {
	ttf.Init()

	assetManager := AssetManager{
		textures:       loaders.LoadTextures(directory),
		animatedModels: loaders.LoadAnimatedModels(directory),
	}

	return &assetManager
}

func (a *AssetManager) GetTexture(name string) *textures.Texture {
	if _, ok := a.textures[name]; !ok {
		panic(fmt.Sprintf("could not find texture %s", name))
	}
	return a.textures[name]
}

func (a *AssetManager) GetAnimatedModel(name string) *animation.ModelSpecification {
	if _, ok := a.animatedModels[name]; !ok {
		panic(fmt.Sprintf("could not find animated model %s", name))
	}
	return a.animatedModels[name]
}
