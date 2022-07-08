package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
)

func NewProjectile(position mgl64.Vec3) *EntityImpl {
	modelName := "fireball"

	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetModel(modelName)

	m := model.NewModel(modelSpec)

	meshComponent := &components.MeshComponent{
		Scale:       mgl64.Ident4(),
		Orientation: mgl64.Ident4(),
		Model:       m,
	}

	capsule := collider.NewCapsuleFromModel(m)
	colliderComponent := &components.ColliderComponent{
		SkipMovementResolution: true,
		CapsuleCollider:        &capsule,
		Contacts:               map[int]*collision.Contact{},
	}

	physicsComponent := &components.PhysicsComponent{
		IgnoreGravity: true,
		Impulses:      map[string]types.Impulse{},
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
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
