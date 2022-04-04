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

var (
	defaultXR          = mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
	defaultYR          = mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	defaultOrientation = defaultYR.Mul4(defaultXR)
	defaultScale       = mgl64.Scale3D(25, 25, 25)
)

func NewScene() *EntityImpl {
	// return NewRigidBody("scene_building", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeScene, "color_grid")
	return NewRigidBody("scene_building", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeScene, "color_grid")
}

func NewSlime() *EntityImpl {
	return NewRigidBody("slime_kevin", defaultScale, defaultOrientation, types.EntityTypeStaticSlime, "default")
}

func NewStaticRigidBody() *EntityImpl {
	return NewRigidBody("cubetest2", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeStaticRigidBody, "default")
}

func NewDynamicRigidBody() *EntityImpl {
	return NewRigidBody("guard", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeDynamicRigidBody, "color_grid")
}

func NewRigidBody(modelName string, Scale mgl64.Mat4, Orientation mgl64.Mat4, entityType types.EntityType, textureName string) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Orientation: mgl64.QuatIdent(),
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetAnimatedModel(modelName)

	var texture *textures.Texture

	m := model.NewModel(modelSpec)

	if utils.IsClient() {
		m.Prepare()
		texture = assetManager.GetTexture(textureName)
	}

	meshComponent := &components.MeshComponent{
		Texture:     texture,
		Scale:       Scale,
		Orientation: Orientation,
		Model:       m,
	}

	// triMesh := collider.NewTriMesh(m.Mesh.Vertices())
	// colliderComponent := &components.ColliderComponent{
	// 	TriMeshCollider: &triMesh,
	// }

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
		// colliderComponent,
		physicsComponent,
	}

	// if m.Animation != nil {
	// 	fmt.Println("rigid body with animation", modelName)
	// 	animationPlayer := animation.NewAnimationPlayer(m.Animations)
	// 	animationPlayer.PlayAnimation("Idle")
	// 	animationComponent := &components.AnimationComponent{
	// 		Player: animationPlayer,
	// 	}
	// 	componentList = append(componentList, animationComponent)
	// }

	entity := NewEntity(
		"rigidbody",
		entityType,
		components.NewComponentContainer(componentList...),
	)

	return entity
}
