package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Animation struct {
	numFrames       int
	secondsPerFrame float64
	frames          []*sdl.Texture
	metaData        MetaData
}

type MetaData struct {
	Fps       int `json:"fps"`
	NumFrames int `json:"num_frames"`
	Name      string
}

type Manager struct {
	icons      map[string]*sdl.Texture
	fonts      map[string]*ttf.Font
	animations map[string]*Animation
}

func NewAssetManager(renderer *sdl.Renderer, directory string) *Manager {
	ttf.Init()

	assetManager := Manager{
		icons:      loadTextures(filepath.Join(directory, "icons"), renderer),
		fonts:      loadFonts(filepath.Join(directory, "fonts")),
		animations: loadAnimations(filepath.Join(directory, "animations"), renderer),
	}

	return &assetManager
}

func (assetManager *Manager) GetTexture(filename string) *sdl.Texture {
	return assetManager.icons[filename]
}

func (assetManager *Manager) GetFont(filename string) *ttf.Font {
	return assetManager.fonts[filename]
}

func (assetManager *Manager) GetAnimation(animation string, frame int) *sdl.Texture {
	return assetManager.animations[animation].frames[frame]
}

func (assetManager *Manager) GetAnimationMetaData(animation string) MetaData {
	return assetManager.animations[animation].metaData
}

func loadAnimations(directory string, renderer *sdl.Renderer) map[string]*Animation {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	animations := map[string]*Animation{}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		animations[file.Name()] = loadAnimation(filepath.Join(directory, file.Name()), renderer)
	}

	return animations
}

func loadAnimation(directory string, renderer *sdl.Renderer) *Animation {
	metaDataFile, err := os.Open(filepath.Join(directory, "meta.json"))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bytes, err := ioutil.ReadAll(metaDataFile)
	var metaData MetaData
	json.Unmarshal(bytes, &metaData)

	metaData.Name = filepath.Base(directory)

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

	a := Animation{}
	a.numFrames = metaData.NumFrames
	a.secondsPerFrame = 1 / float64(metaData.Fps)
	a.frames = frames
	a.metaData = metaData

	return &a
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

		font, err := ttf.OpenFont(fontPath, 15)
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

		m[file.Name()] = texture
	}

	return m
}

// loadAssets(renderer, "assets/icons")
