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
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/models"
	"github.com/kkevinchou/ant/lib/pathing"
	"github.com/kkevinchou/ant/logger"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	width                       = 800
	height                      = 600
	floorPanelDimension         = 1
	sensitivity         float64 = 0.3
)

var (
	textureMap map[string]uint32

	cameraX            float64 = -20
	cameraY            float64 = 15
	cameraZ            float64 = -5
	cameraRotationY    float64 = 90
	cameraRotationX    float64 = 20
	cameraRotationXMax float64 = 60
)

type Renderable interface {
	interfaces.Positionable
	GetRenderData() components.RenderData
}

type Renderables []Renderable

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	assetManager *lib.AssetManager
	renderables  Renderables
	textureMap   map[string]uint32
	modelMap     map[string]*models.Model
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

	sdl.SetRelativeMouseMode(true)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.ColorMaterial(gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE)

	gl.Enable(gl.LIGHTING)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.1, 0.1, 0.1, 1}
	diffuse := []float32{1, 1, 1, 1}
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

	oak, err := models.NewModel("_assets/obj/Oak_Green_01.obj")
	if err != nil {
		panic("Failed to load oak model")
	}
	renderSystem.modelMap = map[string]*models.Model{
		"oak": oak,
	}

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) CameraView(x, y int) {
	cameraRotationY += float64(x) * sensitivity
	cameraRotationX += float64(y) * sensitivity

	if cameraRotationX < -cameraRotationXMax {
		cameraRotationX = -cameraRotationXMax
	}

	if cameraRotationX > cameraRotationXMax {
		cameraRotationX = cameraRotationXMax
	}
}

func (r *RenderSystem) MoveCamera(v vector.Vector3) {
	forwardX, forwardY, forwardZ := forward()
	// Moving backwards
	forwardX *= -v.Z
	forwardY *= -v.Z
	forwardZ *= -v.Z

	rightX, rightY, rightZ := right()
	rightX *= -v.X
	rightY *= -v.X
	rightZ *= -v.X

	cameraX += forwardX + rightX
	cameraY += forwardY + rightY + v.Y
	cameraZ += forwardZ + rightZ
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
		renderData := renderable.GetRenderData()
		if !renderData.IsVisible() {
			continue
		}
		if rData, ok := renderData.(*components.TextureRenderData); ok {
			position := renderable.Position()
			texture := r.textureMap[rData.ID]
			drawQuad(texture, float32(position.X), float32(position.Y), float32(position.Z))
		} else if _, ok := renderData.(*components.ModelRenderData); ok {
			position := renderable.Position()
			fmt.Println(position)
		} else if _, ok := renderData.(*pathing.NavMeshRenderData); ok {
			var ok bool
			var navmesh *pathing.NavMesh
			if navmesh, ok = renderable.(*pathing.NavMesh); !ok {
				logger.Debug("FAILED TO CAST NAVMESH")
				continue
			}

			polygons := navmesh.Polygons()
			for i, polygon := range polygons {
				color := make([]float32, 3)
				gl.Begin(gl.POLYGON)
				gl.Normal3f(0, 1, 0)
				for _, point := range polygon.Points() {
					if i%2 == 0 {
						color[0], color[1], color[2] = 0, 0, 0
					} else {
						color[0], color[1], color[2] = 1, 1, 1
					}
					gl.Color3f(color[0], color[1], color[2])
					gl.Vertex3f(float32(point.X), float32(point.Y), float32(point.Z))
				}
				gl.End()
			}
		}
	}

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

func right() (float64, float64, float64) {
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

	v3 = v3.Normalize()

	return v3.X, v3.Y, v3.Z
}

func drawFloor() {
	width := 21
	height := 21
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			x := (i - int(math.Floor(float64(width)/2))) * floorPanelDimension
			y := (j - int(math.Floor(float64(height)/2))) * floorPanelDimension
			drawFloorPanel(float32(x), float32(y), (i+j)%2 == 0)
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

	halfDimension := float32(floorPanelDimension) / 2
	gl.Normal3f(0, 1, 0)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x-halfDimension, 0, z-halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x-halfDimension, 0, z+halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x+halfDimension, 0, z+halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x+halfDimension, 0, z-halfDimension)

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
