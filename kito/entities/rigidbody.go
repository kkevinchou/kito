package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
)

var (
	defaultXR          = mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
	defaultYR          = mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	defaultOrientation = defaultYR.Mul4(defaultXR)
	defaultScale       = mgl64.Scale3D(25, 25, 25)
)

func NewScene() *EntityImpl {
	return NewRigidBody("scene", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeScene)
	// return NewRigidBody("scene_giga_flat", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeScene)
}

func NewSlime() *EntityImpl {
	return NewRigidBody("slime_kevin", defaultScale, defaultOrientation, types.EntityTypeStaticSlime)
}

func NewStaticRigidBody() *EntityImpl {
	return NewRigidBody("cubetest2", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeStaticRigidBody)
}

func NewDynamicRigidBody() *EntityImpl {
	return NewRigidBody("guard", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeDynamicRigidBody)
}

func NewRigidBody(modelName string, Scale mgl64.Mat4, Orientation mgl64.Mat4, entityType types.EntityType) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Orientation: mgl64.QuatIdent(),
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetModel(modelName)

	m := model.NewModel(modelSpec)

	meshComponent := &components.MeshComponent{
		Scale:       Scale,
		Orientation: Orientation,
		Model:       m,
	}

	triMesh := collider.NewTriMesh(m)
	boundingBox := collider.BoundingBoxFromModel(m)

	colliderComponent := &components.ColliderComponent{
		TriMeshCollider:     &triMesh,
		BoundingBoxCollider: boundingBox,
		Contacts:            map[int]bool{},
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	physicsComponent := &components.PhysicsComponent{
		Static: true,
	}

	componentList := []components.Component{
		transformComponent,
		renderComponent,
		&components.NetworkComponent{},
		meshComponent,
		colliderComponent,
		physicsComponent,
	}

	entity := NewEntity(
		"rigidbody",
		entityType,
		components.NewComponentContainer(componentList...),
	)

	return entity
}
