package render

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kkevinchou/kito/kito/console"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
)

func (s *RenderSystem) networkInfoUIComponent() {
	metricsRegistry := s.world.MetricsRegistry()
	predictionMiss := int(metricsRegistry.GetOneSecondSum("predictionMiss"))
	// serverPosition := int(metricsRegistry.GetLatest("serverPositionDiff"))
	predictionHit := int(metricsRegistry.GetOneSecondSum("predictionHit"))
	ping := int(metricsRegistry.GetOneSecondAverage("ping"))
	// updateMessageSize := int(metricsRegistry.GetOneSecondSum("update_message_size")) / 1000
	// updateCount := int(metricsRegistry.GetOneSecondSum("update_message_count"))
	// newInput := int(metricsRegistry.GetOneSecondSum("newinput"))

	if imgui.CollapsingHeaderV("Network", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("Ping", ping)
		uiTableRow("Predictions Hit", predictionHit)
		uiTableRow("Predictions Miss", predictionMiss)
		// uiTableRow("Server Position", serverPosition)
		// uiTableRow("Update Count", updateCount)
		// uiTableRow("Update Size", updateMessageSize)
		// uiTableRow("Inputs Sent", newInput)
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
	// frameCatchup := int(metricsRegistry.GetOneSecondSum("frameCatchup"))
	if imgui.CollapsingHeaderV("General", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
		imgui.BeginTableV("", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		uiTableRow("FPS", fps)
		// uiTableRow("Frame Catchup", frameCatchup)
		uiTableRow("CF", s.world.CommandFrame())
		imgui.EndTable()
	}
}

func (s *RenderSystem) debugWindow() {
	imgui.SetNextWindowBgAlpha(0.5)
	imgui.BeginV("Debug", nil, imgui.WindowFlagsNoFocusOnAppearing)
	s.generalInfoComponent()
	s.networkInfoUIComponent()
	s.entityInfoUIComponent()
	// s.lightingUIComponent(s.shadowMap.DepthTexture())
	imgui.SetItemDefaultFocus()
	if imgui.IsWindowFocused() {
		s.world.SetFocusedWindow(types.WindowDebug)
	}
	imgui.End()
}

func inputFilterCallback(data imgui.InputTextCallbackData) int32 {
	if data.EventChar() == '`' {
		return 1
	}
	return 0
}

func (s *RenderSystem) consoleWindow() {
	// imgui.SetNextWindowFocus()
	// imgui.BeginV("Console", nil, imgui.WindowFlagsNoFocusOnAppearing)
	imgui.BeginV("Console", nil, imgui.WindowFlagsNone)

	imgui.PushItemWidth(-1)
	imgui.PushStyleColor(imgui.StyleColorFrameBg, imgui.Vec4{X: 0.5, Y: 0.5, Z: 0.5, W: 1})
	for i, consoleItem := range console.GlobalConsole.ConsoleItems {
		imgui.Textf("%d: %s", i, consoleItem.Command)
	}
	imgui.PopStyleColor()
	imgui.Separator()

	flags := imgui.InputTextFlagsEnterReturnsTrue | imgui.InputTextFlagsCallbackCharFilter
	value := imgui.InputTextV("input", &console.GlobalConsole.Input, flags, inputFilterCallback)
	if value {
		console.GlobalConsole.Send()
		imgui.SetKeyboardFocusHereV(-1)
	}

	imgui.PopItemWidth()
	if imgui.IsWindowFocused() {
		s.world.SetFocusedWindow(types.WindowConsole)
	}

	imgui.End()
}

func uiTableRow(label string, value any) {
	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	imgui.Text(label)
	imgui.TableSetColumnIndex(1)
	imgui.Text(fmt.Sprintf("%v", value))
}
