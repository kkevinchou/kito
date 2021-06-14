package loaders

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	_ "image/png"
)

type AssetMetaData struct {
	Name string
	Path string
}

func getAssetMetaData(directory string, subDirectories []string, extensions map[string]interface{}) []AssetMetaData {
	var subPaths []string
	for _, subDir := range subDirectories {
		subPaths = append(subPaths, path.Join(directory, subDir))
	}

	var metaDataCollection []AssetMetaData

	for _, subDir := range subPaths {
		files, err := os.ReadDir(subDir)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		for _, file := range files {
			extension := filepath.Ext(file.Name())
			if _, ok := extensions[extension]; !ok {
				continue
			}

			path := filepath.Join(subDir, file.Name())
			name := file.Name()[0 : len(file.Name())-len(extension)]

			metaDataCollection = append(metaDataCollection, AssetMetaData{Name: name, Path: path})
		}
	}

	return metaDataCollection
}
