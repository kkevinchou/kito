package render

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/noise"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Camera interface {
	UpdateView(mgl64.Vec2)
	Position() mgl64.Vec3
	View() mgl64.Vec2
}

const (
	width               int32 = 1024
	height              int32 = 760
	floorPanelDimension       = 1
	renderDistance            = 500.0
	skyboxSize                = 500

	aspectRatio         = float32(width) / float32(height)
	fovy                = float32(90.0 / aspectRatio)
	near        float32 = 1
	far         float32 = 1000
)

var (
	textureMap map[string]uint32

	noiseMap [][]float64 = noise.GenerateNoiseMap(100, 100)
)

type World interface {
	GetCamera() entities.Entity
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	camera       entities.Entity
	assetManager *lib.AssetManager
	textureMap   map[string]uint32
	world        World
	lights       []*Light
	shaders      map[string]*shaders.Shader
	skybox       *SkyBox
	floor        *Quad

	entities []entities.Entity
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

func NewRenderSystem(world World, assetManager *lib.AssetManager) *RenderSystem {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf("Failed to init SDL", err))
	}

	// Enable hints for multisampling which allows opengl to use the default
	// multisampling algorithms implemented by the OpenGL rasterizer
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

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
		world:        world,
		skybox:       NewSkyBox(300),
		floor:        NewQuad(quadZeroY),
	}

	sdl.SetRelativeMouseMode(false)
	sdl.GLSetSwapInterval(1)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.MULTISAMPLE)

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

	depthShader, err := shaders.NewShader("shaders/depth.vs", "shaders/depth.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load depth shader %s", err))
	}

	depthValueShader, err := shaders.NewShader("shaders/basictexture.vs", "shaders/depthvalue.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load texture test shader %s", err))
	}

	renderSystem.shaders = map[string]*shaders.Shader{
		"basic":      basicShader,
		"skybox":     skyBoxShader,
		"model":      modelShader,
		"depth":      depthShader,
		"depthValue": depthValueShader,
	}

	asdfdepthMapFBO, asdfdepthTexture, err = initializeShadowMap(width, height)
	if err != nil {
		panic(err)
	}

	return &renderSystem
}

func (s *RenderSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.RenderComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

var asdfdepthMapFBO, asdfdepthTexture uint32

func (s *RenderSystem) renderToDepthMap() {

	// var near float32 = 1
	// var far float32 = 7.5

	// lightProjectionMatrix := mgl32.Ortho(-10, 10, -10, 10, near, far)
	// lightViewMatrix := utils.QuatF64ToQuatF32(utils.QuatLookAt(mgl64.Vec3{-2, 4, -1}, mgl64.Vec3{0, -1, -1}, mgl64.Vec3{0, 1, 0})).Mat4()
	// lightSpaceMatrix := lightProjectionMatrix.Mul4(lightViewMatrix)

	// shader.Use()
	// shader.SetUniformMat4("lightSpaceMatrix", lightSpaceMatrix)
	// shader.SetUniformMat4("model", mgl32.Translate3D(0, 3, 0))

	gl.Viewport(0, 0, width, height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, asdfdepthMapFBO)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	orthoMatrix := mgl32.Ortho(-10, 10, -10, 10, near, far)
	s.renderScene(orthoMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (s *RenderSystem) Update(delta time.Duration) {
	// render depth map
	s.renderToDepthMap()

	// regular render
	gl.Viewport(0, 0, width, height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// projectionMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, near, far)
	orthoMatrix := mgl32.Ortho(-10, 10, -10, 10, near, far)
	s.renderSceneTest(orthoMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	s.window.GLSwap()
}

func (s *RenderSystem) renderSceneTest(perspectiveMatrix mgl32.Mat4) {
	shader := s.shaders["depth"]

	// lightViewMatrix := utils.QuatF64ToQuatF32(utils.QuatLookAt(mgl64.Vec3{-2, 4, -1}, mgl64.Vec3{0, -1, -1}, mgl64.Vec3{0, 1, 0})).Mat4()
	// lightSpaceMatrix := perspectiveMatrix.Mul4(mgl32.QuatRotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}).Mat4())
	// fmt.Println(mgl32.QuatRotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}).Rotate(mgl32.Vec3{0, 0, -1}))

	// shader.Use()
	// shader.SetUniformMat4("lightSpaceMatrix", lightSpaceMatrix)
	// shader.SetUniformMat4("model", mgl32.Translate3D(0, 2, 0))
	// _ = lightSpaceMatrix

	shader.Use()
	shader.SetUniformMat4("lightSpaceMatrix", perspectiveMatrix)
	shader.SetUniformMat4("model", mgl32.Translate3D(0, 0, -1).Mul4(mgl32.Scale3D(10, 10, 10)))

	// fmt.Println(lightSpaceMatrix.Mul4x1(mgl32.Vec4{-0.5, 0, -0.5, 1.0}))

	q := NewQuad(quadZeroZ)
	gl.BindVertexArray(q.GetVAO())
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (s *RenderSystem) renderScene(perspectiveMatrix mgl32.Mat4) {
	camera := s.world.GetCamera()
	componentContainer := camera.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent

	// We use the inverse to move the universe in the opposite direction of where the camera is looking
	cameraViewQuaternion := utils.QuatF64ToQuatF32(transformComponent.ViewQuaternion)
	cameraViewMatrix := cameraViewQuaternion.Inverse().Mat4()
	cameraPosition := transformComponent.Position

	floorModelMatrix := createModelMatrix(
		mgl32.Scale3D(100, 100, 100),
		mgl32.Ident4(),
		mgl32.Ident4(),
	)

	viewTranslationMatrix := mgl32.Translate3D(float32(-cameraPosition.X()), float32(-cameraPosition.Y()), float32(-cameraPosition.Z()))
	viewMatrix := cameraViewMatrix.Mul4(viewTranslationMatrix)

	vPosition := mgl32.Vec3{float32(cameraPosition[0]), float32(cameraPosition[1]), float32(cameraPosition[2])}

	drawTextureToQuad(s.shaders["depthValue"], asdfdepthTexture, mgl32.Translate3D(0, 10, 0), viewMatrix, perspectiveMatrix)
	drawSkyBox(s.skybox, s.shaders["skybox"], s.textureMap, mgl32.Ident4(), cameraViewMatrix, perspectiveMatrix)
	drawMesh(s.floor, s.shaders["basic"], floorModelMatrix, viewMatrix, perspectiveMatrix, vPosition)

	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		renderData := componentContainer.RenderComponent.GetRenderData()
		entityPosition := componentContainer.TransformComponent.Position
		rotation := componentContainer.TransformComponent.ViewQuaternion

		if !renderData.IsVisible() {
			continue
		}

		if rData, ok := renderData.(*components.ModelRenderData); ok {
			if rData.Animated {

				// accounting for blender change of axis
				xr := mgl32.QuatRotate(mgl32.DegToRad(90), mgl32.Vec3{1, 0, 0}).Mat4()
				yr := mgl32.QuatRotate(mgl32.DegToRad(180), mgl32.Vec3{0, 1, 0}).Mat4()

				meshModelMatrix := createModelMatrix(
					mgl32.Ident4(),
					utils.QuatF64ToQuatF32(rotation).Mat4().Mul4(xr.Mul4(yr)),
					mgl32.Translate3D(float32(entityPosition.X()), float32(entityPosition.Y()), float32(entityPosition.Z())),
				)

				animationComponent := componentContainer.AnimationComponent
				drawAnimatedMesh(animationComponent.AnimatedModel.Mesh, animationComponent.AnimationTransforms, s.textureMap["cowboy"], s.shaders["model"], meshModelMatrix, viewMatrix, perspectiveMatrix, vPosition)
			}
		} else if _, ok := renderData.(*components.BlockRenderData); ok {
		}
	}
}
