package animation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type Animation struct {
	// State
	currentFrame int
	timeCounter  float64

	// Data
	numFrames       int
	secondsPerFrame float64
}

type MetaData struct {
	Fps       int `json:"fps"`
	NumFrames int `json:"num_frames"`
}

func Load(directory string) *Animation {
	metaDataFile, err := os.Open(filepath.Join(directory, "meta.json"))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bytes, err := ioutil.ReadAll(metaDataFile)
	var metaData MetaData
	json.Unmarshal(bytes, &metaData)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		framePath := filepath.Join(directory, file.Name())
		frame, err := img.Load(framePath)
		if err != nil {
			log.Fatal(err.Error())
		}

		texture, err := renderer.CreateTextureFromSurface(frame)
		if err != nil {
			log.Fatal(err.Error())
		}

		_ = texture
	}

	a := Animation{}
	a.currentFrame = 0
	a.numFrames = 50
	a.secondsPerFrame = 1 / float64(24)

	return &a
}

func (a *Animation) GetFrame() int {
	return a.currentFrame
}

func (a *Animation) Update(delta time.Duration) *sdl.Texture {
	a.timeCounter += delta.Seconds()
	for a.timeCounter >= a.secondsPerFrame {
		a.timeCounter -= a.secondsPerFrame
		a.currentFrame += 1
	}
	a.currentFrame = a.currentFrame % a.numFrames
	return nil
}
