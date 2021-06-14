package assets

import (
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
	"github.com/kkevinchou/kito/lib/types"
)

func loadTextures(directory string) map[string]*types.Texture {
	var subDirectories []string = []string{"images", "collada", "icons"}

	extensions := map[string]interface{}{
		".png": nil,
	}

	textureMap := map[string]*types.Texture{}
	assetMetaData := getAssetMetaData(directory, subDirectories, extensions)

	for _, metaData := range assetMetaData {
		textureID := newTexture(metaData.Path)
		textureMap[metaData.Name] = &types.Texture{ID: textureID}
	}

	return textureMap
}

func loadAnimatedModels(directory string) map[string]*animation.ModelSpecification {
	var subDirectories []string = []string{"collada"}

	extensions := map[string]interface{}{
		".dae": nil,
	}

	animationMap := map[string]*animation.ModelSpecification{}
	assetMetaData := getAssetMetaData(directory, subDirectories, extensions)

	for _, metaData := range assetMetaData {
		parsedCollada, err := collada.ParseCollada(metaData.Path)
		if err != nil {
			panic(err)
		}
		animationMap[metaData.Name] = parsedCollada
	}

	return animationMap
}
