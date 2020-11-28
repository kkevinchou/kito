package render

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/common/enums"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/math/matrix"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/lib/models"
	"github.com/kkevinchou/kito/lib/pathing"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Viewer interface {
	Position() vector.Vector3
	View() vector.Vector
}

const (
	width               = 1280
	height              = 1024
	floorPanelDimension = 1
	renderDistance      = 300.0
)

var (
	textureMap map[string]uint32
)

type Game interface {
	GetGameMode() enums.GameMode
}

type Renderable interface {
	interfaces.Positionable
	GetRenderData() components.RenderData
}

type Renderables []Renderable

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	viewer       Viewer
	assetManager *lib.AssetManager
	renderables  Renderables
	textureMap   map[string]uint32
	modelMap     map[string]*models.Model
	game         Game
	lights       []*Light
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

func NewRenderSystem(game Game, assetManager *lib.AssetManager, viewer Viewer) *RenderSystem {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf("Failed to init SDL", err))
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create window", err))
	}

	_, err = window.GLCreateContext()
	if err != nil {
		panic(fmt.Sprintf("Failed to create context", err))
	}

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to init OpenGL %s", err))
	}

	renderSystem := RenderSystem{
		assetManager: assetManager,
		window:       window,
		viewer:       viewer,
		game:         game,
	}

	sdl.SetRelativeMouseMode(false)
	sdl.GLSetSwapInterval(1)

	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.ColorMaterial(gl.FRONT, gl.AMBIENT_AND_DIFFUSE)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	light0 := NewLight(gl.LIGHT0)
	renderSystem.lights = append(renderSystem.lights, light0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-0.5, 0.5, -0.375, 0.375, 1.0, renderDistance)
	gl.PushMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	_ = initFont()
	// skyboxTexture := newCubeMap("_assets/images/skybox.png")
	highGrassTexture := newTexture("_assets/icons/high-grass.png")
	mushroomGilsTexture := newTexture("_assets/icons/mushroom-gills.png")
	workerTexture := newTexture("_assets/icons/worker.png")
	lightTexture := newTexture("_assets/icons/light.png")
	leftTexture := newTexture("_assets/images/left.png")
	rightTexture := newTexture("_assets/images/right.png")
	frontTexture := newTexture("_assets/images/front.png")
	backTexture := newTexture("_assets/images/back.png")
	topTexture := newTexture("_assets/images/top.png")
	bottomTexture := newTexture("_assets/images/bottom.png")
	renderSystem.textureMap = map[string]uint32{
		"high-grass":     highGrassTexture,
		"mushroom-gills": mushroomGilsTexture,
		"worker":         workerTexture,
		"light":          lightTexture,
		// "skybox":         skyboxTexture,
		"left":   leftTexture,
		"right":  rightTexture,
		"front":  frontTexture,
		"back":   backTexture,
		"top":    topTexture,
		"bottom": bottomTexture,
	}

	oak, err := models.NewModel("_assets/obj/Oak_Green_01.obj")
	if err != nil {
		panic(fmt.Sprintf("Failed to load oak model %s", err))
	}
	renderSystem.modelMap = map[string]*models.Model{
		"oak": oak,
	}

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

func (r *RenderSystem) Update(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if r.game.GetGameMode() == enums.GameModePlaying {
		viewerPosition := r.viewer.Position()
		viewerView := r.viewer.View()

		// Set up the Model View matrix.  Based on how much the camera has moved,
		// translate the entire world
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Rotatef(float32(viewerView.X), 1, 0, 0)
		gl.Rotatef(float32(viewerView.Y), 0, 1, 0)

		// draw skybox without consideration for camera translation
		skyboxSize := 50
		drawSkyBox(r.textureMap, float32(0), float32(-skyboxSize), float32(0), 100, false)

		gl.Translatef(float32(-viewerPosition.X), float32(-viewerPosition.Y), float32(-viewerPosition.Z))

		for _, light := range r.lights {
			light.Update(delta)
			texture := r.textureMap["light"]
			position := light.Position()
			drawCube(texture, float32(position.X), float32(position.Y), float32(position.Z), 1, false)
		}

		for _, renderable := range r.renderables {
			renderData := renderable.GetRenderData()
			if !renderData.IsVisible() {
				continue
			}
			if rData, ok := renderData.(*components.TextureRenderData); ok {
				position := renderable.Position()
				texture := r.textureMap[rData.ID]
				drawCube(texture, float32(position.X), float32(position.Y), float32(position.Z), 1, true)
			} else if rData, ok := renderData.(*components.ItemRenderData); ok {
				position := renderable.Position()
				texture := r.textureMap[rData.ID]
				drawCube(texture, float32(position.X), float32(position.Y), float32(position.Z), 1, true)
			} else if _, ok := renderData.(*components.ModelRenderData); ok {
			} else if _, ok := renderData.(*pathing.NavMeshRenderData); ok {
				if navMesh, ok := renderable.(*pathing.NavMesh); !ok {
					RenderNavMesh(navMesh)
				} else {
					panic("FAILED TO CAST NAVMESH")
				}
			}
		}
	} else {
		fmt.Println("Editor Mode")
	}

	r.renderModel(r.modelMap["oak"])
	r.window.GLSwap()
}

func RenderNavMesh(navMesh *pathing.NavMesh) {
	polygons := navMesh.Polygons()
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

// x, y represents the x,y coordinate on the window. The output is a 3d position in world coordinates
func (r *RenderSystem) GetWorldPoint(x, y float64) vector.Vector3 {
	// Get the projection matrix
	pMatrixValues := make([]float32, 16)
	gl.GetFloatv(gl.PROJECTION_MATRIX, &pMatrixValues[0])

	// Get the model view matrix
	mvMatrixValues := make([]float32, 16)
	gl.GetFloatv(gl.MODELVIEW, &mvMatrixValues[0])

	mvMatrix := matrix.Mat4FromValues(mvMatrixValues)
	pMatrix := matrix.Mat4FromValues(pMatrixValues)

	// Convert the screen coordinate to normalised device coordinates
	NDCPoint := mgl32.Vec4{(2.0*float32(x))/width - 1, 1 - (2.0*float32(y))/height, -1, 1}
	worldPoint := pMatrix.Mul4(mvMatrix).Inv().Mul4x1(NDCPoint)

	// Normalize on W
	worldPoint = mgl32.Vec4{worldPoint[0] / worldPoint[3], worldPoint[1] / worldPoint[3], worldPoint[2] / worldPoint[3], 1}

	// Extract the 3D vector
	worldPointVector := vector.Vector3{X: float64(worldPoint[0]), Y: float64(worldPoint[1]), Z: float64(worldPoint[2])}

	return worldPointVector
}
