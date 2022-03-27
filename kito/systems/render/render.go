package render

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/metrics"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fovx float64 = 90
	near float64 = 1
	far  float64 = 3000

	// shadow map parameters
	shadowMapDimension   int     = 15000
	shadowDistanceFactor float64 = .4 // proportion of view fustrum to include in shadow cuboid
	shadowmapZOffset             = 400
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayerEntity() entities.Entity
	MetricsRegistry() *metrics.MetricsRegistry
	CommandFrame() int
}

type Platform interface {
	NewFrame()
	DisplaySize() [2]float32
	FramebufferSize() [2]float32
}

type RenderSystem struct {
	*base.BaseSystem
	window    *sdl.Window
	world     World
	skybox    *SkyBox
	floor     *Quad
	shadowMap *ShadowMap

	width       int
	height      int
	aspectRatio float64
	fovY        float64

	imguiRenderer *ImguiOpenGL4Renderer
	platform      Platform

	entities []entities.Entity
}

func init() {
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
}

func NewRenderSystem(world World, window *sdl.Window, platform Platform, imguiIO imgui.IO, width, height int) *RenderSystem {
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
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	aspectRatio := float64(width) / float64(height)
	shadowMap, err := NewShadowMap(shadowMapDimension, shadowMapDimension, far*shadowDistanceFactor)
	if err != nil {
		panic(fmt.Sprintf("failed to create shadow map %s", err))
	}

	imguiRenderer, err := NewImguiOpenGL4Renderer(imguiIO)
	if err != nil {
		panic(err)
	}

	renderSystem := RenderSystem{
		BaseSystem: &base.BaseSystem{},
		window:     window,
		world:      world,
		skybox:     NewSkyBox(float32(far)),
		floor:      NewQuad(quadZeroY),
		shadowMap:  shadowMap,

		width:       width,
		height:      height,
		aspectRatio: aspectRatio,
		fovY:        mgl64.RadToDeg(2 * math.Atan(math.Tan(mgl64.DegToRad(fovx)/2)/aspectRatio)),

		platform:      platform,
		imguiRenderer: imguiRenderer,
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
		return
	}

	// configure camera viewer context
	viewerViewMatrix := transformComponent.Orientation.Mat4()
	viewTranslationMatrix := mgl64.Translate3D(transformComponent.Position.X(), transformComponent.Position.Y(), transformComponent.Position.Z())

	cameraViewerContext := ViewerContext{
		Position:    transformComponent.Position,
		Orientation: transformComponent.Orientation,

		InverseViewMatrix: viewTranslationMatrix.Mul4(viewerViewMatrix).Inv(),
		ProjectionMatrix:  mgl64.Perspective(mgl64.DegToRad(s.fovY), s.aspectRatio, near, far),
	}

	// configure light viewer context
	modelSpaceFrustumPoints := CalculateFrustumPoints(transformComponent.Position, transformComponent.Orientation, near, far, s.fovY, s.aspectRatio, shadowDistanceFactor)
	// NOTE: for some reason, using a negative angle makes shadow calculation wonky.
	lightOrientation := mgl64.QuatRotate(mgl64.DegToRad(45), mgl64.Vec3{0, 1, 0})
	lightOrientation = lightOrientation.Mul(mgl64.QuatRotate(mgl64.DegToRad(310), mgl64.Vec3{1, 0, 0}))
	lightPosition, lightProjectionMatrix := ComputeDirectionalLightProps(lightOrientation.Mat4(), modelSpaceFrustumPoints, shadowmapZOffset)
	lightViewMatrix := mgl64.Translate3D(lightPosition.X(), lightPosition.Y(), lightPosition.Z()).Mul4(lightOrientation.Mat4()).Inv()

	lightViewerContext := ViewerContext{
		Position:    lightPosition,
		Orientation: lightOrientation,

		InverseViewMatrix: lightViewMatrix,
		ProjectionMatrix:  lightProjectionMatrix,
	}

	lightContext := LightContext{
		DirectionalLightDir: lightOrientation.Rotate(mgl64.Vec3{0, 0, -1}),
		// this should be the inverse of the transforms applied to the viewer context
		// if the viewer moves along -y, the universe moves along +y
		LightSpaceMatrix: lightProjectionMatrix.Mul4(lightViewMatrix),
	}

	s.renderToDepthMap(lightViewerContext, lightContext)
	s.renderToDisplay(cameraViewerContext, lightContext)
	s.platform.NewFrame()
	s.renderImgui()

	s.window.GLSwap()
}

func (s *RenderSystem) renderToDepthMap(viewerContext ViewerContext, lightContext LightContext) {
	defer resetGLRenderSettings()
	s.shadowMap.Prepare()

	s.renderScene(viewerContext, lightContext, true)
}

func (s *RenderSystem) renderToDisplay(viewerContext ViewerContext, lightContext LightContext) {
	defer resetGLRenderSettings()

	gl.Viewport(0, 0, int32(s.width), int32(s.height))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.renderScene(viewerContext, lightContext, false)
}

var f [3]float32
var inputText string

func (s *RenderSystem) renderImgui() {
	imgui.NewFrame()

	imgui.SetNextWindowBgAlpha(0.5)
	imgui.BeginV("Console", nil, imgui.WindowFlagsNoFocusOnAppearing)
	s.generalInfoComponent()
	s.networkInfoUIComponent()
	s.entityInfoUIComponent()
	imgui.End()

	imgui.Render()
	s.imguiRenderer.Render(s.platform.DisplaySize(), s.platform.FramebufferSize(), imgui.RenderedDrawData())
}

// renderScene renders a scene from the perspective of a viewer
func (s *RenderSystem) renderScene(viewerContext ViewerContext, lightContext LightContext, shadowPass bool) {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	assetManager := d.AssetManager()

	// render a debug shadow map for viewing
	// drawHUDTextureToQuad(viewerContext, shaderManager.GetShaderProgram("depthDebug"), s.shadowMap.DepthTexture(), 0.4)
	// drawHUDTextureToQuad(viewerContext, shaderManager.GetShaderProgram("quadtex"), textTexture, 0.4)

	if !shadowPass {
		drawSkyBox(
			viewerContext,
			s.skybox,
			shaderManager.GetShaderProgram("skybox"),
			assetManager.GetTexture("front"),
			assetManager.GetTexture("top"),
			assetManager.GetTexture("left"),
			assetManager.GetTexture("right"),
			assetManager.GetTexture("bottom"),
			assetManager.GetTexture("back"),
		)
	}

	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		entityPosition := componentContainer.TransformComponent.Position
		orientation := componentContainer.TransformComponent.Orientation
		translation := mgl64.Translate3D(entityPosition.X(), entityPosition.Y(), entityPosition.Z())
		renderComponent := componentContainer.RenderComponent
		meshComponent := componentContainer.MeshComponent

		if !renderComponent.IsVisible {
			continue
		}

		meshModelMatrix := createModelMatrix(
			meshComponent.Scale,
			orientation.Mat4().Mul4(meshComponent.Orientation),
			// mgl64.Ident4(),
			translation,
		)

		shader := "model_static"
		if componentContainer.AnimationComponent != nil {
			shader = "model"
		}

		drawModel(
			viewerContext,
			lightContext,
			s.shadowMap,
			shaderManager.GetShaderProgram(shader),
			componentContainer.MeshComponent,
			componentContainer.AnimationComponent,
			meshModelMatrix,
			orientation.Mat4().Mul4(meshComponent.Orientation),
		)

		if settings.DebugRenderCollisionVolume {
			if componentContainer.ColliderComponent != nil {
				if componentContainer.ColliderComponent.CapsuleCollider != nil {
					// lots of hacky rendering stuff to get the rectangle to billboard
					center := mgl64.Vec3{componentContainer.TransformComponent.Position.X(), 0, componentContainer.TransformComponent.Position.Z()}
					viewerArtificialCenter := mgl64.Vec3{viewerContext.Position.X(), 0, viewerContext.Position.Z()}
					vecToViewer := viewerArtificialCenter.Sub(center).Normalize()
					billboardModelMatrix := translation.Mul4(mgl64.QuatBetweenVectors(mgl64.Vec3{0, 0, 1}, vecToViewer).Mat4())
					drawCapsuleCollider(
						viewerContext,
						lightContext,
						shaderManager.GetShaderProgram("basicsolid"),
						componentContainer.ColliderComponent.CapsuleCollider,
						billboardModelMatrix,
					)
				} else if componentContainer.ColliderComponent.TransformedTriMeshCollider != nil {
					drawTriMeshCollider(
						viewerContext,
						lightContext,
						shaderManager.GetShaderProgram("basicsolid"),
						componentContainer.ColliderComponent.TransformedTriMeshCollider,
					)
				}
			}
		}
	}

	// var renderText string
	// drawText(shaderManager.GetShaderProgram("quadtex"), assetManager.GetFont("robotomono-regular"), renderText, 0.8, 0)
}

func (s *RenderSystem) Update(delta time.Duration) {
}
