package render

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  int32 = 1024
	height int32 = 760

	aspectRatio         = float64(width) / float64(height)
	fovx        float64 = 90
	near        float64 = 1
	far         float64 = 500

	// shadow map parameters
	shadowMapDimension   int     = 8000
	shadowDistanceFactor float64 = 0.3 // proportion of view fustrum to include in shadow cuboid
)

var (
	fovy float64 = mgl64.RadToDeg(2 * math.Atan(math.Tan(mgl64.DegToRad(fovx)/2)/aspectRatio))
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type RenderSystem struct {
	*base.BaseSystem
	window       *sdl.Window
	assetManager *assets.AssetManager
	world        World
	skybox       *SkyBox
	floor        *Quad
	shadowMap    *ShadowMap

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
		panic(fmt.Sprintf("failed to create window %s", err))
	}

	_, err = window.GLCreateContext()
	if err != nil {
		panic(fmt.Sprintf("failed to create context %s", err))
	}

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("failed to init OpenGL %s", err))
	}

	renderSystem := RenderSystem{
		BaseSystem: &base.BaseSystem{},
		window:     window,
		world:      world,
		skybox:     NewSkyBox(600),
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

	renderSystem.shadowMap, err = NewShadowMap(shadowMapDimension, shadowMapDimension, far*shadowDistanceFactor)
	if err != nil {
		panic(fmt.Sprintf("failed to create shadow map %s", err))
	}

	return &renderSystem
}

func (s *RenderSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.RenderComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *RenderSystem) GetCameraTransform() *components.TransformComponent {
	singleton := s.world.GetSingleton()
	if singleton.CameraID == 0 {
		return nil
	}
	camera, err := s.world.GetEntityByID(singleton.CameraID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	componentContainer := camera.GetComponentContainer()
	return componentContainer.TransformComponent
}

func (s *RenderSystem) Render(delta time.Duration) {
	transformComponent := s.GetCameraTransform()
	if transformComponent == nil {
		fmt.Println("camera not found in Render()")
		return
	}

	// configure camera viewer context
	viewerViewMatrix := transformComponent.Orientation.Mat4()
	viewTranslationMatrix := mgl64.Translate3D(transformComponent.Position.X(), transformComponent.Position.Y(), transformComponent.Position.Z())

	cameraViewerContext := ViewerContext{
		Position:    transformComponent.Position,
		Orientation: transformComponent.Orientation,

		InverseViewMatrix: viewTranslationMatrix.Mul4(viewerViewMatrix).Inv(),
		ProjectionMatrix:  mgl64.Perspective(mgl64.DegToRad(fovy), aspectRatio, near, far),
	}

	// configure light viewer context
	modelSpaceFrustumPoints := CalculateFrustumPoints(transformComponent.Position, transformComponent.Orientation, near, far, fovy, aspectRatio, shadowDistanceFactor)
	lightOrientation := mgl64.QuatRotate(mgl64.DegToRad(-150), mgl64.Vec3{1, 0, 0})
	lightPosition, lightProjectionMatrix := ComputeDirectionalLightProps(lightOrientation.Mat4(), modelSpaceFrustumPoints)
	lightViewMatrix := mgl64.Translate3D(lightPosition.X(), lightPosition.Y(), lightPosition.Z()).Mul4(lightOrientation.Mat4()).Inv()

	lightViewerContext := ViewerContext{
		Position:    lightPosition,
		Orientation: lightOrientation,

		InverseViewMatrix: lightViewMatrix,
		ProjectionMatrix:  lightProjectionMatrix,
	}

	lightContext := LightContext{
		DirectionalLightDir: lightOrientation.Rotate(mgl64.Vec3{0, 0, -1}),
		LightSpaceMatrix:    lightProjectionMatrix.Mul4(lightViewMatrix),
	}

	s.renderToDepthMap(lightViewerContext, lightContext)
	s.renderToDisplay(cameraViewerContext, lightContext)

	s.window.GLSwap()
}

func (s *RenderSystem) renderToDepthMap(viewerContext ViewerContext, lightContext LightContext) {
	defer resetGLRenderSettings()
	s.shadowMap.Prepare()
	s.renderScene(viewerContext, lightContext)
}

func (s *RenderSystem) renderToDisplay(viewerContext ViewerContext, lightContext LightContext) {
	defer resetGLRenderSettings()

	gl.Viewport(0, 0, width, height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.renderScene(viewerContext, lightContext)
}

// renderScene renders a scene from the perspective of a viewer
func (s *RenderSystem) renderScene(viewerContext ViewerContext, lightContext LightContext) {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()

	// render a debug shadow map for viewing
	drawHUDTextureToQuad(viewerContext, shaderManager.GetShaderProgram("depthDebug"), s.shadowMap.DepthTexture(), 0.4)

	drawSkyBox(
		viewerContext,
		s.skybox,
		shaderManager.GetShaderProgram("skybox"),
		s.assetManager.GetTexture("front"),
		s.assetManager.GetTexture("top"),
		s.assetManager.GetTexture("left"),
		s.assetManager.GetTexture("right"),
		s.assetManager.GetTexture("bottom"),
		s.assetManager.GetTexture("back"),
	)

	floorModelMatrix := createModelMatrix(
		mgl64.Scale3D(500, 500, 500),
		mgl64.Ident4(),
		mgl64.Ident4(),
	)

	drawThingy(
		viewerContext,
		lightContext,
		s.shadowMap,
		shaderManager.GetShaderProgram("basicShadow"),
		s.assetManager.GetTexture("default"),
		s.floor,
		floorModelMatrix,
	)

	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		renderData := componentContainer.RenderComponent.GetRenderData()
		entityPosition := componentContainer.TransformComponent.Position
		orientation := componentContainer.TransformComponent.Orientation
		translation := mgl64.Translate3D(entityPosition.X(), entityPosition.Y(), entityPosition.Z())

		if !renderData.IsVisible() {
			continue
		}

		if modelRenderData, ok := renderData.(*components.ModelRenderData); ok {
			var meshModelMatrix mgl64.Mat4

			if modelRenderData.ID == "cowboy" || modelRenderData.ID == "box" || modelRenderData.ID == "cube_shifted" || modelRenderData.ID == "cube_static" {
				// // accounting for blender change of axis
				xr := mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
				yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
				meshModelMatrix = createModelMatrix(
					mgl64.Ident4(),
					orientation.Mat4().Mul4(yr).Mul4(xr),
					translation,
				)
			} else {
				// meshModelMatrix = createModelMatrix(
				// 	mgl64.Scale3D(0.07, 0.07, 0.07),
				// 	mgl64.Ident4(),
				// 	mgl64.Translate3D(entityPosition.X(), entityPosition.Y(), entityPosition.Z()),
				// )
				meshModelMatrix = createModelMatrix(
					mgl64.Ident4(),
					orientation.Mat4(),
					translation,
				)
			}

			drawModel(
				viewerContext,
				lightContext,
				s.shadowMap,
				shaderManager.GetShaderProgram(modelRenderData.ShaderProgram),
				componentContainer.MeshComponent,
				componentContainer.AnimationComponent,
				meshModelMatrix,
			)
		} else if _, ok := renderData.(*components.BlockRenderData); ok {
		}
	}
}

func (s *RenderSystem) Update(delta time.Duration) {
}
