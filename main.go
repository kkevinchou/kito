package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	_ "net/http/pprof"

	"github.com/kkevinchou/kito/kito"
	"github.com/kkevinchou/kito/kito/config"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

const (
	modeLocal  string = "LOCAL"
	modeClient string = "CLIENT"
	modeServer string = "SERVER"
)

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("failed to load config.json, using defaults: %s\n", err)
	} else {
		configBytes, err := io.ReadAll(configFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = configFile.Close(); err != nil {
			fmt.Println(err)
			return
		}

		var configSettings config.Config
		err = json.Unmarshal(configBytes, &configSettings)
		if err != nil {
			fmt.Println(err)
			return
		}

		loadConfig(configSettings)
	}

	var mode string = modeClient
	if len(os.Args) > 1 {
		mode = strings.ToUpper(os.Args[1])
		if mode != modeLocal && mode != modeClient && mode != modeServer {
			panic(fmt.Sprintf("unexpected mode %s", mode))
		}
	}

	// if mode == modeClient {
	// 	f, err := os.Create("cpuprofile")
	// 	if err != nil {
	// 		log.Fatal("could not create CPU profile: ", err)
	// 	}
	// 	defer f.Close() // error handling omitted for example
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		log.Fatal("could not start CPU profile: ", err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	if settings.PProfEnabled {
		go func() {
			if mode == modeClient {
				log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", settings.PProfClientPort), nil))
			} else {
				log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", settings.PProfServerPort), nil))
			}
		}()
	}

	fmt.Println("starting game on mode:", mode)
	if mode == modeClient {
		game := kito.NewClientGame("_assets", "shaders")
		platform := input.NewSDLPlatform()
		game.Start(platform.PollInput)
	} else if mode == modeServer {
		game := kito.NewServerGame("_assets")
		game.Start(input.NullInputPoller)
	}

	sdl.Quit()
}

func loadConfig(configSettings config.Config) {
	settings.Host = configSettings.ServerIP
	settings.Port = configSettings.ServerPort
	settings.Width = configSettings.Width
	settings.Height = configSettings.Height
}
