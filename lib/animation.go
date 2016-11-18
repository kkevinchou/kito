package lib

import (
	"time"

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
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type AnimationState struct {
	animationDefinition *AnimationDefinition
	currentFrame        int
	timeCounter         float64
	secondsPerFrame     float64
}

func (a *AnimationState) GetCurrentFrame() *sdl.Texture {
	return a.animationDefinition.GetFrame(a.currentFrame)
}

func (a *AnimationState) Update(delta time.Duration) {
	a.timeCounter += delta.Seconds()
	for a.timeCounter >= a.secondsPerFrame {
		a.timeCounter -= a.secondsPerFrame
		a.currentFrame += 1
	}
	a.currentFrame = a.currentFrame % a.animationDefinition.NumFrames()
}

func CreateStateFromAnimationDef(animationDefinition *AnimationDefinition) *AnimationState {
	secondsPerFrame := 1 / float64(animationDefinition.Fps())
	return &AnimationState{
		animationDefinition: animationDefinition,
		secondsPerFrame:     secondsPerFrame,
	}
}

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

type AssetManager struct {
	icons      map[string]*sdl.Texture
	fonts      map[string]*ttf.Font
	animations map[string]*AnimationDefinition
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *AssetManager {
	ttf.Init()

	assetManager := AssetManager{
		icons:      loadTextures(filepath.Join(directory, "icons"), renderer),
		fonts:      loadFonts(filepath.Join(directory, "fonts")),
		animations: loadAnimations(filepath.Join(directory, "animations"), renderer),
	}

	return &assetManager
}

func (assetManager *AssetManager) GetTexture(filename string) *sdl.Texture {
	return assetManager.icons[filename]
}

func (assetManager *AssetManager) GetFont(filename string) *ttf.Font {
	return assetManager.fonts[filename]
}

func (assetManager *AssetManager) GetAnimation(animation string) *AnimationDefinition {
	return assetManager.animations[animation]
}

func loadFonts(directory string) map[string]*ttf.Font {
	fonts := make(map[string]*ttf.Font)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		fontPath := filepath.Join(directory, file.Name())

		font, err := ttf.OpenFont(fontPath, 24)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fonts[file.Name()] = font
	}

	return fonts
}

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
