package render

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	width  = 800
	height = 600
)

var (
	textureMap map[string]uint32

	cameraX         float64 = 0
	cameraY         float64 = 2
	cameraZ         float64 = 8
	cameraRotationY float64 = 0
	cameraRotationX float64 = 0
)

type Renderable interface {
	interfaces.Positionable
	Texture() string
	Visible() bool
}

type Renderables []Renderable

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	assetManager *lib.AssetManager
	renderables  Renderables
	textureMap   map[string]uint32
}

func initFont() *ttf.Font {
	ttf.Init()

	font, err := ttf.OpenFont("_assets/fonts/courier_new.ttf", 30)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Font not found")
	}

	return font
}

func NewRenderSystem(window *sdl.Window, assetManager *lib.AssetManager) *RenderSystem {
	renderSystem := RenderSystem{
		assetManager: assetManager,
		window:       window,
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.ColorMaterial(gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE)
	gl.Enable(gl.COLOR_MATERIAL)

	gl.Enable(gl.LIGHTING)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{0.5, 0.5, 0.5, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-0.5, 0.5, -0.5, 0.5, 1.0, 100.0)
	gl.PushMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	_ = initFont()
	highGrassTexture := newTexture("_assets/icons/high-grass.png")
	mushroomGilsTexture := newTexture("_assets/icons/mushroom-gills.png")
	workerTexture := newTexture("_assets/icons/worker.png")
	renderSystem.textureMap = map[string]uint32{
		"high-grass":     highGrassTexture,
		"mushroom-gills": mushroomGilsTexture,
		"worker":         workerTexture,
	}

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) CameraView(x, y int) {
	normalizedX := (float64(x) - (float64(width) / 2)) / width
	normalizedY := (float64(y) - (float64(height) / 2)) / height
	cameraRotationY = normalizedX * 180
	cameraRotationX = normalizedY * 180
}

func (r *RenderSystem) MoveCamera(v vector.Vector3) {
	forwardX, forwardY, forwardZ := forward()
	if v.Z == 1 {
		forwardX *= -1
		forwardY *= -1
		forwardZ *= -1
	}

	leftX, leftY, leftZ := left()
	if v.Z == 1 {
		leftX *= -1
		leftY *= -1
		leftZ *= -1
	}

	cameraX += forwardX + leftX
	cameraY += forwardY + leftY
	cameraZ += forwardZ + leftZ
}

func (r *RenderSystem) Update(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Rotatef(float32(cameraRotationX), 1, 0, 0)
	gl.Rotatef(float32(cameraRotationY), 0, 1, 0)
	gl.Translatef(float32(-cameraX), float32(-cameraY), float32(-cameraZ))
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])

	for _, renderable := range r.renderables {
		if !renderable.Visible() {
			continue
		}
		position := renderable.Position()
		texture := r.textureMap[renderable.Texture()]
		drawQuad(texture, float32(position.X), float32(position.Y), float32(position.Z))
	}
	drawFloor()

	sdl.GL_SwapWindow(r.window)
}

func toRadians(degrees float64) float64 {
	return degrees / 180 * math.Pi
}

func forward() (float64, float64, float64) {
	xRadianAngle := -toRadians(cameraRotationX)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x := math.Cos(yRadianAngle) * math.Cos(xRadianAngle)
	y := math.Sin(xRadianAngle)
	z := -math.Sin(yRadianAngle) * math.Cos(xRadianAngle)

	return x, y, z
}

func left() (float64, float64, float64) {
	xRadianAngle := -toRadians(cameraRotationX)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x, y, z := math.Cos(yRadianAngle), math.Sin(xRadianAngle), -math.Sin(yRadianAngle)

	v1 := vector.Vector3{x, math.Abs(y), z}
	v2 := vector.Vector3{x, 0, z}
	v3 := v1.Cross(v2)

	if v3.X == 0 && v3.Y == 0 && v3.Z == 0 {
		v3 = vector.Vector3{v2.Z, 0, -v2.X}
	}
	return v3.X, v3.Y, v3.Z
}

func drawFloor() {
	width := 11
	height := 11
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			drawFloorPanel(float32(i)*2-float32(width)+1, float32(j)*2-float32(height)+1, (i+j)%2 == 0)
		}
	}
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

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
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

func drawFloorPanel(x, z float32, black bool) {
	color := make([]float32, 3)
	if black {
		color[0] = 0
		color[1] = 0
		color[2] = 0
	} else {
		color[0] = 1
		color[1] = 1
		color[2] = 1
	}

	gl.Begin(gl.QUADS)

	gl.Normal3f(0, 0, 1)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(-1+x, 0, -1+z)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(1+x, 0, -1+z)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(1+x, 0, 1+z)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(-1+x, 0, 1+z)

	gl.End()
}

func drawQuad(texture uint32, x, y, z float32) {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	// // FRONT
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z+0.5)

	// BACK
	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+1, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z+-0.5)

	// BOTTOM
	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+0, z+0.5)

	// RIGHT
	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+0.5, y+1, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z-0.5)

	// LEFT
	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)

}
