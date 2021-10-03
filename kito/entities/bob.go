package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/textures"
)

func NewBob(position mgl64.Vec3) *EntityImpl {
	modelName := "guard_running"
	shaderProgram := "model"
	textureName := "default"

	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetAnimatedModel(modelName)

	var m *model.Model
	var vao uint32
	var vertexCount int
	var texture *textures.Texture

	if utils.IsClient() {
		m = model.NewModel(modelSpec)
		vao = m.VAO()
		vertexCount = m.VertexCount()
		texture = assetManager.GetTexture(textureName)
	} else {
		m = model.NewPlaceholderModel(modelSpec)
	}

	animationComponent := &components.AnimationComponent{
		Animation: m.Animation,
	}
	_ = animationComponent

	yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
		ShaderProgram:    shaderProgram,
		Scale:            mgl64.Scale3D(0.07, 0.07, 0.07),
		Orientation:      yr,
	}

	physicsComponent := &components.PhysicsComponent{
		Impulses: map[string]types.Impulse{},
	}

	thirdPersonControllerComponent := &components.ThirdPersonControllerComponent{
		Controlled: true,
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		animationComponent,
		physicsComponent,
		thirdPersonControllerComponent,
		meshComponent,
	}

	if utils.IsClient() {
		entityComponents = append(entityComponents, renderComponent)
	}

	entity := NewEntity(
		"bob",
		types.EntityTypeBob,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
