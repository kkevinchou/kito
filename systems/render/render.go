package render

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	texture uint32

	cameraX         float32 = 0
	cameraY         float32 = 50
	cameraZ         float32 = 0
	cameraRotationY float32 = 0
	cameraRotationX float32 = 90
)

type Renderable interface {
	interfaces.Positionable
	Render(*lib.AssetManager, *sdl.Renderer)
	UpdateRenderComponent(time.Duration)
	GetRenderPriority() int
	GetY() float64
}

type Renderables []Renderable

func (r Renderables) Len() int {
	return len(r)
}

func (r Renderables) Less(i, j int) bool {
	if r[i].GetRenderPriority() == r[j].GetRenderPriority() {
		return r[i].GetY() < r[j].GetY()
	}
	return r[i].GetRenderPriority() < r[j].GetRenderPriority()
}

func (r Renderables) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	assetManager *lib.AssetManager
	renderables  Renderables
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
	texture = newTexture("_assets/icons/F.png")

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) Update(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Rotatef(cameraRotationX, 1, 0, 0)
	gl.Rotatef(cameraRotationY, 0, 1, 0)
	gl.Translatef(-cameraX, -cameraY, -cameraZ)
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	drawFloor()

	for _, renderable := range r.renderables {
		position := renderable.Position()
		drawQuad(texture, float32(position.X), 0, float32(position.Y))
	}

	sdl.GL_SwapWindow(r.window)
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

	// // BACK
	// gl.Normal3f(0, 0, -1)
	// gl.TexCoord2f(0, 0)
	// gl.Vertex3f(-1, 1, -1)
	// gl.TexCoord2f(0, 1)
	// gl.Vertex3f(-1, -1, -1)
	// gl.TexCoord2f(1, 1)
	// gl.Vertex3f(1, -1, -1)
	// gl.TexCoord2f(1, 1)
	// gl.Vertex3f(1, 0, -1)

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5) // A
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+1, z+0.5) // C
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+1, z+0.5) // D
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z+-0.5) // B

	// // BOTTOM
	// gl.Normal3f(0, -1, 0)
	// gl.TexCoord2f(0, 1)
	// gl.Vertex3f(-1, -1, -1)
	// gl.TexCoord2f(1, 1)
	// gl.Vertex3f(1, -1, -1)
	// gl.TexCoord2f(1, 0)
	// gl.Vertex3f(1, -1, 1)
	// gl.TexCoord2f(0, 0)
	// gl.Vertex3f(-1, -1, 1)

	// // RIGHT
	// gl.Normal3f(1, 0, 0)
	// gl.TexCoord2f(1, 1)
	// gl.Vertex3f(1, -1, -1)
	// gl.TexCoord2f(1, 0)
	// gl.Vertex3f(1, 1, -1)
	// gl.TexCoord2f(0, 0)
	// gl.Vertex3f(1, 1, 1)
	// gl.TexCoord2f(0, 1)
	// gl.Vertex3f(1, -1, 1)

	// // LEFT
	// gl.Normal3f(-1, 0, 0)
	// gl.TexCoord2f(0, 1)
	// gl.Vertex3f(-1, -1, -1)
	// gl.TexCoord2f(1, 1)
	// gl.Vertex3f(-1, -1, 1)
	// gl.TexCoord2f(1, 0)
	// gl.Vertex3f(-1, 1, 1)
	// gl.TexCoord2f(0, 0)
	// gl.Vertex3f(-1, 1, -1)

	gl.End()
}
