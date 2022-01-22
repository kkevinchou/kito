package entities

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/textures"
)

var (
	defaultXR          = mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0}).Mat4()
	defaultYR          = mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	defaultOrientation = defaultYR.Mul4(defaultXR)
	defaultScale       = mgl64.Scale3D(25, 25, 25)
)

func NewScene(position mgl64.Vec3) *EntityImpl {
	shaderProgram := "model_static"
	return NewRigidBody(position, "sceneuv", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeScene, "color_grid", shaderProgram)
}

func NewSlime(position mgl64.Vec3) *EntityImpl {
	shaderProgram := "model"
	return NewRigidBody(position, "slime_kevin", defaultScale, defaultOrientation, types.EntityTypeStaticSlime, "default", shaderProgram)
}

func NewStaticRigidBody(position mgl64.Vec3) *EntityImpl {
	shaderProgram := "model_static"
	return NewRigidBody(position, "cubetest2", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeStaticRigidBody, "default", shaderProgram)
}

func NewDynamicRigidBody(position mgl64.Vec3) *EntityImpl {
	shaderProgram := "model"
	return NewRigidBody(position, "guard", mgl64.Ident4(), mgl64.Ident4(), types.EntityTypeDynamicRigidBody, "color_grid", shaderProgram)
}

func NewRigidBody(position mgl64.Vec3, modelName string, Scale mgl64.Mat4, Orientation mgl64.Mat4, entityType types.EntityType, textureName string, shaderProgram string) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
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

	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
		ShaderProgram:    shaderProgram,
		Scale:            Scale,
		Orientation:      Orientation,
		Material:         m.Mesh.Material(),
	}

	// triMesh := collider.NewBoxTriMesh(40, 50, 20)
	triMesh := collider.NewTriMesh(m.Mesh.Vertices())
	colliderComponent := &components.ColliderComponent{
		TriMeshCollider: &triMesh,
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

	if m.Animation != nil {
		fmt.Println("rigid body with animation", modelName)
		animationComponent := &components.AnimationComponent{
			Animation: m.Animation,
		}
		componentList = append(componentList, animationComponent)
	}

	entity := NewEntity(
		"rigidbody",
		entityType,
		components.NewComponentContainer(componentList...),
	)

	return entity
}
