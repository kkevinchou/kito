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

func NewProjectile(position mgl64.Vec3) *EntityImpl {
	modelName := "human"
	textureName := "color_grid"

	transformComponent := &components.TransformComponent{
		Position:    position,
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

	// animationComponent := &components.AnimationComponent{
	// 	Animation: m.Animation,
	// }

	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
		Scale:            mgl64.Ident4(),
		Orientation:      mgl64.Ident4(),
	}

	capsule := collider.NewCapsuleFromMeshVertices(m.Mesh.Vertices())
	colliderComponent := &components.ColliderComponent{
		CapsuleCollider: &capsule,
	}

	physicsComponent := &components.PhysicsComponent{
		IgnoreGravity: true,
		Impulses:      map[string]types.Impulse{},
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		// animationComponent,
		physicsComponent,
		meshComponent,
		colliderComponent,
		renderComponent,
	}

	entity := NewEntity(
		"projectile",
		types.EntityTypeProjectile,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
