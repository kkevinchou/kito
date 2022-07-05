package render

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kkevinchou/kito/kito/utils"
)

func (s *RenderSystem) networkInfoUIComponent() {
	metricsRegistry := s.world.MetricsRegistry()
	predictionMiss := int(metricsRegistry.GetOneSecondSum("predictionMiss"))
	serverPosition := int(metricsRegistry.GetLatest("serverPositionDiff"))
	predictionHit := int(metricsRegistry.GetOneSecondSum("predictionHit"))
	ping := int(metricsRegistry.GetOneSecondAverage("ping"))
	updateMessageSize := int(metricsRegistry.GetOneSecondSum("update_message_size")) / 1000
	updateCount := int(metricsRegistry.GetOneSecondSum("update_message_count"))
	newInput := int(metricsRegistry.GetOneSecondSum("newinput"))

	if imgui.CollapsingHeaderV("Network", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("Ping", ping)
		uiTableRow("Predictions Hit", predictionHit)
		uiTableRow("Predictions Miss", predictionMiss)
		uiTableRow("Server Position", serverPosition)
		uiTableRow("Update Count", updateCount)
		uiTableRow("Update Size", updateMessageSize)
		uiTableRow("Inputs Sent", newInput)
		imgui.EndTable()
	}
}

func (s *RenderSystem) lightingUIComponent(textureID uint32) {
	if imgui.CollapsingHeaderV("Lighting", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.ImageV(imgui.TextureID(textureID), imgui.Vec2{X: 160, Y: 90}, imgui.Vec2{X: 0, Y: 1}, imgui.Vec2{X: 1, Y: 0}, imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1}, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	}
}

func (s *RenderSystem) entityInfoUIComponent() {
	entity := s.world.GetPlayerEntity()
	componentContainer := entity.GetComponentContainer()
	entityPosition := componentContainer.TransformComponent.Position
	orientation := componentContainer.TransformComponent.Orientation
	velocity := componentContainer.ThirdPersonControllerComponent.Velocity

	if imgui.CollapsingHeaderV("Entity", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("ID", entity.GetID())
		uiTableRow("Position", utils.PPrintVec(entityPosition))
		uiTableRow("Velocity", utils.PPrintVec(velocity))
		uiTableRow("Orientation", utils.PPrintQuatAsVec(orientation))
		imgui.EndTable()
	}
}

func (s *RenderSystem) generalInfoComponent() {
	metricsRegistry := s.world.MetricsRegistry()
	fps := int(metricsRegistry.GetOneSecondSum("fps"))
	if imgui.CollapsingHeaderV("General", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("FPS", fps)
		uiTableRow("CF", s.world.CommandFrame())
		imgui.EndTable()
	}
}

func uiTableRow(label string, value any) {
	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	imgui.Text(label)
	imgui.TableSetColumnIndex(1)
	imgui.Text(fmt.Sprintf("%v", value))
}
