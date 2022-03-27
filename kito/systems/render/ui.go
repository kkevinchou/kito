package render

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kkevinchou/kito/kito/utils"
)

func (s *RenderSystem) networkInfoUIComponent() {
	metricsRegistry := s.world.MetricsRegistry()
	predictionMiss := int(metricsRegistry.GetOneSecondSum("predictionMiss"))
	predictionHit := int(metricsRegistry.GetOneSecondSum("predictionHit"))
	ping := int(metricsRegistry.GetOneSecondAverage("ping"))
	updateMessageSize := int(metricsRegistry.GetOneSecondSum("update_message_size")) / 1000
	updateCount := int(metricsRegistry.GetOneSecondSum("update_message_count"))
	newInput := int(metricsRegistry.GetOneSecondSum("newinput"))

	if imgui.CollapsingHeaderV("Network Info", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("Ping", ping)
		uiTableRow("Predictions Hit", predictionHit)
		uiTableRow("Predictions Miss", predictionMiss)
		uiTableRow("Update Count", updateCount)
		uiTableRow("Update Size", updateMessageSize)
		uiTableRow("Inputs Sent", newInput)
		imgui.EndTable()
	}
}

func (s *RenderSystem) entityInfoUIComponent() {
	entity := s.world.GetPlayerEntity()
	componentContainer := entity.GetComponentContainer()
	entityPosition := componentContainer.TransformComponent.Position
	orientation := componentContainer.TransformComponent.Orientation
	velocity := componentContainer.ThirdPersonControllerComponent.Velocity

	if imgui.CollapsingHeaderV("Entity Info", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
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
	if imgui.CollapsingHeaderV("General Info", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
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
