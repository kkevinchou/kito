package behavior

import "time"

type AiStateModifierFunction func(state AiState)

type AiStateModifier struct {
	modifier AiStateModifierFunction
}

func NewAiStateModifier(f AiStateModifierFunction) *AiStateModifier {
	return &AiStateModifier{
		modifier: f,
	}
}

func (a *AiStateModifier) Tick(state AiState, delta time.Duration) Status {
	a.modifier(state)
	return SUCCESS
}
