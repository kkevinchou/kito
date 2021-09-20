package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/textures"
	"github.com/kkevinchou/kito/types"
	"github.com/kkevinchou/kito/utils"
)

func NewBob(position mgl64.Vec3) *EntityImpl {
	modelName := "guard_running"
	shaderProgram := "model"
	textureName := "Guard_02__diffuse"

	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
	}

	renderData := &components.ModelRenderData{
		ID:            modelName,
		Visible:       true,
		ShaderProgram: shaderProgram,
	}
	renderComponent := &components.RenderComponent{
		RenderData: renderData,
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

	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
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
