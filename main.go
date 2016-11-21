package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/ant"
	"github.com/kkevinchou/ant/directory"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 800
	height = 600
)

var window *sdl.Window

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

func setupDisplay() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	_, err = sdl.GL_CreateContext(window)
	if err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	setupDisplay()
	defer window.Destroy()

	game := ant.Game{}
	game.Init(window)
	directory := directory.GetDirectory()
	movementSystem := directory.MovementSystem()
	renderSystem := directory.RenderSystem()

	var event sdl.Event
	gameOver := false

	previousTime := time.Now()
	for gameOver != true {
		now := time.Now()
		delta := time.Since(previousTime)
		previousTime = now

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				gameOver = true
			case *sdl.MouseButtonEvent:
				if e.State == 0 { // Mouse Up
					// game.MoveAnt(float64(e.X), float64(e.Y))
					game.PlaceFood(float64(e.X), float64(e.Y))
				}
			case *sdl.KeyUpEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					gameOver = true
				}
			}
		}

		game.Update()
		movementSystem.Update(delta)
		renderSystem.Update(delta)
	}
	sdl.Quit()
}

// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/draw"
// 	_ "image/png"
// 	"log"
// 	"math"
// 	"os"
// 	"runtime"
// 	"time"

// 	"github.com/go-gl/gl/v2.1/gl"
// 	"github.com/kkevinchou/ant/lib/math/vector"
// 	"github.com/veandco/go-sdl2/sdl"
// )

// var (
// 	texture   uint32
// 	rotationX float32
// 	rotationY float32

// 	cameraX         float32 = 0
// 	cameraY         float32 = 0
// 	cameraZ         float32 = 8
// 	cameraRotationY float32 = 0
// 	cameraRotationX float32 = 0

// 	xDelta float32 = 0
// 	zDelta float32 = 0

// 	speed float32 = 1
// )

// var window *sdl.Window
// var renderer *sdl.Renderer
// var translationZ float32
// var translationX float32

// const (
// 	width  = 800
// 	height = 600
// )

// func init() {
// 	// We want to lock the main thread to this goroutine.  Otherwise,
// 	// SDL rendering will randomly panic
// 	//
// 	// For more details: https://github.com/golang/go/wiki/LockOSThread
// 	runtime.LockOSThread()
// }

// func toRadians(degrees float32) float32 {
// 	return math.Pi * degrees / 180.0
// }

// func forward() (float32, float32, float32) {
// 	xRadianAngle := -toRadians(cameraRotationX)
// 	if xRadianAngle < 0 {
// 		xRadianAngle += 2 * math.Pi
// 	}
// 	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
// 	if yRadianAngle < 0 {
// 		yRadianAngle += 2 * math.Pi
// 	}

// 	x := float32(math.Cos(float64(yRadianAngle)) * math.Cos(float64(xRadianAngle)))
// 	y := float32(math.Sin(float64(xRadianAngle)))
// 	z := -float32(math.Sin(float64(yRadianAngle)) * math.Cos(float64(xRadianAngle)))

// 	return x, y, z
// }

// func left() (float32, float32, float32) {
// 	xRadianAngle := -toRadians(cameraRotationX)
// 	if xRadianAngle < 0 {
// 		xRadianAngle += 2 * math.Pi
// 	}
// 	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
// 	if yRadianAngle < 0 {
// 		yRadianAngle += 2 * math.Pi
// 	}

// 	x, y, z := math.Cos(float64(yRadianAngle)), math.Sin(float64(xRadianAngle)), -math.Sin(float64(yRadianAngle))

// 	v1 := vector.Vector3{x, math.Abs(y), z}
// 	v2 := vector.Vector3{x, 0, z}
// 	v3 := v1.Cross(v2)

// 	fmt.Println(v3)
// 	if v3.X == 0 && v3.Y == 0 && v3.Z == 0 {
// 		v3 = vector.Vector3{v2.Z, 0, -v2.X}
// 	}
// 	return float32(v3.X), float32(v3.Y), float32(v3.Z)

// }

// func setupDisplay() {
// 	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
// 		panic(err)
// 	}

// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	_, err = sdl.GL_CreateContext(window)
// 	if err != nil {
// 		panic(err)
// 	}

// 	texture = newTexture("_assets/icons/F.png")

// 	setupScene()

// 	gameOver := false
// 	for !gameOver {
// 		time.Sleep(17 * time.Millisecond)

// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch e := event.(type) {
// 			case *sdl.QuitEvent:
// 				gameOver = true

// 			case *sdl.MouseMotionEvent:
// 				// cameraRotationY = (float32(e.Y) / height) * 360
// 				// cameraRotationX = (float32(e.X) / width) * 360
// 			case *sdl.KeyDownEvent:
// 				if e.Keysym.Sym == sdl.K_ESCAPE {
// 					gameOver = true
// 				}

// 				if e.Keysym.Sym == sdl.K_SPACE {
// 					cameraY += 1
// 				}

// 				if e.Keysym.Sym == sdl.K_w {
// 					x, y, z := forward()

// 					cameraX += x
// 					cameraY += y
// 					cameraZ += z
// 				}

// 				if e.Keysym.Sym == sdl.K_s {
// 					x, y, z := forward()

// 					cameraX -= x
// 					cameraY -= y
// 					cameraZ -= z
// 				}

// 				if e.Keysym.Sym == sdl.K_a {
// 					x, y, z := left()
// 					cameraX += x
// 					cameraY += y
// 					cameraZ += z
// 				}

// 				if e.Keysym.Sym == sdl.K_d {
// 					x, y, z := left()

// 					cameraX -= x
// 					cameraY -= y
// 					cameraZ -= z
// 				}

// 				if e.Keysym.Sym == sdl.K_LEFT {
// 					cameraRotationY -= 5
// 					if cameraRotationY < 0 {
// 						cameraRotationY += 360
// 					}
// 				}

// 				if e.Keysym.Sym == sdl.K_RIGHT {
// 					cameraRotationY += 5
// 					if cameraRotationY < 0 {
// 						cameraRotationY += 360
// 					}
// 				}

// 				if e.Keysym.Sym == sdl.K_UP {
// 					cameraRotationX -= 3
// 					if cameraRotationY >= 360 {
// 						cameraRotationY -= 360
// 					}
// 				}

// 				if e.Keysym.Sym == sdl.K_DOWN {
// 					cameraRotationX += 3
// 					if cameraRotationY >= 360 {
// 						cameraRotationY -= 360
// 					}
// 				}
// 			}
// 		}

// 		drawScene()

// 		sdl.GL_SwapWindow(window)
// 	}

// }

// func main() {
// 	setupDisplay()
// }

// func newTexture(file string) uint32 {
// 	imgFile, err := os.Open(file)
// 	if err != nil {
// 		log.Fatalf("texture %q not found on disk: %v\n", file, err)
// 	}
// 	img, _, err := image.Decode(imgFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rgba := image.NewRGBA(img.Bounds())
// 	if rgba.Stride != rgba.Rect.Size().X*4 {
// 		panic("unsupported stride")
// 	}
// 	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

// 	var texture uint32
// 	gl.Enable(gl.TEXTURE_2D)
// 	gl.GenTextures(1, &texture)
// 	gl.BindTexture(gl.TEXTURE_2D, texture)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
// 	gl.TexImage2D(
// 		gl.TEXTURE_2D,
// 		0,
// 		gl.RGBA,
// 		int32(rgba.Rect.Size().X),
// 		int32(rgba.Rect.Size().Y),
// 		0,
// 		gl.RGBA,
// 		gl.UNSIGNED_BYTE,
// 		gl.Ptr(rgba.Pix))

// 	return texture
// }

// func setupScene() {
// 	gl.Enable(gl.DEPTH_TEST)
// 	gl.Enable(gl.LIGHTING)

// 	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
// 	gl.ClearDepth(1)
// 	gl.DepthFunc(gl.LEQUAL)

// 	ambient := []float32{0.5, 0.5, 0.5, 1}
// 	diffuse := []float32{1, 1, 1, 1}
// 	lightPosition := []float32{-5, 5, 10, 0}
// 	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
// 	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
// 	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
// 	gl.Enable(gl.LIGHT0)

// 	gl.MatrixMode(gl.PROJECTION)
// 	gl.LoadIdentity()
// 	gl.Frustum(-0.5, 0.5, -0.5, 0.5, 1.0, 100.0)
// 	gl.PushMatrix()
// 	gl.MatrixMode(gl.MODELVIEW)
// 	gl.LoadIdentity()

// }

// func drawScene() {
// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

// 	gl.MatrixMode(gl.MODELVIEW)
// 	gl.LoadIdentity()
// 	drawCenter()
// 	gl.Rotatef(cameraRotationX, 1, 0, 0)
// 	gl.Rotatef(cameraRotationY, 0, 1, 0)
// 	gl.Translatef(-cameraX, -cameraY, -cameraZ)
// 	lightPosition := []float32{-5, 5, 10, 0}
// 	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])

// 	// yRadians := (cameraRotationY / 180 * 3.14159)
// 	// xDelta = float32(math.Sin(float64(yRadians)))
// 	// zDelta = float32(math.Cos(float64(yRadians)))
// 	// gl.Translatef(xDelta, 0, zDelta)
// 	// gl.Rotatef(cameraRotationY, 0, 1, 0)

// 	// gl.Rotatef(rotationX, 1, 0, 0)

// 	drawQuad()
// }

// func drawCenter() {

// 	// // gl.Disable(gl.LIGHTING)
// 	// gl.LineWidth(2.5)
// 	// gl.Color3f(256.0, 0.0, 0.0)
// 	// gl.Begin(gl.LINES)
// 	// gl.Vertex3f(0, 10, 0)
// 	// gl.Vertex3f(0, 0, 0)
// 	// gl.End()
// 	// // gl.Enable(gl.LIGHTING)

// 	gl.Disable(gl.LIGHTING)
// 	gl.Begin(gl.QUADS)

// 	gl.Color3f(256, 0.0, 0.0)
// 	gl.Vertex3f(-0.005, 0.005, -1)
// 	gl.Color3f(256, 0.0, 0.0)
// 	gl.Vertex3f(0.005, 0.005, -1)
// 	gl.Color3f(256, 0.0, 0.0)
// 	gl.Vertex3f(0.005, -0.005, -1)
// 	gl.Color3f(256, 0.0, 0.0)
// 	gl.Vertex3f(-0.005, -0.005, -1)

// 	gl.End()
// 	gl.Enable(gl.LIGHTING)
// }

// func drawQuad() {
// 	rotationX += 0.5
// 	rotationY += 0.5

// 	gl.BindTexture(gl.TEXTURE_2D, texture)

// 	gl.Color4f(1, 1, 1, 1)

// 	gl.Begin(gl.QUADS)

// 	// // FRONT
// 	gl.Normal3f(0, 0, 1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-10, -1, -10)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(10, -1, -10)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(10, -1, 10)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-10, -1, 10)

// 	gl.End()

// 	gl.Begin(gl.QUADS)

// 	// // FRONT
// 	gl.Normal3f(0, 0, 1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-1, 1, 1)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(1, 1, 1)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(1, -1, 1)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-1, -1, 1)

// 	// BACK
// 	gl.Normal3f(0, 0, -1)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-1, 1, -1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-1, -1, -1)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(1, -1, -1)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(1, 1, -1)

// 	// TOP
// 	gl.Normal3f(0, 1, 0)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-1, 1, -1) // A
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-1, 1, 1) // C
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(1, 1, 1) // D
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(1, 1, -1) // B

// 	// BOTTOM
// 	gl.Normal3f(0, -1, 0)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-1, -1, -1)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(1, -1, -1)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(1, -1, 1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-1, -1, 1)

// 	// RIGHT
// 	gl.Normal3f(1, 0, 0)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(1, -1, -1)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(1, 1, -1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(1, 1, 1)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(1, -1, 1)

// 	// LEFT
// 	gl.Normal3f(-1, 0, 0)
// 	gl.TexCoord2f(0, 1)
// 	gl.Vertex3f(-1, -1, -1)
// 	gl.TexCoord2f(1, 1)
// 	gl.Vertex3f(-1, -1, 1)
// 	gl.TexCoord2f(1, 0)
// 	gl.Vertex3f(-1, 1, 1)
// 	gl.TexCoord2f(0, 0)
// 	gl.Vertex3f(-1, 1, -1)

// 	gl.End()
// }
