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

func NewRigidBody(position mgl64.Vec3) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetAnimatedModel("slime_kevin")
	var m *model.Model
	var vao uint32
	var vertexCount int
	var texture *textures.Texture

	if utils.IsClient() {
		m = model.NewMeshedModel(modelSpec)
		vao = m.VAO()
		vertexCount = m.VertexCount()
		texture = assetManager.GetTexture("default")
	} else {
		m = model.NewPlaceholderModel(modelSpec)
	}

	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
	}

	renderData := &components.ModelRenderData{
		ID:            "slime_kevin",
		Visible:       true,
		ShaderProgram: "model_static",
	}
	renderComponent := &components.RenderComponent{
		RenderData: renderData,
	}

	entity := NewEntity(
		"rigidbody",
		types.EntityTypeRigidBody,
		components.NewComponentContainer(
			transformComponent,
			renderComponent,
			&components.NetworkComponent{},
			meshComponent,
		),
	)

	return entity
}
