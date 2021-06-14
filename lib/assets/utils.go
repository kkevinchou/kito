package assets

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"path/filepath"

	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type AssetMetaData struct {
	Name string
	Path string
}

func newTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	// is vertically flipped if directly read into opengl texture
	nrgba := imaging.FlipV(img)

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), nrgba, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
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
