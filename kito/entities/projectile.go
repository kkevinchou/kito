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

	var texture *textures.Texture

	m := model.NewModel(modelSpec)
	m.Prepare()

	if utils.IsClient() {
		texture = assetManager.GetTexture(textureName)
	}

	// animationComponent := &components.AnimationComponent{
	// 	Animation: m.Animation,
	// }

	meshComponent := &components.MeshComponent{
		Texture:     texture,
		Scale:       mgl64.Ident4(),
		Orientation: mgl64.Ident4(),
		Model:       m,
	}

	// capsule := collider.NewCapsuleFromMeshVertices(m.Mesh.Vertices())
	// colliderComponent := &components.ColliderComponent{
	// 	CapsuleCollider: &capsule,
	// }

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
		// colliderComponent,
		renderComponent,
	}

	entity := NewEntity(
		"projectile",
		types.EntityTypeProjectile,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
