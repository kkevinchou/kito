package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type AnimationMetaDataJson struct {
	Fps       int `json:"fps"`
	NumFrames int `json:"num_frames"`
}

type AnimationDefinition struct {
	numFrames int
	fps       int
	frames    []*sdl.Texture
}

func (a *AnimationDefinition) NumFrames() int {
	return a.numFrames
}

func (a *AnimationDefinition) Fps() int {
	return a.fps
}

func (a *AnimationDefinition) GetFrame(frame int) *sdl.Texture {
	return a.frames[frame]
}

func loadAnimations(directory string, renderer *sdl.Renderer) map[string]*AnimationDefinition {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	animations := map[string]*AnimationDefinition{}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		animations[file.Name()] = loadAnimation(filepath.Join(directory, file.Name()), renderer)
	}

	return animations
}

func loadAnimation(directory string, renderer *sdl.Renderer) *AnimationDefinition {
	metaDataFile, err := os.Open(filepath.Join(directory, "meta.json"))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bytes, err := ioutil.ReadAll(metaDataFile)
	var metaData AnimationMetaDataJson
	json.Unmarshal(bytes, &metaData)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	frames := make([]*sdl.Texture, metaData.NumFrames)

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".png" {
			continue
		}

		framePath := filepath.Join(directory, file.Name())
		texture, err := img.LoadTexture(renderer, framePath)
		if err != nil {
			log.Fatal(err.Error())
		}

		frameIndex, err := strconv.Atoi(strings.Split(file.Name(), ".")[0])
		frames[frameIndex] = texture
	}

	a := AnimationDefinition{
		frames:    frames,
		numFrames: metaData.NumFrames,
		fps:       metaData.Fps,
	}

	return &a
}
