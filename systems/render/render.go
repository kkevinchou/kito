package render

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/noise"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/veandco/go-sdl2/sdl"
)

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
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type RenderSystem struct {
	*base.BaseSystem
	renderer     *sdl.Renderer
	window       *sdl.Window
	camera       entities.Entity
	assetManager *assets.AssetManager
	world        World
	lights       []*Light
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
		panic(fmt.Sprintf("Failed to init SDL %s", err))
	}

	// Enable hints for multisampling which allows opengl to use the default
	// multisampling algorithms implemented by the OpenGL rasterizer
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create window %s", err))
	}

	_, err = window.GLCreateContext()
	if err != nil {
		panic(fmt.Sprintf("Failed to create context %s", err))
	}

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to init OpenGL %s", err))
	}

	renderSystem := RenderSystem{
		BaseSystem: &base.BaseSystem{},
		window:     window,
		world:      world,
		skybox:     NewSkyBox(300),
		floor:      NewQuad(quadZeroY),
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
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()

	shader := shaderManager.GetShaderProgram("depth")

	lightPosition := mgl64.Vec3{0, 40, 40}
	lightOrientation := mgl64.QuatRotate(mgl64.DegToRad(-30), mgl64.Vec3{1, 0, 0})

	orthoMatrix := mgl32.Ortho(-100, 100, -100, 100, near, far)
	// orthoMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, near, far)
	lightViewMatrix := mgl64.Translate3D(lightPosition.X(), lightPosition.Y(), lightPosition.Z()).Mul4(lightOrientation.Mat4()).Inv()

	shader.Use()
	shader.SetUniformMat4("lightPerspective", orthoMatrix)
	shader.SetUniformMat4("view", utils.Mat4F64ToMat4F32(lightViewMatrix))
	shader.SetUniformMat4("model", mgl32.Ident4())

	gl.CullFace(gl.FRONT)
	gl.Viewport(0, 0, depthBufferWidth, depthBufferHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.depthMapFBO)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	s.renderScene(orthoMatrix, lightPosition, lightOrientation, lightViewMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.CullFace(gl.BACK)

	return mgl64.Ortho(-100, 100, -100, 100, float64(near), float64(far)).Mul4(lightViewMatrix)
}

func (s *RenderSystem) Update(delta time.Duration) {
}

func (s *RenderSystem) Render(delta time.Duration) {
	singleton := s.world.GetSingleton()
	if singleton.CameraID == 0 {
		fmt.Println("camera not found in Render()")
		return
	}

	camera, err := s.world.GetEntityByID(singleton.CameraID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// render depth map
	lightViewMatrix := s.renderToDepthMap()

	// regular render
	gl.Viewport(0, 0, width, height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	componentContainer := camera.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent

	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(fovy), aspectRatio, near, far)
	s.renderScene(projectionMatrix, transformComponent.Position, transformComponent.Orientation, lightViewMatrix)

	// orthoMatrix := mgl32.Ortho(-10, 10, -10, 10, near, far)
	// s.renderSceneTest(orthoMatrix)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	s.window.GLSwap()
}

func (s *RenderSystem) renderScene(perspectiveMatrix mgl32.Mat4, viewerPosition mgl64.Vec3, viewerQuaternion mgl64.Quat, lightSpaceMatrix mgl64.Mat4) {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()

	// We use the inverse to move the universe in the opposite direction of where the camera is looking
	viewerOrientation := utils.QuatF64ToQuatF32(viewerQuaternion)
	viewerViewMatrix := viewerOrientation.Inverse().Mat4()

	floorModelMatrix := createModelMatrix(
		mgl32.Scale3D(100, 100, 100),
		mgl32.Ident4(),
		mgl32.Ident4(),
	)

	viewTranslationMatrix := mgl32.Translate3D(float32(-viewerPosition.X()), float32(-viewerPosition.Y()), float32(-viewerPosition.Z()))
	viewMatrix := viewerViewMatrix.Mul4(viewTranslationMatrix)

	vPosition := mgl32.Vec3{float32(viewerPosition[0]), float32(viewerPosition[1]), float32(viewerPosition[2])}

	drawTextureToQuad(shaderManager.GetShaderProgram("depthDebug"), s.depthTexture, mgl32.Translate3D(0, 10, 0), viewMatrix, perspectiveMatrix)
	drawSkyBox(
		s.skybox,
		shaderManager.GetShaderProgram("skybox"),
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
	drawMesh(s.floor, shaderManager.GetShaderProgram("basicShadow"), floorModelMatrix, viewMatrix, perspectiveMatrix, vPosition, lightSpaceMatrix, s.depthTexture)

	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		renderData := componentContainer.RenderComponent.GetRenderData()
		entityPosition := componentContainer.TransformComponent.Position
		rotation := componentContainer.TransformComponent.Orientation

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
					shaderManager.GetShaderProgram("model"),
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
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	shader := shaderManager.GetShaderProgram("depth")

	lightViewMatrix := mgl32.QuatRotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0}).Mat4().Inv()

	shader.Use()
	shader.SetUniformMat4("lightPerspective", perspectiveMatrix)
	shader.SetUniformMat4("view", lightViewMatrix)
	shader.SetUniformMat4("model", mgl32.Translate3D(0, -2, 0).Mul4(mgl32.Scale3D(10, 10, 10)))

	q := NewQuad(quadZeroY)
	gl.BindVertexArray(q.GetVAO())
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
