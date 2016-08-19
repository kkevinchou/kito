package assets

import (
	"path/filepath"

	"github.com/kkevinchou/ant/lib"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Manager struct {
	icons      map[string]*sdl.Texture
	fonts      map[string]*ttf.Font
	animations map[string]*animation.AnimationDefinition
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *Manager {
	ttf.Init()

	assetManager := Manager{
		icons:      loadTextures(filepath.Join(directory, "icons"), renderer),
		fonts:      loadFonts(filepath.Join(directory, "fonts")),
		animations: lib.LoadAnimations(filepath.Join(directory, "animations"), renderer),
	}

	return &assetManager
}

func (assetManager *Manager) GetTexture(filename string) *sdl.Texture {
	return assetManager.icons[filename]
}

func (assetManager *Manager) GetFont(filename string) *ttf.Font {
	return assetManager.fonts[filename]
}

func (assetManager *Manager) GetAnimation(animation string) *lib.AnimationDefinition {
	return assetManager.animations[animation]
}
