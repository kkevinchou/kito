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

func NewRigidBody(position mgl64.Vec3) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Position:    position,
		Orientation: mgl64.QuatIdent(),
	}

	assetManager := directory.GetDirectory().AssetManager()
	modelSpec := assetManager.GetAnimatedModel("cowboy")
	var m *model.Model
	var vao uint32
	var vertexCount int
	var texture *textures.Texture

	if utils.IsClient() {
		m = model.NewMeshedModel(modelSpec)
		vao = m.VAO()
		vertexCount = m.VertexCount()
		texture = assetManager.GetTexture("character Texture")
	} else {
		m = model.NewPlaceholderModel(modelSpec)
	}

	meshComponent := &components.MeshComponent{
		ModelVAO:         vao,
		ModelVertexCount: vertexCount,
		Texture:          texture,
	}

	renderData := &components.ModelRenderData{
		ID:      "thing",
		Visible: true,
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
