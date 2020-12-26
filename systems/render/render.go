package render

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/lib/models"
	"github.com/kkevinchou/kito/lib/noise"
	"github.com/kkevinchou/kito/lib/pathing"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Viewer interface {
	UpdateView(vector.Vector)
	Position() vector.Vector3
	View() vector.Vector
}

const (
	width               int32 = 1024
	height              int32 = 760
	floorPanelDimension       = 1
	renderDistance            = 500.0
	skyboxSize                = 500

	aspectRatio = float32(width) / float32(height)
	fovy        = float32(90.0 / aspectRatio)
)

var (
	textureMap map[string]uint32

	noiseMap [][]float64 = noise.GenerateNoiseMap(100, 100)
)

type Game interface {
}

type Renderable interface {
	types.Positionable
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
	shaders      map[string]*shaders.Shader
	skybox       *SkyBox
	floor        *Quad

	animator *animation.Animator
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
		skybox:       NewSkyBox(300),
		floor:        NewQuad(nil),
	}

	sdl.SetRelativeMouseMode(false)
	sdl.GLSetSwapInterval(1)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	// gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)

	parsedCollada, err := collada.ParseCollada("_assets/collada/model.dae")
	if err != nil {
		panic(err)
	}
	animatedModel := animation.NewAnimatedModel(parsedCollada, 50, 3)
	renderSystem.animator = animation.NewAnimator(animatedModel, parsedCollada.Animation)

	_ = initFont()
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
	cowboyTexture := newTexture("_assets/collada/diffuse.png")
	renderSystem.textureMap = map[string]uint32{
		"high-grass":     highGrassTexture,
		"mushroom-gills": mushroomGilsTexture,
		"worker":         workerTexture,
		"light":          lightTexture,

		// skybox
		"left":   leftTexture,
		"right":  rightTexture,
		"front":  frontTexture,
		"back":   backTexture,
		"top":    topTexture,
		"bottom": bottomTexture,
		"cowboy": cowboyTexture,
	}

	oak, err := models.NewModel("_assets/obj/Oak_Green_01.obj")
	if err != nil {
		panic(fmt.Sprintf("Failed to load oak model %s", err))
	}

	torus, err := models.NewModel("_assets/obj/torus.obj")
	if err != nil {
		panic(fmt.Sprintf("Failed to load oak model %s", err))
	}

	land, err := models.NewModel("_assets/obj/land.obj")
	if err != nil {
		panic(fmt.Sprintf("Failed to load oak model %s", err))
	}

	renderSystem.modelMap = map[string]*models.Model{
		"oak":   oak,
		"torus": torus,
		"land":  land,
	}

	basicShader, err := shaders.NewShader("shaders/basic.vs", "shaders/basic.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load basic shader %s", err))
	}

	skyBoxShader, err := shaders.NewShader("shaders/skybox.vs", "shaders/skybox.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load skybox shader %s", err))
	}

	modelShader, err := shaders.NewShader("shaders/model.vs", "shaders/model.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load model shader %s", err))
	}

	renderSystem.shaders = map[string]*shaders.Shader{
		"basic":  basicShader,
		"skybox": skyBoxShader,
		"model":  modelShader,
	}

	return &renderSystem
}

func (r *RenderSystem) Register(renderable Renderable) {
	r.renderables = append(r.renderables, renderable)
}

var updateOnce bool

func (r *RenderSystem) Update(delta time.Duration) {
	if !updateOnce {
		r.animator.Update(delta)
		updateOnce = true
	}
	// r.viewer.UpdateView(vector.Vector{X: 5, Y: 0})
	viewerPosition := r.viewer.Position()
	viewerView := r.viewer.View()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	verticalViewRotationMatrix := mgl32.QuatRotate(mgl32.DegToRad(float32(viewerView.X)), mgl32.Vec3{1, 0, 0}).Mat4()
	horizontalViewRotationMatrix := mgl32.QuatRotate(mgl32.DegToRad(float32(viewerView.Y)), mgl32.Vec3{0, 1, 0}).Mat4()

	floorModelMatrix := createModelMatrix(
		// mgl32.Scale3D(10, 10, 10),
		mgl32.Scale3D(1, 1, 1),
		mgl32.Ident4(),
		mgl32.Ident4(),
	)
	floorModelMatrix = horizontalViewRotationMatrix.Mul4(floorModelMatrix)

	viewTranslationMatrix := mgl32.Translate3D(float32(-viewerPosition.X), float32(-viewerPosition.Y), float32(-viewerPosition.Z))
	viewMatrix := verticalViewRotationMatrix.Mul4(viewTranslationMatrix)

	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, 1, 1000)

	meshModelMatrix := createModelMatrix(
		// mgl32.Scale3D(3, 3, 3),
		mgl32.Ident4(),
		horizontalViewRotationMatrix,
		// horizontalViewRotationMatrix.Mul4(mgl32.QuatRotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}).Mat4()),
		mgl32.Ident4(),
	)
	drawMesh(r, r.textureMap["cowboy"], r.shaders["model"], meshModelMatrix, viewMatrix, projectionMatrix, viewerPosition)
	drawSkyBox(r.skybox, r.shaders["skybox"], r.textureMap, mgl32.Ident4(), verticalViewRotationMatrix.Mul4(horizontalViewRotationMatrix), projectionMatrix)
	drawQuad(r.floor, r.shaders["basic"], floorModelMatrix, viewMatrix, projectionMatrix, viewerPosition)

	for _, renderable := range r.renderables {
		renderData := renderable.GetRenderData()
		if !renderData.IsVisible() {
			continue
		}
		if rData, ok := renderData.(*components.TextureRenderData); ok {
			_ = rData
			// position := renderable.Position()
			// texture := r.textureMap[rData.ID]
			// drawCube(texture, float32(position.X), float32(position.Y), float32(position.Z), 1, true)
		} else if rData, ok := renderData.(*components.ItemRenderData); ok {
			_ = rData
			// position := renderable.Position()
			// texture := r.textureMap[rData.ID]
			// drawCube(texture, float32(position.X), float32(position.Y), float32(position.Z), 1, true)
		} else if _, ok := renderData.(*components.ModelRenderData); ok {
		} else if _, ok := renderData.(*pathing.NavMeshRenderData); ok {
			// if navMesh, ok := renderable.(*pathing.NavMesh); ok {
			// 	RenderNavMesh(navMesh)
			// } else {
			// 	panic("FAILED TO CAST NAVMESH")
			// }
		}

		// temp code, force rendering oak tree
		// r.renderModel(r.modelMap["oak"], vector.Vector3{X: 0, Y: 0, Z: 0})
		// r.renderModel(r.modelMap["land"], vector.Vector3{X: 0, Y: 0, Z: 0})

		// width := float32(len(noiseMap[0]))
		// height := float32(len(noiseMap))
		// var edgeLength float32 = 1
		// RenderNoiseMap(noiseMap, -(width*edgeLength)/2, -(height*edgeLength)/2, edgeLength)
	}

	gl.UseProgram(0)

	r.window.GLSwap()
}

// // x, y represents the x,y coordinate on the window. The output is a 3d position in world coordinates
// func (r *RenderSystem) GetWorldPoint(x, y float64) vector.Vector3 {
// 	// Get the projection matrix
// 	pMatrixValues := make([]float32, 16)
// 	gl.GetFloatv(gl.PROJECTION_MATRIX, &pMatrixValues[0])

// 	// Get the model view matrix
// 	mvMatrixValues := make([]float32, 16)
// 	gl.GetFloatv(gl.MODELVIEW, &mvMatrixValues[0])

// 	mvMatrix := matrix.Mat4FromValues(mvMatrixValues)
// 	pMatrix := matrix.Mat4FromValues(pMatrixValues)

// 	// Convert the screen coordinate to normalised device coordinates
// 	NDCPoint := mgl32.Vec4{(2.0*float32(x))/width - 1, 1 - (2.0*float32(y))/height, -1, 1}
// 	worldPoint := pMatrix.Mul4(mvMatrix).Inv().Mul4x1(NDCPoint)

// 	// Normalize on W
// 	worldPoint = mgl32.Vec4{worldPoint[0] / worldPoint[3], worldPoint[1] / worldPoint[3], worldPoint[2] / worldPoint[3], 1}

// 	// Extract the 3D vector
// 	worldPointVector := vector.Vector3{X: float64(worldPoint[0]), Y: float64(worldPoint[1]), Z: float64(worldPoint[2])}

// 	return worldPointVector
// }
