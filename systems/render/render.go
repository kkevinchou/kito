package render

import (
	"fmt"
	_ "image/png"
	"math"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/singleton"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fovx float64 = 90
	near float64 = 1
	far  float64 = 500

	// shadow map parameters
	shadowMapDimension   int     = 8000
	shadowDistanceFactor float64 = .8 // proportion of view fustrum to include in shadow cuboid
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
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

	entities []entities.Entity
}

func NewRenderSystem(world World, window *sdl.Window, width, height int) *RenderSystem {
	aspectRatio := float64(width) / float64(height)
	renderSystem := RenderSystem{
		BaseSystem: &base.BaseSystem{},
		window:     window,
		world:      world,
		skybox:     NewSkyBox(600),
		floor:      NewQuad(quadZeroY),

		width:       width,
		height:      height,
		aspectRatio: aspectRatio,
		fovY:        mgl64.RadToDeg(2 * math.Atan(math.Tan(mgl64.DegToRad(fovx)/2)/aspectRatio)),
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

	var err error
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
		ProjectionMatrix:  mgl64.Perspective(mgl64.DegToRad(s.fovY), s.aspectRatio, near, far),
	}

	// configure light viewer context
	modelSpaceFrustumPoints := CalculateFrustumPoints(transformComponent.Position, transformComponent.Orientation, near, far, s.fovY, s.aspectRatio, shadowDistanceFactor)
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

	gl.Viewport(0, 0, int32(s.width), int32(s.height))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.renderScene(viewerContext, lightContext)
}

// renderScene renders a scene from the perspective of a viewer
func (s *RenderSystem) renderScene(viewerContext ViewerContext, lightContext LightContext) {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	assetManager := d.AssetManager()

	// render a debug shadow map for viewing
	drawHUDTextureToQuad(viewerContext, shaderManager.GetShaderProgram("depthDebug"), s.shadowMap.DepthTexture(), 0.4)

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
		assetManager.GetTexture("default"),
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

			if strings.Contains(modelRenderData.ID, "slime") {
				// // accounting for blender change of axis
				xr := mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
				yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
				meshModelMatrix = createModelMatrix(
					mgl64.Scale3D(25, 25, 25),
					orientation.Mat4().Mul4(yr).Mul4(xr),
					translation,
				)
			} else if strings.Contains(modelRenderData.ID, "guard") {
				// meshModelMatrix = createModelMatrix(
				// 	mgl64.Scale3D(0.07, 0.07, 0.07),
				// 	mgl64.Ident4(),
				// 	mgl64.Translate3D(entityPosition.X(), entityPosition.Y(), entityPosition.Z()),
				// )
				yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
				meshModelMatrix = createModelMatrix(
					mgl64.Scale3D(0.07, 0.07, 0.07),
					orientation.Mat4().Mul4(yr),
					translation,
				)
			} else {
				// // accounting for blender change of axis
				xr := mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
				yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
				meshModelMatrix = createModelMatrix(
					mgl64.Ident4(),
					orientation.Mat4().Mul4(yr).Mul4(xr),
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
