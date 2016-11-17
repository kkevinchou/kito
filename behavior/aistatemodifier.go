package behavior

import "time"

type AiStateModifierFunction func(state AIState)

type AiStateModifier struct {
	modifier AiStateModifierFunction
}

func NewAiStateModifier(f AiStateModifierFunction) *AiStateModifier {
	return &AiStateModifier{
		modifier: f,
	}
}

func (a *AiStateModifier) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	a.modifier(state)
	return nil, SUCCESS
}

func (a *AiStateModifier) Reset() {}
