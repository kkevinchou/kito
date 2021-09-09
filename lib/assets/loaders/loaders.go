package loaders

import (
	"fmt"
	"strings"

	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/assets/loaders/collada"
	"github.com/kkevinchou/kito/lib/assets/loaders/gltextures"
	"github.com/kkevinchou/kito/lib/textures"
	"github.com/kkevinchou/kito/lib/utils"
)

func LoadTextures(directory string) map[string]*textures.Texture {
	var subDirectories []string = []string{"images", "collada", "icons"}

	extensions := map[string]interface{}{
		".png": nil,
	}

	textureMap := map[string]*textures.Texture{}
	fileMetaData := utils.GetFileMetaData(directory, subDirectories, extensions)

	for _, metaData := range fileMetaData {
		textureID := gltextures.NewTexture(metaData.Path)
		textureMap[metaData.Name] = &textures.Texture{ID: textureID}
	}

	return textureMap
}

func LoadAnimatedModels(directory string) map[string]*animation.ModelSpecification {
	var subDirectories []string = []string{"collada"}

	extensions := map[string]interface{}{
		".dae": nil,
	}

	animationMap := map[string]*animation.ModelSpecification{}
	fileMetaData := utils.GetFileMetaData(directory, subDirectories, extensions)

	for _, metaData := range fileMetaData {
		if strings.HasPrefix(metaData.Name, "_") {
			continue
		}
		parsedCollada, err := collada.ParseCollada(metaData.Path)
		if err != nil {
			fmt.Println("failed to parse collada for", metaData.Path, ", error:", err)
			continue
		}
		animationMap[metaData.Name] = parsedCollada
	}

	return animationMap
}
