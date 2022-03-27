package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/textures"
)

const (
	startY = 50
)

func NewBob() *EntityImpl {
	modelName := "human"
	textureName := "color_grid"

	transformComponent := &components.TransformComponent{
		Position:    mgl64.Vec3{580, startY, -210},
		Orientation: mgl64.QuatIdent(),
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetAnimatedModel(modelName)

	var vao uint32
	var texture *textures.Texture

	m := model.NewModel(modelSpec)
	vertexCount := m.VertexCount()

	if utils.IsClient() {
		vao = m.Bind()
		texture = assetManager.GetTexture(textureName)
	}

	animationComponent := &components.AnimationComponent{
		CurrentAnimation: "Idle",
		Animation:        m.Animations["Idle"],
		Animations:       m.Animations,
	}
	_ = animationComponent

	yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
		// Scale:            mgl64.Scale3D(1, 1, 1),
		Scale: mgl64.Scale3D(10, 10, 10),
		// Orientation: mgl64.Ident4(),
		Orientation: yr,
		Material:    m.Mesh.Material(),
		PBRMaterial: m.Mesh.PBRMaterial(),
	}

	capsule := collider.NewCapsule(mgl64.Vec3{0, 12, 0}, mgl64.Vec3{0, 3, 0}, 3)
	colliderComponent := &components.ColliderComponent{
		CapsuleCollider: &capsule,
	}

	thirdPersonControllerComponent := &components.ThirdPersonControllerComponent{
		Controlled: true,
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		animationComponent,
		thirdPersonControllerComponent,
		meshComponent,
		colliderComponent,
		renderComponent,
	}

	entity := NewEntity(
		"bob",
		types.EntityTypeBob,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
