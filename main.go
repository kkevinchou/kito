package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kkevinchou/kito/kito"
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
	var mode string = modeLocal
	if len(os.Args) > 1 {
		mode = strings.ToUpper(os.Args[1])
		if mode != modeLocal && mode != modeClient && mode != modeServer {
			panic(fmt.Sprintf("unexpected mode %s", mode))
		}
	}

	fmt.Println("starting game on mode:", mode)
	if mode == modeLocal {
		serverGame := kito.NewServerGame("_assets")

		go func() {
			serverGame.Start(input.NullInputPoller)
		}()

		game := kito.NewClientGame("_assets", "shaders")
		inputPoller := input.NewSDLInputPoller()
		game.Start(inputPoller.PollInput)
	} else if mode == modeClient {
		game := kito.NewClientGame("_assets", "shaders")
		inputPoller := input.NewSDLInputPoller()
		game.Start(inputPoller.PollInput)

	} else if mode == modeServer {
		game := kito.NewServerGame("_assets")
		game.Start(input.NullInputPoller)
	}

	sdl.Quit()
}
