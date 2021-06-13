package render

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/noise"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/veandco/go-sdl2/sdl"
)

type Camera interface {
	UpdateView(mgl64.Vec2)
	Position() mgl64.Vec3
	View() mgl64.Vec2
}

const (
	width               int32 = 1024
	height              int32 = 760
	depthBufferWidth    int32 = 2048
	depthBufferHeight   int32 = 1520
	floorPanelDimension       = 1
	renderDistance            = 500.0
	skyboxSize                = 500

	aspectRatio         = float32(width) / float32(height)
	fovy                = float32(90.0 / aspectRatio)
	near        float32 = 1
	far         float32 = 1000
)

var (
	noiseMap [][]float64 = noise.GenerateNoiseMap(100, 100)
)

type World interface {
	GetCamera() entities.Entity
}

type RenderSystem struct {
	renderer     *sdl.Renderer
	window       *sdl.Window
	camera       entities.Entity
	assetManager *assets.AssetManager
	world        World
	lights       []*Light
	shaders      map[string]*shaders.Shader
	skybox       *SkyBox
	floor        *Quad

	depthMapFBO  uint32
	depthTexture uint32

	entities []entities.Entity
}

func (r *RenderSystem) SetAssetManager(assetManager *assets.AssetManager) {
	r.assetManager = assetManager
}

func NewRenderSystem(world World) *RenderSystem {
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
		window: window,
		world:  world,
		skybox: NewSkyBox(300),
		floor:  NewQuad(quadZeroY),
	}

	sdl.SetRelativeMouseMode(false)
	sdl.GLSetSwapInterval(1)

	gl.ClearColor(1.0, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.MULTISAMPLE)

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

	depthDebugShader, err := shaders.NewShader("shaders/basictexture.vs", "shaders/depthvalue.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load depth debug shader %s", err))
	}

	basicShadowShader, err := shaders.NewShader("shaders/basicshadow.vs", "shaders/basicshadow.fs")
	if err != nil {
		panic(fmt.Sprintf("Failed to load basic shadow shader %s", err))
	}

	renderSystem.shaders = map[string]*shaders.Shader{
		"basic":       basicShader,
		"basicShadow": basicShadowShader,
		"skybox":      skyBoxShader,
		"model":       modelShader,
		"depth":       depthShader,
		"depthDebug":  depthDebugShader,
	}

	depthMapFBO, depthTexture, err := initializeShadowMap(depthBufferWidth, depthBufferHeight)
	if err != nil {
		panic(err)
	}

	renderSystem.depthMapFBO = depthMapFBO
	renderSystem.depthTexture = depthTexture

	return &renderSystem
}

func (s *RenderSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.RenderComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *RenderSystem) renderToDepthMap() mgl64.Mat4 {
	shader := s.shaders["depth"]

	lightPosition := mgl64.Vec3{0, 40, 40}
	lightViewQuaternion := mgl64.QuatRotate(mgl64.DegToRad(-30), mgl64.Vec3{1, 0, 0})

	orthoMatrix := mgl32.Ortho(-100, 100, -100, 100, near, far)
	// orthoMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, near, far)
	lightViewMatrix := mgl64.Translate3D(lightPosition.X(), lightPosition.Y(), lightPosition.Z()).Mul4(lightViewQuaternion.Mat4()).Inv()

	shader.Use()
	shader.SetUniformMat4("lightPerspective", orthoMatrix)
	shader.SetUniformMat4("view", utils.Mat4F64ToMat4F32(lightViewMatrix))
	shader.SetUniformMat4("model", mgl32.Ident4())

	gl.CullFace(gl.FRONT)
	gl.Viewport(0, 0, depthBufferWidth, depthBufferHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.depthMapFBO)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	s.renderScene(orthoMatrix, lightPosition, lightViewQuaternion, lightViewMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.CullFace(gl.BACK)

	return mgl64.Ortho(-100, 100, -100, 100, float64(near), float64(far)).Mul4(lightViewMatrix)
}

func (s *RenderSystem) Update(delta time.Duration) {
	// render depth map
	lightViewMatrix := s.renderToDepthMap()

	// regular render
	gl.Viewport(0, 0, width, height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	camera := s.world.GetCamera()
	componentContainer := camera.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent

	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, near, far)
	s.renderScene(projectionMatrix, transformComponent.Position, transformComponent.ViewQuaternion, lightViewMatrix)

	// orthoMatrix := mgl32.Ortho(-10, 10, -10, 10, near, far)
	// s.renderSceneTest(orthoMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	s.window.GLSwap()
}

func (s *RenderSystem) renderScene(perspectiveMatrix mgl32.Mat4, viewerPosition mgl64.Vec3, viewerQuaternion mgl64.Quat, lightSpaceMatrix mgl64.Mat4) {
	// We use the inverse to move the universe in the opposite direction of where the camera is looking
	viewerViewQuaternion := utils.QuatF64ToQuatF32(viewerQuaternion)
	viewerViewMatrix := viewerViewQuaternion.Inverse().Mat4()

	floorModelMatrix := createModelMatrix(
		mgl32.Scale3D(100, 100, 100),
		mgl32.Ident4(),
		mgl32.Ident4(),
	)

	viewTranslationMatrix := mgl32.Translate3D(float32(-viewerPosition.X()), float32(-viewerPosition.Y()), float32(-viewerPosition.Z()))
	viewMatrix := viewerViewMatrix.Mul4(viewTranslationMatrix)

	vPosition := mgl32.Vec3{float32(viewerPosition[0]), float32(viewerPosition[1]), float32(viewerPosition[2])}

	drawTextureToQuad(s.shaders["depthDebug"], s.depthTexture, mgl32.Translate3D(0, 10, 0), viewMatrix, perspectiveMatrix)
	drawSkyBox(
		s.skybox,
		s.shaders["skybox"],
		s.assetManager.GetTexture("front"),
		s.assetManager.GetTexture("top"),
		s.assetManager.GetTexture("left"),
		s.assetManager.GetTexture("right"),
		s.assetManager.GetTexture("bottom"),
		s.assetManager.GetTexture("back"),
		mgl32.Ident4(),
		viewerViewMatrix,
		perspectiveMatrix,
	)
	drawMesh(s.floor, s.shaders["basicShadow"], floorModelMatrix, viewMatrix, perspectiveMatrix, vPosition, lightSpaceMatrix, s.depthTexture)

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
				drawAnimatedMesh(
					animationComponent.AnimatedModel.Mesh,
					animationComponent.AnimationTransforms,
					s.assetManager.GetTexture("diffuse"),
					s.shaders["model"],
					meshModelMatrix,
					viewMatrix,
					perspectiveMatrix,
					vPosition,
				)
			}
		} else if _, ok := renderData.(*components.BlockRenderData); ok {
		}
	}
}

func (s *RenderSystem) renderSceneTest(perspectiveMatrix mgl32.Mat4) {
	shader := s.shaders["depth"]

	lightViewMatrix := mgl32.QuatRotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}).Mat4().Inv()

	shader.Use()
	shader.SetUniformMat4("lightPerspective", perspectiveMatrix)
	shader.SetUniformMat4("view", lightViewMatrix)
	shader.SetUniformMat4("model", mgl32.Translate3D(0, -2, 0).Mul4(mgl32.Scale3D(10, 10, 10)))

	q := NewQuad(quadZeroY)
	gl.BindVertexArray(q.GetVAO())
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
